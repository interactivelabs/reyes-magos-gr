// Command classify-toy-ages reads the JSON files written by scrape-toys and
// asks TypeSafe (Jev) for a realistic age range for each toy, since the
// manufacturer's "Edad recomendada" on Amazon is often wrong ("100 años y
// más", "5 meses y más" for a remote-control dinosaur...).
//
// Usage:
//
//	TYPESAFE_API_KEY=... go run ./scripts/classify-toy-ages -in ./scraped
//
// It uses the same TURSO_* environment variables as the app to read each
// toy's current age range; nothing is written to the database.
//
// The stages and questions sent to the model live in rules.json, next to this
// file, so they can evolve without touching the code; pass -rules to try
// another file. The model only picks semantic options (developmental stages,
// yes/no); turning those into age_min/age_max numbers and applying the
// safety rules happens here in code.
package main

import (
	"cmp"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/scripts/internal/typesafe"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
)

type ToyFile struct {
	ToyID    int64    `json:"toy_id"`
	ToyName  string   `json:"toy_name"`
	Title    string   `json:"title"`
	Features []string `json:"features"`
	Details  []struct {
		Name       string         `json:"name"`
		Attributes map[string]any `json:"attributes"`
	} `json:"details"`
}

// Result is what gets written for each toy.
type Result struct {
	ToyID              int64              `json:"toy_id"`
	ToyName            string             `json:"toy_name"`
	ManufacturerAge    string             `json:"manufacturer_age,omitempty"`
	CurrentAgeMin      *int64             `json:"current_age_min,omitempty"`
	CurrentAgeMax      *int64             `json:"current_age_max,omitempty"`
	AgeMin             int64              `json:"age_min"`
	AgeMax             int64              `json:"age_max"`
	NeedsReview        bool               `json:"needs_review"`
	ReviewReasons      []string           `json:"review_reasons,omitempty"`
	YoungestConfidence float64            `json:"youngest_confidence"`
	OldestConfidence   float64            `json:"oldest_confidence"`
	SmallParts         float64            `json:"small_parts"`
	ManufacturerOK     *float64           `json:"manufacturer_plausible,omitempty"`
	YoungestProbs      map[string]float64 `json:"youngest_probabilities"`
	OldestProbs        map[string]float64 `json:"oldest_probabilities"`
}

//go:embed rules.json
var defaultRules []byte

// Rules is the content of rules.json.
//
// Stages describe each floor (youngest) by the skills the toy demands and
// each ceiling (oldest) by the kind of toy and why older kids drop it, so
// options differ in meaning, not number: Jev tells "5 to 6" from "7 to 8"
// poorly, but tells a push truck from a karaoke machine well.
type Rules struct {
	// ManufacturerAgeAttribute is the scraped attribute holding the
	// manufacturer's recommended age, sent apart as manufacturer_age.
	ManufacturerAgeAttribute string `json:"manufacturer_age_attribute"`
	// NoisyAttributes carry no information about who the toy is for; any
	// attribute whose name contains one of them is dropped.
	NoisyAttributes []string `json:"noisy_attributes"`
	YoungestStages  []Stage  `json:"youngest_stages"`
	OldestStages    []Stage  `json:"oldest_stages"`
	// Questions must include youngest and oldest (their criteria come from
	// the stages) and small_parts.
	Questions map[string]typesafe.Question `json:"questions"`
	// ManufacturerQuestions are only asked when the toy has a manufacturer
	// age; manufacturer_plausible is read back into the results.
	ManufacturerQuestions map[string]typesafe.Question `json:"manufacturer_questions"`
}

// Stage is a Choice option: the model sees the key and description, code
// uses the age.
type Stage struct {
	Key         string `json:"key"`
	Age         int64  `json:"age"`
	Description any    `json:"description"`
}

