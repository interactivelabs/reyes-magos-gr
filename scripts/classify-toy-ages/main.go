// Command classify-toy-ages reads the JSON files written by scrape-toys and
// asks TypeSafe (Jev) for a realistic age range for each toy, since the
// manufacturer's "Edad recomendada" on Amazon is often wrong ("100 años y
// más", "5 meses y más" for a remote-control dinosaur...).
//
// Usage:
//
//	TYPESAFE_API_KEY=... go run ./scripts/classify-toy-ages -in ./scraped
//
// The model only picks semantic options (developmental stages, yes/no);
// turning those into age_min/age_max numbers and applying the safety rules
// happens here in code.
package main

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"reyes-magos-gr/scripts/internal/typesafe"
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

// stage is a Choice option: the model sees the key and description, code
// uses the age.
type stage struct {
	key  string
	age  int64
	desc any
}

// youngestStages, like oldestStages, describe each floor by the skills the
// toy demands and the kind of toy, so options differ in meaning, not number.
var youngestStages = []stage{
	{"newborn", 0, map[string]string{
		"suitable_from":   "Birth, babies under 12 months.",
		"skills_needed":   "None beyond looking, grasping, shaking or mouthing.",
		"typical_toys":    "Rattles, teethers, soft cloth books, crib mobiles, tummy-time mats, baby toys labelled 'desde 0 meses' or 'recién nacido'.",
		"too_hard_before": "Nothing: it is built for newborns and safe to put in the mouth.",
	}},
	{"one_year", 1, map[string]string{
		"suitable_from":   "About 1 year, once a baby sits steadily, crawls or starts walking.",
		"skills_needed":   "Pressing big buttons, pushing, pulling, banging; no pretend play or rules.",
		"typical_toys":    "Baby laptops and phones with lights and songs, activity cubes, push walkers, stacking rings, chunky blocks, 'mi primer' toys for 6 to 12 months and up.",
		"too_hard_before": "A baby who cannot yet sit up or grab deliberately gets nothing from it.",
	}},
	{"two_years", 2, map[string]string{
		"suitable_from":   "About 2 years, a walking toddler.",
		"skills_needed":   "Simple imitation and pretend play, carrying and pushing toys around, working buckles or big knobs; still mouths things.",
		"typical_toys":    "Plush dolls with Montessori buckles, toy wagons and big push vehicles, chunky puzzles, toddler ride-ons, simple shakers and drums.",
		"too_hard_before": "A 1 year old cannot yet use it on purpose or walk with it.",
	}},
	{"three_years", 3, map[string]string{
		"suitable_from":   "About 3 years, a preschooler who no longer puts toys in the mouth.",
		"skills_needed":   "Role play with small accessories, taking turns, simple on/off electronics, assembling a few pieces with an adult.",
		"typical_toys":    "Play tents and castles, toy vehicle sets with small cars, styling heads, dinosaur sets, magnetic fishing games, kids' cameras and phones with simple menus, karaoke microphones.",
		"too_hard_before": "It has small loose pieces or needs pretend play and hand control that a 2 year old does not have yet.",
	}},
	{"four_years", 4, map[string]string{
		"suitable_from":   "About 4 to 5 years.",
		"skills_needed":   "Following simple rules, aiming, steering with a remote control, careful use of cosmetics or art materials.",
		"typical_toys":    "Washable kids' makeup kits, remote-control animals, sports games with goals and scoring, toy drum kits with stool and sticks.",
		"too_hard_before": "A 3 year old cannot follow its rules, aim, or use the materials without making a mess or misusing them.",
	}},
	{"six_years", 6, map[string]string{
		"suitable_from":   "About 6 years, a school-age child.",
		"skills_needed":   "Multi-step building, reading a little, shooting foam darts safely, using a screwdriver tool.",
		"typical_toys":    "Magnetic rod-and-ball building kits, take-apart dinosaur kits with screwdriver, foam dart blasters and electronic targets, vanity sets with many small cosmetics.",
		"too_hard_before": "Younger children cannot follow the steps or are at risk with the small or projectile parts.",
	}},
	{"eight_years", 8, map[string]string{
		"suitable_from":   "About 8 years.",
		"skills_needed":   "Reading instructions independently, fine motor precision, patience for detailed projects.",
		"typical_toys":    "Detailed construction models, strategy board games, science kits, art sets with acrylic paint, brushes and many media.",
		"too_hard_before": "A 6 or 7 year old would need an adult to do most of it.",
	}},
	{"ten_years", 10, map[string]string{
		"suitable_from":   "About 10 years, a preteen.",
		"skills_needed":   "Using real tools and materials responsibly with little supervision.",
		"typical_toys":    "Hobby-grade gadgets, advanced electronics or coding kits, professional-grade art sets.",
		"too_hard_before": "Younger children would misuse or be frustrated by real tools and materials.",
	}},
	{"teen", 13, map[string]string{
		"suitable_from":   "13 years and older.",
		"skills_needed":   "Adult-level judgment.",
		"typical_toys":    "Products meant for teens and adults, not designed for children at all.",
		"too_hard_before": "It is unsafe or inappropriate for children under 13.",
	}},
}

