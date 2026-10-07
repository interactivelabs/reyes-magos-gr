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
//  4. The same alignment maps database categories missing from
//     existing_categories ("Carros") to the rules category that covers them
//     ("Vehículos"). Per toy, the database one is flagged when it clearly
//     scores higher, so a rule can be added, and dropped otherwise; when
//     both fit about equally well, the drop is flagged too.
//  5. Candidates that enough toys need with high confidence are promoted:
//     they move from candidate_categories to existing_categories in the
//     rules file and are assigned like any existing category.
package main

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"math"
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

// defaultRulesPath is where promotions are saved when -rules isn't given,
// relative to the repository root like the usage above.
const defaultRulesPath = "scripts/classify-toy-categories/rules.json"

// Rules is the content of rules.json.
type Rules struct {
	// ExistingDescriptions explains the categories already in the database,
	// so the model doesn't have to guess from terse (or misspelled) names.
	// The keys must match the database exactly because the catalog filters
	// with LIKE. Categories missing here are sent with their name only, and
	// keys missing from the database are promoted candidates still waiting
	// for their toys to be updated.
	ExistingDescriptions map[string]string `json:"existing_categories"`
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
	if err := json.Unmarshal(data, &r); err != nil {
		return Rules{}, err
	}
	if r.ExistingDescriptions == nil {
		r.ExistingDescriptions = map[string]string{}
	}
	return r, nil
}

func saveRules(path string, rules Rules) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rules); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
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
	ToyID     int64    `json:"toy_id"`
	ToyName   string   `json:"toy_name"`
	Current   string   `json:"current_category"`
	Proposed  string   `json:"proposed_category"`
	Assigned  []string `json:"assigned"`
	Uncertain []string `json:"uncertain,omitempty"`
	Suggested []string `json:"suggested_new,omitempty"`
	Uncovered bool     `json:"uncovered"`
	// Outscoring maps an unruled database category to the rules category
	// that covers the same toys but scored lower for this toy: a hint that
	// the database one may deserve an existing_categories entry.
	Outscoring map[string]string `json:"unruled_outscoring,omitempty"`
	// Similar maps an unruled database category to the rules category that
	// covers the same toys when both fit this toy about equally well. The
	// rules one is kept, but the pair is worth a look.
	Similar    map[string]string `json:"unruled_similar,omitempty"`
	Existing   Scores            `json:"existing_scores"`
	Candidates Scores            `json:"candidate_scores"`
}

// Scores maps a category to its probability. It is written highest first, so
// the best fits lead each toy in the output.
type Scores map[string]float64

