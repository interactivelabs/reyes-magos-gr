// Command classify-toy-categories reads the JSON files written by scrape-toys
// and asks TypeSafe (Jev) which catalog categories each toy belongs to, and
// which new categories would be worth creating.
//
// Usage:
//
//	TYPESAFE_API_KEY=... go run ./scripts/classify-toy-categories -in ./scraped
//
// It uses the same TURSO_* environment variables as the app to read the
// current categories; nothing is written to the database.
//
// The category descriptions, candidates and questions live in rules.json,
// next to this file, so they can evolve without touching the code; pass
// -rules to try another file.
//
// Jev picks from options, it doesn't invent names, so new categories come
// from candidate_categories in the rules:
//  1. One request aligns every candidate with the existing categories and
//     drops those the catalog already covers ("Vehículos" vs "Carros").
//  2. One request per toy asks a yes/no question per existing category and
//     per remaining candidate; a toy can belong to several.
//  3. Code assigns existing categories over the threshold and suggests a
//     candidate when enough toys need it.
package main

import (
	"cmp"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/scripts/internal/typesafe"
	"reyes-magos-gr/store"
)

//go:embed rules.json
var defaultRules []byte

// Rules is the content of rules.json.
type Rules struct {
	// ExistingDescriptions explains the categories already in the database,
	// so the model doesn't have to guess from terse (or misspelled) names.
	// The keys must match the database exactly because the catalog filters
	// with LIKE. Categories missing here are sent with their name only.
	ExistingDescriptions map[string]string `json:"existing_descriptions"`
	// Placeholders are values stored in the category column that aren't
	// real categories.
	Placeholders []string `json:"placeholders"`
	// CandidateCategories are names the script may suggest. Add freely: the
	// alignment step drops any that an existing category already covers.
	CandidateCategories map[string]string `json:"candidate_categories"`
	// RelevantAttributes are the scraped attributes that say what kind of
	// toy it is.
	RelevantAttributes []string `json:"relevant_attributes"`
	// RankingAttribute holds Amazon's bestseller rankings; RootRankings are
	// its top-level departments, too broad to say anything.
	RankingAttribute string   `json:"ranking_attribute"`
	RootRankings     []string `json:"root_rankings"`
	// AlignQuestion is asked once per candidate, with the candidate added to
	// its instructions and the existing categories to its criteria.
	AlignQuestion typesafe.Question `json:"align_question"`
	// BelongsQuestion is asked per toy for every category, with the category
	// added to its instructions.
	BelongsQuestion typesafe.Question `json:"belongs_question"`
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
	return r, json.Unmarshal(data, &r)
}

// rankingCategory pulls "Excavadoras Infantiles" out of "nº3 en Excavadoras Infantiles".
var rankingCategory = regexp.MustCompile(`^nº[\d,.]+ en (.+?)(?: \(.*\))?$`)

type ToyFile struct {
	ToyID    int64    `json:"toy_id"`
	ToyName  string   `json:"toy_name"`
	Title    string   `json:"title"`
	Features []string `json:"features"`
	Details  []struct {
		Attributes map[string]any `json:"attributes"`
	} `json:"details"`
}

type ToyResult struct {
	ToyID      int64              `json:"toy_id"`
	ToyName    string             `json:"toy_name"`
	Current    string             `json:"current_category"`
	Proposed   string             `json:"proposed_category"`
	Assigned   []string           `json:"assigned"`
	Uncertain  []string           `json:"uncertain,omitempty"`
	Suggested  []string           `json:"suggested_new,omitempty"`
	Uncovered  bool               `json:"uncovered"`
	Existing   map[string]float64 `json:"existing_scores"`
	Candidates map[string]float64 `json:"candidate_scores"`
}

type Suggestion struct {
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Toys        []int64 `json:"toys"`
	Uncovered   []int64 `json:"uncovered_toys"`
	Recommended bool    `json:"recommended"`
}

type Output struct {
	Existing    []string          `json:"existing_categories"`
	Duplicates  map[string]string `json:"candidates_already_covered"`
	Suggestions []Suggestion      `json:"suggestions"`
	Toys        []ToyResult       `json:"toys"`
}