func loadRules(path string) (Rules, error) {
	data := defaultRules
	if path != "" {
		var err error
		if data, err = os.ReadFile(path); err != nil {
			return Rules{}, err
		}
	}
	var r Rules
	if err := json.Unmarshal(data, &r); err != nil {
		return Rules{}, err
	}
	if len(r.YoungestStages) == 0 || len(r.OldestStages) == 0 {
		return Rules{}, errors.New("rules need youngest_stages and oldest_stages")
	}
	for _, k := range []string{"youngest", "oldest", "small_parts"} {
		if _, ok := r.Questions[k]; !ok {
			return Rules{}, fmt.Errorf("rules are missing the %q question", k)
		}
	}
	return r, nil
}

func main() {
	inDir := flag.String("in", "scraped", "directory with the toy-*.json files from scrape-toys")
	out := flag.String("out", "scraped/ages.json", "where to write the results")
	rulesPath := flag.String("rules", "", "rules file to use instead of the embedded rules.json")
	minConfidence := flag.Float64("min-confidence", 0.5, "flag toys whose youngest/oldest confidence is below this for review")
	workers := flag.Int("workers", 4, "concurrent requests")
	flag.Parse()

	rules, err := loadRules(*rulesPath)
	if err != nil {
		log.Fatal("Error loading rules: ", err)
	}
	client, err := typesafe.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	current := loadCurrent()

	paths, err := filepath.Glob(filepath.Join(*inDir, "toy-*.json"))
	if err != nil {
		log.Fatal(err)
	}

	results := make([]Result, len(paths))
	errs := make([]error, len(paths))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range *workers {
		wg.Go(func() {
			for i := range jobs {
				results[i], errs[i] = classify(client, rules, paths[i], *minConfidence)
			}
		})
	}
	for i := range paths {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	var ok []Result
	for i, err := range errs {
		if err != nil {
			log.Printf("%s: %v", paths[i], err)
			continue
		}
		r := results[i]
		if t, found := current[r.ToyID]; found {
			r.CurrentAgeMin, r.CurrentAgeMax = &t.AgeMin, &t.AgeMax
		}
		ok = append(ok, r)
	}
	sort.Slice(ok, func(i, j int) bool { return ok[i].ToyID < ok[j].ToyID })

	data, err := json.MarshalIndent(ok, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatal(err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tACTUAL\tRANGE\tMIN CONF\tMAX CONF\tSMALL PARTS\tFABRICANTE\tFAB OK\tREVIEW\tTOY")
	for _, r := range ok {
		manufacturerOK := "-"
		if r.ManufacturerOK != nil {
			manufacturerOK = fmt.Sprintf("%.2f", *r.ManufacturerOK)
		}
		currentRange := "-"
		if r.CurrentAgeMin != nil {
			currentRange = fmt.Sprintf("%d-%d", *r.CurrentAgeMin, *r.CurrentAgeMax)
		}
		fmt.Fprintf(tw, "%d\t%s\t%d-%d\t%.2f\t%.2f\t%.2f\t%s\t%s\t%s\t%s\n",
			r.ToyID, currentRange, r.AgeMin, r.AgeMax, r.YoungestConfidence, r.OldestConfidence, r.SmallParts,
			cmp.Or(r.ManufacturerAge, "-"), manufacturerOK, strings.Join(r.ReviewReasons, "; "), truncate(r.ToyName, 50))
	}
	tw.Flush()
	log.Printf("classified %d/%d toys → %s", len(ok), len(paths), *out)
}

// loadCurrent returns the toys in the database by id, so the proposed ranges
// can be compared with what the catalog has today.
func loadCurrent() map[int64]models.Toy {
	db, connector, dir, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	defer connector.Close()
	defer db.Close()
	if _, err := connector.Sync(); err != nil {
		log.Fatal("Error syncing database: ", err)
	}

	toys, err := store.NewToysStore(db).GetToys()
	if err != nil {
		log.Fatal("Error reading toys: ", err)
	}
	current := make(map[int64]models.Toy, len(toys))
	for _, t := range toys {
		current[t.ToyID] = t
	}
	return current
}

func classify(client *typesafe.Client, rules Rules, path string, minConfidence float64) (Result, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Result{}, err
	}
	var toy ToyFile
	if err := json.Unmarshal(raw, &toy); err != nil {
		return Result{}, err
	}

	state, manufacturerAge := buildState(rules, toy)
	questions := maps.Clone(rules.Questions)
	questions["youngest"] = questions["youngest"].WithCriteria(stageCriteria(rules.YoungestStages))
	questions["oldest"] = questions["oldest"].WithCriteria(stageCriteria(rules.OldestStages))
	if manufacturerAge != "" {
		maps.Copy(questions, rules.ManufacturerQuestions)
	}

	resp, err := client.Ask(state, questions)
	if err != nil {
		return Result{}, err
	}

	y, o := resp.Answers["youngest"], resp.Answers["oldest"]
	r := Result{
		ToyID:              toy.ToyID,
		ToyName:            toy.ToyName,
		ManufacturerAge:    manufacturerAge,
		AgeMin:             stageAge(rules.YoungestStages, y.Choice),
		AgeMax:             stageAge(rules.OldestStages, o.Choice),
		YoungestConfidence: y.Confidence,
		OldestConfidence:   o.Confidence,
		SmallParts:         resp.Answers["small_parts"].Noul,
		YoungestProbs:      y.Probabilities,
		OldestProbs:        o.Probabilities,
	}
	if a, ok := resp.Answers["manufacturer_plausible"]; ok {
		r.ManufacturerOK = &a.Noul
	}

	if r.SmallParts > 0.5 && r.AgeMin < 3 {
		r.AgeMin = 3
		r.review("small parts: raised age_min to 3")
	}
	if r.AgeMax <= r.AgeMin {
		r.AgeMax = r.AgeMin + 2
		r.review("oldest stage was not above youngest")
	}
	// Splitting between two neighbouring stages (preschool vs early_school)
	// moves the range by a year or two and is fine; only flag when the
	// probability is spread further than that.
	if r.YoungestConfidence < minConfidence && adjacentMass(rules.YoungestStages, y.Probabilities) < 0.8 {
		r.review(fmt.Sprintf("uncertain youngest (%.2f)", r.YoungestConfidence))
	}
	if r.OldestConfidence < minConfidence && adjacentMass(rules.OldestStages, o.Probabilities) < 0.8 {
		r.review(fmt.Sprintf("uncertain oldest (%.2f)", r.OldestConfidence))
	}
	return r, nil
}

func (r *Result) review(reason string) {
	r.NeedsReview = true
	r.ReviewReasons = append(r.ReviewReasons, reason)
}

// buildState keeps only what helps judge who the toy is for, so the model
// isn't distracted by ASINs, rankings and measurements.
func buildState(rules Rules, toy ToyFile) (map[string]any, string) {
	attrs := map[string]any{}
	var manufacturerAge string
	for _, d := range toy.Details {
		for k, v := range d.Attributes {
			if k == rules.ManufacturerAgeAttribute {
				manufacturerAge, _ = v.(string)
				continue
			}
			if !isNoisy(rules.NoisyAttributes, k) {
				attrs[k] = v
			}
		}
	}
	state := map[string]any{
		"name":       toy.ToyName,
		"title":      toy.Title,
		"features":   toy.Features,
		"attributes": attrs,
	}
	if manufacturerAge != "" {
		state["manufacturer_age"] = manufacturerAge
	}
	return state, manufacturerAge
}

func isNoisy(noisy []string, key string) bool {
	for _, n := range noisy {
		if strings.Contains(key, n) {
			return true
		}
	}
	return false
}

func stageCriteria(stages []Stage) map[string]any {
	c := make(map[string]any, len(stages))
	for _, s := range stages {
		c[s.Key] = s.Description
	}
	return c
}

// adjacentMass is the highest probability held by two neighbouring stages.
func adjacentMass(stages []Stage, probs map[string]float64) float64 {
	var best float64
	for i := 1; i < len(stages); i++ {
		best = max(best, probs[stages[i-1].Key]+probs[stages[i].Key])
	}
	return best
}

func stageAge(stages []Stage, key string) int64 {
	for _, s := range stages {
		if s.Key == key {
			return s.Age
		}
	}
	return 0
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}
