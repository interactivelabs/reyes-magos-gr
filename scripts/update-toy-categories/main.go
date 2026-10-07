// Command update-toy-categories writes the categories proposed by
// classify-toy-categories (scraped/categories.json) to the toys table.
//
// Usage:
//
//	go run ./scripts/update-toy-categories               # dev, from TURSO_*
//	go run ./scripts/update-toy-categories -production   # prod, from PROD_TURSO_*
//
// By default it uses the same TURSO_* environment variables as the app. With
// -production every one of them must be overridden by its PROD_ counterpart
// (PROD_TURSO_DATABASE_URL, ...), so prod is never reached by accident with
// the dev settings.
//
// Each toy's category is replaced, not added to, with its assigned
// categories. When classify-toy-categories flagged a database category
// without a rule against the rules category covering it (unruled_similar or
// unruled_outscoring, "Carros" vs "Vehículos"), the choice is prompted. The
// uncertain categories are offered too and only added when picked. Toys left
// without any category are skipped rather than cleared. Nothing is written
// until the final list of changes is confirmed.
package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"maps"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
)

var tursoEnv = []string{
	"TURSO_DATABASE_URL",
	"TURSO_AUTH_TOKEN",
	"TURSO_ENCRYPTION_KEY",
	"TURSO_LOCAL_DB",
}

// Result is the part of classify-toy-categories' output this script needs.
type Result struct {
	ToyID      int64              `json:"toy_id"`
	ToyName    string             `json:"toy_name"`
	Assigned   []string           `json:"assigned"`
	Uncertain  []string           `json:"uncertain,omitempty"`
	Outscoring map[string]string  `json:"unruled_outscoring,omitempty"`
	Similar    map[string]string  `json:"unruled_similar,omitempty"`
	Existing   map[string]float64 `json:"existing_scores"`
}

type Output struct {
	Toys []Result `json:"toys"`
}

type Update struct {
	Toy        models.Toy
	Categories []string
	Reviewed   bool
}

func main() {
	in := flag.String("in", "scraped/categories.json", "categories written by classify-toy-categories")
	production := flag.Bool("production", false, "write to prod, reading the TURSO_* settings from PROD_TURSO_*")
	flag.Parse()

	if *production {
		useProductionEnv()
	}
	target := "dev"
	if *production {
		target = "PRODUCTION"
	}
	if u, err := url.Parse(os.Getenv("TURSO_DATABASE_URL")); err == nil && u.Host != "" {
		target += " (" + u.Host + ")"
	}

	results, err := loadResults(*in)
	if err != nil {
		log.Fatal("Error reading ", *in, ": ", err)
	}

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
	log.Printf("database: %s", target)

	toysStore := store.NewToysStore(db)
	toys, err := toysStore.GetToys()
	if err != nil {
		log.Fatal("Error reading toys: ", err)
	}
	current := make(map[int64]models.Toy, len(toys))
	for _, t := range toys {
		current[t.ToyID] = t
	}

	input := bufio.NewScanner(os.Stdin)
	var updates []Update
	for _, r := range results {
		toy, ok := current[r.ToyID]
		if !ok {
			log.Printf("toy %d is not in the database, skipping", r.ToyID)
			continue
		}
		u := Update{Toy: toy, Categories: slices.Clone(r.Assigned)}
		if len(r.Similar) > 0 || len(r.Outscoring) > 0 || len(r.Uncertain) > 0 {
			present(r, toy)
			var keep bool
			if u.Categories, keep = review(input, r); !keep {
				fmt.Println("  skipped")
				continue
			}
			u.Reviewed = true
		}
		if len(u.Categories) == 0 {
			log.Printf("toy %d has no category to assign, skipping", r.ToyID)
			continue
		}
		slices.Sort(u.Categories)
		u.Categories = slices.Compact(u.Categories)
		if !slices.Equal(u.Categories, split(toy.Category)) {
			updates = append(updates, u)
		}
	}

	if len(updates) == 0 {
		log.Print("nothing to update")
		return
	}
	printUpdates(updates)
	if !confirm(input, fmt.Sprintf("\nWrite %d updates to %s? [y/N]: ", len(updates), target)) {
		log.Print("nothing written")
		return
	}

	failed := 0
	for _, u := range updates {
		// Same format the app stores: comma-separated, no spaces.
		if err := toysStore.UpdateToyCategory(u.Toy.ToyID, strings.Join(u.Categories, ",")); err != nil {
			log.Printf("toy %d: %v", u.Toy.ToyID, err)
			failed++
		}
	}
	log.Printf("updated %d/%d toys in %s", len(updates)-failed, len(updates), target)
	if failed > 0 {
		os.Exit(1)
	}
}