func main() {
	inDir := flag.String("in", "scraped", "directory with the toy-*.json files from scrape-toys")
	out := flag.String("out", "scraped/categories.json", "where to write the results")
	rulesPath := flag.String("rules", "", "rules file to use instead of the embedded rules.json")
	assignAt := flag.Float64(
		"assign",
		0.6,
		"assign a category when its probability is at least this",
	)
	uncertainAt := flag.Float64(
		"uncertain",
		0.4,
		"list a category as uncertain when its probability is at least this",
	)
	minToys := flag.Int(
		"min-toys",
		2,
		"recommend a new category when at least this many toys need it",
	)
	alignAt := flag.Float64(
		"align",
		0.6,
		"treat a candidate as covered when an existing category matches with at least this probability",
	)
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
	existing, current := loadCatalog(rules)
	log.Printf("existing categories: %s", strings.Join(existing, ", "))

	candidates, duplicates, err := alignCandidates(client, rules, existing, *alignAt)
	if err != nil {
		log.Fatal("Error aligning candidates: ", err)
	}
	for c, e := range duplicates {
		log.Printf("candidate %q is already covered by %q", c, e)
	}

	paths, err := filepath.Glob(filepath.Join(*inDir, "toy-*.json"))
	if err != nil {
		log.Fatal(err)
	}
	results := make([]ToyResult, len(paths))
	errs := make([]error, len(paths))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for range *workers {
		wg.Go(func() {
			for i := range jobs {
				results[i], errs[i] = classify(client, rules, paths[i], existing, candidates)
			}
		})
	}
	for i := range paths {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	var toys []ToyResult
	for i, err := range errs {
		if err != nil {
			log.Printf("%s: %v", paths[i], err)
			continue
		}
		toys = append(toys, decide(results[i], current, *assignAt, *uncertainAt))
	}
	sort.Slice(toys, func(i, j int) bool { return toys[i].ToyID < toys[j].ToyID })

	output := Output{
		Existing:    existing,
		Duplicates:  duplicates,
		Suggestions: suggest(toys, candidates, *minToys),
		Toys:        toys,
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatal(err)
	}

	printToys(toys)
	printSuggestions(output.Suggestions)
	log.Printf("classified %d/%d toys → %s", len(toys), len(paths), *out)
}

// loadCatalog returns the real categories in the database and each toy's
// current category value.
func loadCatalog(rules Rules) ([]string, map[int64]string) {
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

	toysStore := store.NewToysStore(db)
	categories, err := toysStore.GetCategories()
	if err != nil {
		log.Fatal("Error reading categories: ", err)
	}
	categories = slices.DeleteFunc(categories, func(c string) bool {
		return c == "" || slices.Contains(rules.Placeholders, c)
	})

	toys, err := toysStore.GetToys()
	if err != nil {
		log.Fatal("Error reading toys: ", err)
	}
	current := make(map[int64]string, len(toys))
	for _, t := range toys {
		current[t.ToyID] = t.Category
	}
	return categories, current
}

func describe(name string, descriptions map[string]string) string {
	return cmp.Or(descriptions[name], name)
}

// alignCandidates asks, for every candidate, which existing category already
// covers the same kind of toy, and returns the candidates that are new.
func alignCandidates(
	client *typesafe.Client,
	rules Rules,
	existing []string,
	alignAt float64,
) (map[string]string, map[string]string, error) {
	existingState := map[string]string{}
	criteria := map[string]any{}
	for _, e := range existing {
		existingState[e] = describe(e, rules.ExistingDescriptions)
		criteria[e] = existingState[e]
	}
	align := rules.AlignQuestion.WithCriteria(criteria)

	names := slices.Sorted(maps.Keys(rules.CandidateCategories))
	questions := map[string]typesafe.Question{}
	for i, name := range names {
		questions[fmt.Sprint(i)] = align.With("candidate", map[string]string{
			"name":        name,
			"description": rules.CandidateCategories[name],
		})
	}

	resp, err := client.Ask(map[string]any{"existing_categories": existingState}, questions)
	if err != nil {
		return nil, nil, err
	}

	candidates, duplicates := map[string]string{}, map[string]string{}
	for i, name := range names {
		a := resp.Answers[fmt.Sprint(i)]
		if a.Choice != "none" && a.Probabilities[a.Choice] >= alignAt {
			duplicates[name] = a.Choice
			continue
		}
		candidates[name] = rules.CandidateCategories[name]
	}
	return candidates, duplicates, nil
}

func classify(
	client *typesafe.Client,
	rules Rules,
	path string,
	existing []string,
	candidates map[string]string,
) (ToyResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ToyResult{}, err
	}
	var toy ToyFile
	if err := json.Unmarshal(raw, &toy); err != nil {
		return ToyResult{}, err
	}

	// Question ids are only for code; each question carries its full meaning.
	questions := map[string]typesafe.Question{}
	for i, name := range existing {
		questions[fmt.Sprintf("existing_%d", i)] = belongsQuestion(
			rules,
			name,
			describe(name, rules.ExistingDescriptions),
		)
	}
	candidateNames := slices.Sorted(maps.Keys(candidates))
	for i, name := range candidateNames {
		questions[fmt.Sprintf("candidate_%d", i)] = belongsQuestion(rules, name, candidates[name])
	}

	resp, err := client.Ask(buildState(rules, toy), questions)
	if err != nil {
		return ToyResult{}, err
	}

	r := ToyResult{
		ToyID:      toy.ToyID,
		ToyName:    toy.ToyName,
		Existing:   map[string]float64{},
		Candidates: map[string]float64{},
	}
	for i, name := range existing {
		r.Existing[name] = resp.Answers[fmt.Sprintf("existing_%d", i)].Noul
	}
	for i, name := range candidateNames {
		r.Candidates[name] = resp.Answers[fmt.Sprintf("candidate_%d", i)].Noul
	}
	return r, nil
}