// oldestStages describe each ceiling by the kind of toy and why older kids
// drop it, rather than by numbers alone: Jev tells "5 to 6" from "7 to 8"
// poorly, but tells a push truck from a karaoke machine well.
var oldestStages = []stage{
	{"toddlers", 2, map[string]string{
		"designed_for":  "Babies and toddlers up to about 2 years.",
		"typical_toys":  "Rattles, teethers, baby laptops and phones with big buttons, activity cubes, crib and stroller toys, first-words toys, anything named 'bebé' or 'mi primer'.",
		"outgrown_when": "A 3 year old already finds it babyish: it only offers lights, sounds and pressing big buttons.",
	}},
	{"preschool_early", 3, map[string]string{
		"designed_for":  "Toddlers and young preschoolers up to about 3 years.",
		"typical_toys":  "Soft plush dolls with Montessori buckles and zippers, chunky shape sorters, large simple puzzles, pull-along toys, ride-ons for toddlers.",
		"outgrown_when": "Once a child can do pretend play with storylines and small accessories, these feel too simple.",
	}},
	{"preschool", 5, map[string]string{
		"designed_for":  "Preschoolers, about 3 to 5 years.",
		"typical_toys":  "Push vehicles and construction trucks without remote control, simple percussion and musical instrument sets, magnetic fishing games, simple dinosaur and car sets, toy wagons.",
		"outgrown_when": "Children starting primary school want rules, challenge, electronics or more realistic detail; the toy has a single simple action repeated over and over.",
	}},
	{"early_school", 7, map[string]string{
		"designed_for":  "Kindergarten and first school years, about 5 to 7 years.",
		"typical_toys":  "Princess play tents and castles, dress-up sets, styling heads, simple remote-control animals, toy drum kits, drawing and tracing projectors, simple sports games for kids.",
		"outgrown_when": "Pretend play is centered on characters or make-believe that 8 year olds consider childish.",
	}},
	{"middle_childhood", 10, map[string]string{
		"designed_for":  "School-age children, about 6 to 10 years.",
		"typical_toys":  "Large dollhouses with furniture and many accessories, magnetic and take-apart building kits, kids' basketball hoops with scoreboards, karaoke microphones, kids' makeup and vanity sets.",
		"outgrown_when": "Preteens prefer real-looking gadgets or hobby gear; this still looks like a children's toy.",
	}},
	{"preteen", 12, map[string]string{
		"designed_for":  "Older children and preteens, up to about 12 years.",
		"typical_toys":  "Kids' digital cameras and phones with real functions, foam dart blasters and electronic targets, detailed construction kits, games that require reading and strategy.",
		"outgrown_when": "Teenagers want the real adult version instead of a product branded for kids.",
	}},
	{"teen", 16, map[string]string{
		"designed_for":  "Teenagers as well as children, 13 years and older.",
		"typical_toys":  "Professional art supply sets (acrylics, watercolors, markers), real sports equipment, hobby-grade gadgets that adults also use.",
		"outgrown_when": "Not outgrown in childhood: it is a real tool, not a toy for kids.",
	}},
}

// noisyAttrs carry no information about who the toy is for.
var noisyAttrs = []string{"ASIN", "UPC", "EAN", "Núm. de identificación", "Número Modelo", "Número de pieza", "Número de paquetes", "Número de productos", "Clasificación", "Opinión", "Fabricante", "Marca", "marca", "Modelo", "embalaje", "Total del paquete", "Peso", "Dimensiones", "Medidas"}

func main() {
	inDir := flag.String("in", "scraped", "directory with the toy-*.json files from scrape-toys")
	out := flag.String("out", "scraped/ages.json", "where to write the results")
	minConfidence := flag.Float64("min-confidence", 0.5, "flag toys whose youngest/oldest confidence is below this for review")
	workers := flag.Int("workers", 4, "concurrent requests")
	flag.Parse()

	client, err := typesafe.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

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
				results[i], errs[i] = classify(client, paths[i], *minConfidence)
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
		ok = append(ok, results[i])
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
	fmt.Fprintln(tw, "ID\tRANGE\tMIN CONF\tMAX CONF\tSMALL PARTS\tFABRICANTE\tFAB OK\tREVIEW\tTOY")
	for _, r := range ok {
		manufacturerOK := "-"
		if r.ManufacturerOK != nil {
			manufacturerOK = fmt.Sprintf("%.2f", *r.ManufacturerOK)
		}
		fmt.Fprintf(tw, "%d\t%d-%d\t%.2f\t%.2f\t%.2f\t%s\t%s\t%s\t%s\n",
			r.ToyID, r.AgeMin, r.AgeMax, r.YoungestConfidence, r.OldestConfidence, r.SmallParts,
			cmp.Or(r.ManufacturerAge, "-"), manufacturerOK, strings.Join(r.ReviewReasons, "; "), truncate(r.ToyName, 50))
	}
	tw.Flush()
	log.Printf("classified %d/%d toys → %s", len(ok), len(paths), *out)
}

