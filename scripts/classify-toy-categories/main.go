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
// Jev picks from options, it doesn't invent names, so new categories come
// from candidateCategories below:
//  1. One request aligns every candidate with the existing categories and
//     drops those the catalog already covers ("Vehículos" vs "Carros").
//  2. One request per toy asks a yes/no question per existing category and
//     per remaining candidate; a toy can belong to several.
//  3. Code assigns existing categories over the threshold and suggests a
//     candidate when enough toys need it.
package main

import (
	"cmp"
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

// existingDescriptions explains the categories already in the database, so
// the model doesn't have to guess from terse (or misspelled) names. The keys
// must match the database exactly because the catalog filters with LIKE.
// Categories missing here are sent with their name only.
var existingDescriptions = map[string]string{
	"Arte":         "Art and crafts: drawing, painting, coloring, tracing projectors, art supply sets.",
	"Carros":       "Toy cars, trucks and other vehicles, including vehicle sets and transporters.",
	"Construction": "Building and assembling (blocks, magnetic kits, take-apart toys) or construction-site vehicles such as dump trucks and excavators.",
	"Dinosaurio":   "Dinosaur-themed toys of any kind.",
	"Electronico":  "Toys powered by batteries or USB charging whose play depends on electronics: lights, sounds, screens, remote control, cameras.",
	"Maquillahe":   "Kids' makeup, nail polish, cosmetics, beauty and vanity sets.",
	"Muneca":       "Dolls, plush dolls, dollhouses and doll styling heads.",
	"Musica":       "Musical instruments, karaoke and other toys made for making music.",
	"Video Juegos": "Video games, consoles and handheld devices whose main use includes playing digital games.",
}

// placeholders are values stored in the category column that aren't real
// categories.
var placeholders = []string{"Nuevo"}

// candidateCategories are names the script may suggest. Add freely: the
// alignment step drops any that an existing category already covers.
var candidateCategories = map[string]string{
	"Deportes":            "Sports and active play: balls, goals, hoops, rackets, hockey, sports games for kids.",
	"Aire Libre":          "Outdoor play: tents for the garden, water toys, bubbles, sandbox toys, ride-ons for outside.",
	"Peluches":            "Plush and stuffed animals or characters made to cuddle.",
	"Juegos de Mesa":      "Board games, card games and tabletop games with rules for several players.",
	"Ciencia":             "Science, experiments, coding and STEM kits that teach how things work.",
	"Bebés":               "Toys made for babies and toddlers under 2: rattles, activity toys, baby laptops.",
	"Disfraces":           "Costumes, dress-up clothes and role play accessories like crowns and wands.",
	"Juego de Rol":        "Pretend play sets that imitate adult life: kitchens, tools, doctor kits, shops.",
	"Casitas y Carpas":    "Play tents, castles, play houses and structures children get inside.",
	"Figuras de Acción":   "Action figures and character figurines from shows, movies or games.",
	"Rompecabezas":        "Puzzles and brain teasers.",
	"Lanzadores":          "Foam dart blasters, targets and other aiming or shooting toys.",
	"Cámaras":             "Kids' cameras and video recorders for taking photos and videos.",
	"Animales":            "Toys about animals other than dinosaurs: farm, ocean, fishing games, pets.",
	"Princesas":           "Princess-themed toys and characters such as Disney princesses.",
	"Vehículos":           "Toy cars, trucks and other vehicles.",
	"Manualidades":        "Crafts and art activities.",
	"Belleza":             "Makeup, nail polish and beauty sets.",
	"Bloques":             "Building blocks and construction sets.",
	"Instrumentos":        "Musical instruments for kids.",
	"Teléfonos y Tablets": "Kids' phones, tablets and smartwatches with games and apps.",
}

// relevantAttrs are the scraped attributes that say what kind of toy it is.
var relevantAttrs = []string{
	"Nombre Tipo Artículo", "Componentes Incluidos", "Objetivo educativo",
	"Características especiales", "Tema", "Personaje", "Material",
	"Tipo de fuente de alimentación", "¿Requiere pilas?", "Descripción del rango de edad",
}

const rankingAttr = "Clasificación en los más vendidos de Amazon"

// rootRankings are Amazon's top-level departments, too broad to say anything.
var rootRankings = []string{"Juguetes y Juegos", "Hogar y Cocina"}

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
	assignAt := flag.Float64("assign", 0.6, "assign a category when its probability is at least this")
	uncertainAt := flag.Float64("uncertain", 0.4, "list a category as uncertain when its probability is at least this")
	minToys := flag.Int("min-toys", 2, "recommend a new category when at least this many toys need it")
	alignAt := flag.Float64("align", 0.6, "treat a candidate as covered when an existing category matches with at least this probability")
	workers := flag.Int("workers", 4, "concurrent requests")
	flag.Parse()

	client, err := typesafe.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	existing, current := loadCatalog()
	log.Printf("existing categories: %s", strings.Join(existing, ", "))

	candidates, duplicates, err := alignCandidates(client, existing, *alignAt)
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
				results[i], errs[i] = classify(client, paths[i], existing, candidates)
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
func loadCatalog() ([]string, map[int64]string) {
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
		return c == "" || slices.Contains(placeholders, c)
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
func alignCandidates(client *typesafe.Client, existing []string, alignAt float64) (map[string]string, map[string]string, error) {
	criteria := map[string]any{
		"none": "No existing category covers this kind of toy; it would be a genuinely new category.",
	}
	existingState := map[string]string{}
	for _, e := range existing {
		criteria[e] = describe(e, existingDescriptions)
		existingState[e] = describe(e, existingDescriptions)
	}

	names := slices.Sorted(maps.Keys(candidateCategories))
	questions := map[string]any{}
	for i, name := range names {
		questions[fmt.Sprint(i)] = map[string]any{
			"type": "choice",
			"instructions": map[string]any{
				"candidate": map[string]string{"name": name, "description": candidateCategories[name]},
				"question": "Which category in `existing_categories` already groups the same kind of toys as `candidate`, " +
					"so a shopper looking for `candidate` toys would find them there? " +
					"Pick one only if it covers most of the same toys, not if it merely overlaps with a few.",
			},
			"criteria": criteria,
		}
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
		candidates[name] = candidateCategories[name]
	}
	return candidates, duplicates, nil
}

func classify(client *typesafe.Client, path string, existing []string, candidates map[string]string) (ToyResult, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ToyResult{}, err
	}
	var toy ToyFile
	if err := json.Unmarshal(raw, &toy); err != nil {
		return ToyResult{}, err
	}

	// Question ids are only for code; each question carries its full meaning.
	questions := map[string]any{}
	for i, name := range existing {
		questions[fmt.Sprintf("existing_%d", i)] = belongsQuestion(name, describe(name, existingDescriptions))
	}
	candidateNames := slices.Sorted(maps.Keys(candidates))
	for i, name := range candidateNames {
		questions[fmt.Sprintf("candidate_%d", i)] = belongsQuestion(name, candidates[name])
	}

	resp, err := client.Ask(buildState(toy), questions)
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

func belongsQuestion(name, description string) map[string]any {
	return map[string]any{
		"type": "noul",
		"instructions": map[string]any{
			"category": map[string]string{"name": name, "description": description},
			"question": "Would a parent browsing a toy catalog expect to find this toy under `category`?",
		},
		"criteria": map[string]string{
			"true":  "What the toy is, its theme or its main way of playing matches `category.description`.",
			"false": "The toy only shares a minor detail with `category` (e.g. a small light on a toy that isn't electronic), or nothing at all.",
		},
	}
}

// buildState keeps the fields that say what kind of toy it is, including
// Amazon's own bestseller categories.
func buildState(toy ToyFile) map[string]any {
	attrs := map[string]any{}
	var amazonCategories []string
	for _, d := range toy.Details {
		for _, k := range relevantAttrs {
			if v, ok := d.Attributes[k]; ok {
				attrs[k] = v
			}
		}
		if ranks, ok := d.Attributes[rankingAttr].([]any); ok {
			for _, rank := range ranks {
				s, _ := rank.(string)
				if m := rankingCategory.FindStringSubmatch(strings.TrimSpace(s)); m != nil && !slices.Contains(rootRankings, m[1]) {
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
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n",
			t.ToyID, cmp.Or(t.Current, "-"), cmp.Or(t.Proposed, "(ninguna)"),
			scores(t.Uncertain, t.Existing), scores(t.Suggested, t.Candidates), truncate(t.ToyName, 45))
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
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", s.Category, ids(s.Toys), ids(s.Uncovered), s.Recommended)
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