func belongsQuestion(rules Rules, name, description string) typesafe.Question {
	return rules.BelongsQuestion.With("category", map[string]string{"name": name, "description": description})
}

// buildState keeps the fields that say what kind of toy it is, including
// Amazon's own bestseller categories.
func buildState(rules Rules, toy ToyFile) map[string]any {
	attrs := map[string]any{}
	var amazonCategories []string
	for _, d := range toy.Details {
		for _, k := range rules.RelevantAttributes {
			if v, ok := d.Attributes[k]; ok {
				attrs[k] = v
			}
		}
		if ranks, ok := d.Attributes[rules.RankingAttribute].([]any); ok {
			for _, rank := range ranks {
				s, _ := rank.(string)
				if m := rankingCategory.FindStringSubmatch(strings.TrimSpace(s)); m != nil &&
					!slices.Contains(rules.RootRankings, m[1]) {
					amazonCategories = append(amazonCategories, m[1])
				}
			}
		}
	}
	state := map[string]any{
		"name":       toy.ToyName,
		"title":      toy.Title,
		"features":   toy.Features,
		"attributes": attrs,
	}
	if len(amazonCategories) > 0 {
		state["amazon_categories"] = amazonCategories
	}
	return state
}

// decide applies the thresholds to the raw scores.
func decide(r ToyResult, current map[int64]string, assignAt, uncertainAt float64) ToyResult {
	r.Current = current[r.ToyID]
	for _, name := range slices.Sorted(maps.Keys(r.Existing)) {
		switch p := r.Existing[name]; {
		case p >= assignAt:
			r.Assigned = append(r.Assigned, name)
		case p >= uncertainAt:
			r.Uncertain = append(r.Uncertain, name)
		}
	}
	for _, name := range slices.Sorted(maps.Keys(r.Candidates)) {
		if r.Candidates[name] >= assignAt {
			r.Suggested = append(r.Suggested, name)
		}
	}
	r.Uncovered = len(r.Assigned) == 0
	// Same format the app stores: comma-separated, no spaces.
	r.Proposed = strings.Join(r.Assigned, ",")
	return r
}

// suggest groups the toys by candidate category. A candidate is recommended
// when enough toys need it, or when it is the only fit for a toy that no
// existing category covers.
func suggest(toys []ToyResult, candidates map[string]string, minToys int) []Suggestion {
	byName := map[string]*Suggestion{}
	for _, t := range toys {
		for _, name := range t.Suggested {
			s, ok := byName[name]
			if !ok {
				s = &Suggestion{Category: name, Description: candidates[name]}
				byName[name] = s
			}
			s.Toys = append(s.Toys, t.ToyID)
			if t.Uncovered {
				s.Uncovered = append(s.Uncovered, t.ToyID)
			}
		}
	}

	var suggestions []Suggestion
	for _, s := range byName {
		s.Recommended = len(s.Toys) >= minToys || len(s.Uncovered) > 0
		suggestions = append(suggestions, *s)
	}
	sort.Slice(suggestions, func(i, j int) bool {
		a, b := suggestions[i], suggestions[j]
		if len(a.Uncovered) != len(b.Uncovered) {
			return len(a.Uncovered) > len(b.Uncovered)
		}
		if len(a.Toys) != len(b.Toys) {
			return len(a.Toys) > len(b.Toys)
		}
		return a.Category < b.Category
	})
	return suggestions
}

func printToys(toys []ToyResult) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tACTUAL\tPROPUESTA\tDUDOSAS\tNUEVAS\tTOY")
	for _, t := range toys {
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%s\t%s\t%s\n",
			t.ToyID,
			cmp.Or(t.Current, "-"),
			cmp.Or(t.Proposed, "(ninguna)"),
			scores(
				t.Uncertain,
				t.Existing,
			),
			scores(t.Suggested, t.Candidates),
			truncate(t.ToyName, 45),
		)
	}
	tw.Flush()
}

func printSuggestions(suggestions []Suggestion) {
	if len(suggestions) == 0 {
		return
	}
	fmt.Println()
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NUEVA CATEGORÍA\tTOYS\tSIN CATEGORÍA\tRECOMENDADA")
	for _, s := range suggestions {
		fmt.Fprintf(
			tw,
			"%s\t%s\t%s\t%v\n",
			s.Category,
			ids(s.Toys),
			ids(s.Uncovered),
			s.Recommended,
		)
	}
	tw.Flush()
}

func scores(names []string, probs map[string]float64) string {
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = fmt.Sprintf("%s %.2f", n, probs[n])
	}
	return cmp.Or(strings.Join(parts, ", "), "-")
}

func ids(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = fmt.Sprint(id)
	}
	return cmp.Or(strings.Join(parts, ","), "-")
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}