func classify(client *typesafe.Client, path string, minConfidence float64) (Result, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Result{}, err
	}
	var toy ToyFile
	if err := json.Unmarshal(raw, &toy); err != nil {
		return Result{}, err
	}

	state, manufacturerAge := buildState(toy)
	questions := map[string]any{
		"youngest": map[string]any{
			"type": "choice",
			"instructions": map[string]any{
				"question": "What is the youngest age group that can safely play with this toy and get its intended play out of it, judging by what the toy is and the skills it demands (`title`, `features`, `attributes`)?",
				"evidence_rules": []string{
					"Match the toy to the option whose `skills_needed` and `typical_toys` describe it best.",
					"A lower age written for children in `title` or `features` (e.g. \"desde 6 meses\", \"niños de 3 a 12 años\") is good evidence for the floor, unless the toy clearly demands more skill than that.",
					"`manufacturer_age` is often wrong (e.g. \"5 meses y más\" on a remote-control toy); ignore it when it contradicts the kind of toy.",
					"Small loose pieces, magnets, darts and cosmetics need a child who no longer puts things in the mouth; big buttons, lights and songs point to babies.",
				},
			},
			"criteria": stageCriteria(youngestStages),
		},
		"oldest": map[string]any{
			"type": "choice",
			"instructions": map[string]any{
				"question": "What is the oldest age group this toy is designed for, judging by what the toy is and the kind of play it offers (`title`, `features`, `attributes`)?",
				"evidence_rules": []string{
					"Match the toy to the option whose `typical_toys` and `outgrown_when` describe it best.",
					"Sellers stretch the upper age in titles for marketing (e.g. a toddler tracing toy sold \"para niños de 3 a 12 años\"); trust the kind of play over the stated number.",
					"`manufacturer_age` is often wrong (e.g. \"100 años y más\"); ignore it when it contradicts the kind of toy.",
					"Words like \"niños pequeños\", \"bebé\" or \"preescolar\" point to younger groups; real functions (camera, Bluetooth, scoreboard) and many small detailed pieces point to older groups.",
				},
			},
			"criteria": stageCriteria(oldestStages),
		},
		"small_parts": map[string]any{
			"type":         "noul",
			"instructions": "Does this toy come with loose small pieces that a child under 3 could swallow or choke on, such as small cars, beads, magnets, small balls, darts or cosmetics?",
			"criteria": map[string]string{
				"true":  "At least one loose item included in the box is small enough to fit in a toddler's mouth, or is a cosmetic a toddler must not eat.",
				"false": "Every piece is large, soft or permanently attached. Batteries inside a closed compartment and buttons built into the toy do not count.",
			},
		},
	}
	if manufacturerAge != "" {
		questions["manufacturer_plausible"] = map[string]any{
			"type":         "noul",
			"instructions": "Is `manufacturer_age` a realistic description of the ages of children this toy is designed for?",
			"criteria": map[string]string{
				"true":  "The stated age matches the kind of toy, its complexity and its safety.",
				"false": "The stated age is absurd (e.g. 60 or 100 years for a children's toy) or clearly too young or too old for this toy.",
			},
		}
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
		AgeMin:             stageAge(youngestStages, y.Choice),
		AgeMax:             stageAge(oldestStages, o.Choice),
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
	if r.YoungestConfidence < minConfidence && adjacentMass(youngestStages, y.Probabilities) < 0.8 {
		r.review(fmt.Sprintf("uncertain youngest (%.2f)", r.YoungestConfidence))
	}
	if r.OldestConfidence < minConfidence && adjacentMass(oldestStages, o.Probabilities) < 0.8 {
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
func buildState(toy ToyFile) (map[string]any, string) {
	attrs := map[string]any{}
	var manufacturerAge string
	for _, d := range toy.Details {
		for k, v := range d.Attributes {
			if k == "Edad recomendada por el fabricante" {
				manufacturerAge, _ = v.(string)
				continue
			}
			if !isNoisy(k) {
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

func isNoisy(key string) bool {
	for _, n := range noisyAttrs {
		if strings.Contains(key, n) {
			return true
		}
	}
	return false
}

func stageCriteria(stages []stage) map[string]any {
	c := make(map[string]any, len(stages))
	for _, s := range stages {
		c[s.key] = s.desc
	}
	return c
}

// adjacentMass is the highest probability held by two neighbouring stages.
func adjacentMass(stages []stage, probs map[string]float64) float64 {
	var best float64
	for i := 1; i < len(stages); i++ {
		best = max(best, probs[stages[i-1].key]+probs[stages[i].key])
	}
	return best
}

func stageAge(stages []stage, key string) int64 {
	for _, s := range stages {
		if s.key == key {
			return s.age
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