func (s Scores) MarshalJSON() ([]byte, error) {
	names := slices.Collect(maps.Keys(s))
	slices.SortFunc(names, func(a, b string) int {
		return cmp.Or(cmp.Compare(s[b], s[a]), strings.Compare(a, b))
	})
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, name := range names {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, err := json.Marshal(name)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(s[name])
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(value)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
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
	Unruled     map[string]string `json:"unruled_already_covered"`
	Promoted    []string          `json:"promoted,omitempty"`
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
	promoteAt := flag.Float64(
		"promote",
		0.75,
		"promote a candidate to an existing category when at least -min-toys toys belong to it with this probability (0 disables)",
	)
	outscoreAt := flag.Float64(
		"outscore",
		0.5,
		"flag a database category without a rule only when it scores above this",
	)
	outscoreBy := flag.Float64(
		"outscore-by",
		0.1,
		"flag a database category without a rule only when it beats the rules category covering it by more than this",
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

	unruled, err := alignUnruled(client, rules, existing, *alignAt)
	if err != nil {
		log.Fatal("Error aligning database categories: ", err)
	}
	for u, e := range unruled {
		log.Printf("database category %q is covered by %q", u, e)
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

	var raw []ToyResult
	for i, err := range errs {
		if err != nil {
			log.Printf("%s: %v", paths[i], err)
			continue
		}
		raw = append(raw, results[i])
	}
	sort.Slice(raw, func(i, j int) bool { return raw[i].ToyID < raw[j].ToyID })
	decideAll := func() []ToyResult {
		toys := make([]ToyResult, len(raw))
		for i, r := range raw {
			toys[i] = decide(r, current, unruled, *assignAt, *uncertainAt, *outscoreAt, *outscoreBy)
		}
		return toys
	}
	toys := decideAll()

	var promoted []string
	if *promoteAt > 0 {
		promoted = promote(toys, candidates, *promoteAt, *minToys)
	}
	if len(promoted) > 0 {
		for _, name := range promoted {
			rules.ExistingDescriptions[name] = candidates[name]
			delete(rules.CandidateCategories, name)
			delete(candidates, name)
			existing = append(existing, name)
			// The belongs question is the same for both kinds, so the
			// candidate scores count as existing ones.
			for _, r := range raw {
				r.Existing[name] = r.Candidates[name]
				delete(r.Candidates, name)
			}
		}
		toys = decideAll()
		path := cmp.Or(*rulesPath, defaultRulesPath)
		if err := saveRules(path, rules); err != nil {
			log.Fatal("Error saving rules: ", err)
		}
		log.Printf("promoted %s → %s", strings.Join(promoted, ", "), path)
	}

	output := Output{
		Existing:    existing,
		Duplicates:  duplicates,
		Unruled:     unruled,
		Promoted:    promoted,
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

	for _, t := range toys {
		for _, u := range slices.Sorted(maps.Keys(t.Outscoring)) {
			name := t.Outscoring[u]
			log.Printf(
				"toy %d: %q (%.2f) scores below %q (%.2f), which has no existing_categories entry",
				t.ToyID, name, t.Existing[name], u, t.Existing[u],
			)
		}
		for _, u := range slices.Sorted(maps.Keys(t.Similar)) {
			name := t.Similar[u]
			log.Printf(
				"toy %d: %q (%.2f) and %q fit about equally well; kept %q, %q has no existing_categories entry",
				t.ToyID, name, t.Existing[name], u, name, u,
			)
		}
	}
	printToys(toys)
	printSuggestions(output.Suggestions)
	log.Printf("classified %d/%d toys → %s", len(toys), len(paths), *out)
}

// loadCatalog returns the real categories in the database, plus promoted ones
// not stored yet, and each toy's current category value.
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
	for _, name := range slices.Sorted(maps.Keys(rules.ExistingDescriptions)) {
		if !slices.Contains(categories, name) {
			categories = append(categories, name)
		}
	}

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
	duplicates, err := align(client, rules, existing, rules.CandidateCategories, alignAt)
	if err != nil {
		return nil, nil, err
	}
	candidates := map[string]string{}
	for name, description := range rules.CandidateCategories {
		if _, ok := duplicates[name]; !ok {
			candidates[name] = description
		}
	}
	return candidates, duplicates, nil
}

// alignUnruled maps the database categories missing from existing_categories
// to the rules category that covers the same kind of toy, if any.
func alignUnruled(
	client *typesafe.Client,
	rules Rules,
	existing []string,
	alignAt float64,
) (map[string]string, error) {
	unruled := map[string]string{}
	for _, e := range existing {
		if _, ok := rules.ExistingDescriptions[e]; !ok {
			unruled[e] = e
		}
	}
	if len(unruled) == 0 {
		return map[string]string{}, nil
	}
	return align(
		client,
		rules,
		slices.Sorted(maps.Keys(rules.ExistingDescriptions)),
		unruled,
		alignAt,
	)
}

// align asks, for every item, which of the targets already groups the same
// kind of toys, and returns the items that matched with their target.
func align(
	client *typesafe.Client,
	rules Rules,
	targets []string,
	items map[string]string,
	alignAt float64,
) (map[string]string, error) {
	targetState := map[string]string{}
	criteria := map[string]any{}
	for _, t := range targets {
		targetState[t] = describe(t, rules.ExistingDescriptions)
		criteria[t] = targetState[t]
	}
	question := rules.AlignQuestion.WithCriteria(criteria)

	names := slices.Sorted(maps.Keys(items))
	questions := map[string]typesafe.Question{}
	for i, name := range names {
		questions[fmt.Sprint(i)] = question.With("candidate", map[string]string{
			"name":        name,
			"description": items[name],
		})
	}

	resp, err := client.Ask(map[string]any{"existing_categories": targetState}, questions)
	if err != nil {
		return nil, err
	}

	matches := map[string]string{}
	for i, name := range names {
		a := resp.Answers[fmt.Sprint(i)]
		if a.Choice != "none" && a.Probabilities[a.Choice] >= alignAt {
			matches[name] = a.Choice
		}
	}
	return matches, nil
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
		Existing:   Scores{},
		Candidates: Scores{},
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
	return rules.BelongsQuestion.With(
		"category",
		map[string]string{"name": name, "description": description},
	)
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

// decide resolves the unruled database categories against the rules ones
// that cover them and applies the thresholds to the raw scores. An unruled
// category is kept and flagged only when it clearly beats its rules category:
// above outscoreAt and by more than outscoreBy. Otherwise the rules category
// wins, and when the unruled one is still above outscoreAt and within
// outscoreBy of it, the pair is flagged as similar.
func decide(
	r ToyResult,
	current map[int64]string,
	unruled map[string]string,
	assignAt, uncertainAt, outscoreAt, outscoreBy float64,
) ToyResult {
	r.Current = current[r.ToyID]
	// The raw scores are shared between decide calls, so drop from a copy.
	r.Existing = maps.Clone(r.Existing)
	for u, name := range unruled {
		p, diff := r.Existing[u], r.Existing[u]-r.Existing[name]
		if p <= outscoreAt || diff <= outscoreBy {
			if p > outscoreAt && math.Abs(diff) <= outscoreBy {
				if r.Similar == nil {
					r.Similar = map[string]string{}
				}
				r.Similar[u] = name
			}
			delete(r.Existing, u)
			continue
		}
		if r.Outscoring == nil {
			r.Outscoring = map[string]string{}
		}
		r.Outscoring[u] = name
	}
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

// promote returns the candidates that at least minToys toys belong to with
// probability promoteAt or more.
func promote(
	toys []ToyResult,
	candidates map[string]string,
	promoteAt float64,
	minToys int,
) []string {
	var promoted []string
	for _, name := range slices.Sorted(maps.Keys(candidates)) {
		sure := 0
		for _, t := range toys {
			if t.Candidates[name] >= promoteAt {
				sure++
			}
		}
		if sure >= minToys {
			promoted = append(promoted, name)
		}
	}
	return promoted
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