// useProductionEnv replaces every TURSO_* variable with its PROD_ value. All
// of them are required, and the URL must differ from the dev one.
func useProductionEnv() {
	var missing []string
	for _, k := range tursoEnv {
		if os.Getenv("PROD_"+k) == "" {
			missing = append(missing, "PROD_"+k)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("-production requires %s", strings.Join(missing, ", "))
	}
	if os.Getenv("PROD_TURSO_DATABASE_URL") == os.Getenv("TURSO_DATABASE_URL") {
		log.Fatal("PROD_TURSO_DATABASE_URL is the same as TURSO_DATABASE_URL")
	}
	for _, k := range tursoEnv {
		os.Setenv(k, os.Getenv("PROD_"+k))
	}
}

func loadResults(path string) ([]Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out Output
	return out.Toys, json.Unmarshal(data, &out)
}

// split parses the stored category column the same way it is written.
func split(category string) []string {
	var names []string
	for name := range strings.SplitSeq(category, ",") {
		if name = strings.TrimSpace(name); name != "" {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return slices.Compact(names)
}

// present shows what the model concluded so the categories can be decided by
// hand.
func present(r Result, toy models.Toy) {
	fmt.Printf("\n#%d %s\n", r.ToyID, truncate(r.ToyName, 100))
	if toy.SourceURL != "" {
		fmt.Printf("  link:       %s\n", toy.SourceURL)
	}
	fmt.Printf("  current:    %s\n", cmp.Or(toy.Category, "-"))
	fmt.Printf("  assigned:   %s\n", scores(r.Assigned, r.Existing))
}

// review asks which category of each flagged pair to keep and which uncertain
// categories to add; keep is false when the toy is skipped.
func review(input *bufio.Scanner, r Result) (categories []string, keep bool) {
	pairs := maps.Clone(r.Outscoring)
	if pairs == nil {
		pairs = map[string]string{}
	}
	maps.Copy(pairs, r.Similar)

	categories = slices.Clone(r.Assigned)
	for _, unruled := range slices.Sorted(maps.Keys(pairs)) {
		ruled := pairs[unruled]
		kind := "outscores"
		if _, ok := r.Similar[unruled]; ok {
			kind = "similar to"
		}
		fmt.Printf("  %s %s %s\n", score(unruled, r.Existing), kind, score(ruled, r.Existing))
		chosen, keep := askPair(input, ruled, unruled)
		if !keep {
			return nil, false
		}
		categories = slices.DeleteFunc(categories, func(c string) bool {
			return c == ruled || c == unruled
		})
		categories = append(categories, chosen...)
	}

	if len(r.Uncertain) > 0 {
		added, keep := askUncertain(input, r.Uncertain, r.Existing)
		if !keep {
			return nil, false
		}
		categories = append(categories, added...)
	}
	return categories, true
}

// askPair has no default: the point of the prompt is to make the call.
func askPair(input *bufio.Scanner, ruled, unruled string) ([]string, bool) {
	for {
		fmt.Printf("  [1] %s  [2] %s  [b]oth  [s]kip toy: ", ruled, unruled)
		if !input.Scan() {
			log.Fatal("input closed")
		}
		switch strings.ToLower(strings.TrimSpace(input.Text())) {
		case "1":
			return []string{ruled}, true
		case "2":
			return []string{unruled}, true
		case "b":
			return []string{ruled, unruled}, true
		case "s":
			return nil, false
		}
	}
}

// askUncertain offers the uncertain categories by number; Enter adds none.
func askUncertain(input *bufio.Scanner, uncertain []string, probs map[string]float64) ([]string, bool) {
	fmt.Println("  uncertain:")
	for i, name := range uncertain {
		fmt.Printf("    [%d] %s\n", i+1, score(name, probs))
	}
	for {
		fmt.Print("  add (e.g. 1,2) [none, s to skip toy]: ")
		if !input.Scan() {
			log.Fatal("input closed")
		}
		text := strings.TrimSpace(input.Text())
		switch {
		case text == "":
			return nil, true
		case strings.EqualFold(text, "s"):
			return nil, false
		}
		var added []string
		valid := true
		for part := range strings.SplitSeq(text, ",") {
			n, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil || n < 1 || n > len(uncertain) {
				valid = false
				break
			}
			added = append(added, uncertain[n-1])
		}
		if valid {
			return added, true
		}
		fmt.Printf("  enter numbers from 1 to %d separated by commas\n", len(uncertain))
	}
}

func confirm(input *bufio.Scanner, prompt string) bool {
	fmt.Print(prompt)
	if !input.Scan() {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(input.Text()))
	return answer == "y" || answer == "yes"
}

func printUpdates(updates []Update) {
	fmt.Println()
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tACTUAL\tNUEVA\tREVISADA\tTOY")
	for _, u := range updates {
		reviewed := ""
		if u.Reviewed {
			reviewed = "sí"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
			u.Toy.ToyID, cmp.Or(u.Toy.Category, "-"), strings.Join(u.Categories, ","), reviewed, truncate(u.Toy.ToyName, 50))
	}
	tw.Flush()
}

// score shows a category with its probability; unruled_similar categories
// are dropped from the scores, so theirs may be missing.
func score(name string, probs map[string]float64) string {
	if p, ok := probs[name]; ok {
		return fmt.Sprintf("%s %.2f", name, p)
	}
	return name
}

func scores(names []string, probs map[string]float64) string {
	parts := make([]string, len(names))
	for i, n := range names {
		parts[i] = score(n, probs)
	}
	return cmp.Or(strings.Join(parts, ", "), "-")
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}
