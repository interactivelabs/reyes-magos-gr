// Command update-toy-ages writes the age ranges proposed by classify-toy-ages
// (scraped/ages.json) to the toys table.
//
// Usage:
//
//	go run ./scripts/update-toy-ages               # dev, from TURSO_*
//	go run ./scripts/update-toy-ages -production   # prod, from PROD_TURSO_*
//
// By default it uses the same TURSO_* environment variables as the app. With
// -production every one of them must be overridden by its PROD_ counterpart
// (PROD_TURSO_DATABASE_URL, ...), so prod is never reached by accident with
// the dev settings.
//
// Toys flagged needs_review are shown with what the model saw and their ages
// are prompted; Enter keeps the proposed value. Nothing is written until the
// final list of changes is confirmed.
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

// maxAge bounds what can be typed at the prompt.
const maxAge = 18

// Result is the part of classify-toy-ages' output this script needs.
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

type Update struct {
	Toy            models.Toy
	AgeMin, AgeMax int64
	Reviewed       bool
}

func main() {
	in := flag.String("in", "scraped/ages.json", "ages written by classify-toy-ages")
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
		u := Update{Toy: toy, AgeMin: r.AgeMin, AgeMax: r.AgeMax}
		if r.NeedsReview {
			present(r, toy)
			var keep bool
			if u.AgeMin, u.AgeMax, keep = askRange(input, r); !keep {
				fmt.Println("  skipped")
				continue
			}
			u.Reviewed = true
		}
		if u.AgeMin != toy.AgeMin || u.AgeMax != toy.AgeMax {
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
		if err := toysStore.UpdateToyAges(u.Toy.ToyID, u.AgeMin, u.AgeMax); err != nil {
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
	var results []Result
	return results, json.Unmarshal(data, &results)
}

// present shows what the model concluded so the ages can be decided by hand.
func present(r Result, toy models.Toy) {
	fmt.Printf("\n#%d %s\n", r.ToyID, r.ToyName)
	if toy.SourceURL != "" {
		fmt.Printf("  link:          %s\n", toy.SourceURL)
	}
	fmt.Printf("  review:        %s\n", strings.Join(r.ReviewReasons, "; "))
	fmt.Printf("  current:       %d-%d\n", toy.AgeMin, toy.AgeMax)
	fmt.Printf("  proposed:      %d-%d\n", r.AgeMin, r.AgeMax)
	manufacturer := cmp.Or(r.ManufacturerAge, "-")
	if r.ManufacturerOK != nil {
		manufacturer += fmt.Sprintf(" (plausible %.2f)", *r.ManufacturerOK)
	}
	fmt.Printf("  manufacturer:  %s\n", manufacturer)
	fmt.Printf("  small parts:   %.2f\n", r.SmallParts)
	fmt.Printf("  youngest:      %.2f  %s\n", r.YoungestConfidence, top(r.YoungestProbs, 3))
	fmt.Printf("  oldest:        %.2f  %s\n", r.OldestConfidence, top(r.OldestProbs, 3))
}

// askRange prompts for the ages the review is about, keeping the proposed
// value for the other; keep is false when the toy is skipped.
func askRange(input *bufio.Scanner, r Result) (ageMin, ageMax int64, keep bool) {
	askMin, askMax := reviewedAges(r.ReviewReasons)
	ageMin, ageMax = r.AgeMin, r.AgeMax
	if askMin {
		most := int64(maxAge)
		if !askMax {
			most = ageMax
		}
		if ageMin, keep = askAge(input, "age_min", ageMin, 0, most); !keep {
			return 0, 0, false
		}
	}
	if askMax {
		if ageMax, keep = askAge(input, "age_max", max(ageMax, ageMin), ageMin, maxAge); !keep {
			return 0, 0, false
		}
	}
	return ageMin, ageMax, true
}

// reviewedAges maps classify-toy-ages' review reasons to the age they put in
// doubt. A reason it doesn't know asks for both.
func reviewedAges(reasons []string) (askMin, askMax bool) {
	for _, reason := range reasons {
		switch {
		case strings.HasPrefix(reason, "uncertain youngest"), strings.HasPrefix(reason, "small parts"):
			askMin = true
		case strings.HasPrefix(reason, "uncertain oldest"), strings.HasPrefix(reason, "oldest stage"):
			askMax = true
		default:
			return true, true
		}
	}
	if !askMin && !askMax {
		return true, true
	}
	return askMin, askMax
}

func askAge(input *bufio.Scanner, label string, proposed, least, most int64) (int64, bool) {
	for {
		fmt.Printf("  %s [%d, s to skip]: ", label, proposed)
		if !input.Scan() {
			log.Fatal("input closed")
		}
		text := strings.TrimSpace(input.Text())
		switch {
		case text == "":
			return proposed, true
		case strings.EqualFold(text, "s"):
			return 0, false
		}
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil || n < least || n > most {
			fmt.Printf("  enter a number from %d to %d\n", least, most)
			continue
		}
		return n, true
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
		fmt.Fprintf(tw, "%d\t%d-%d\t%d-%d\t%s\t%s\n",
			u.Toy.ToyID, u.Toy.AgeMin, u.Toy.AgeMax, u.AgeMin, u.AgeMax, reviewed, truncate(u.Toy.ToyName, 50))
	}
	tw.Flush()
}

// top lists the n most likely options, highest first.
func top(probs map[string]float64, n int) string {
	names := slices.Collect(maps.Keys(probs))
	slices.SortFunc(names, func(a, b string) int {
		return cmp.Or(cmp.Compare(probs[b], probs[a]), strings.Compare(a, b))
	})
	parts := make([]string, 0, n)
	for _, name := range names[:min(n, len(names))] {
		parts = append(parts, fmt.Sprintf("%s %.2f", name, probs[name]))
	}
	return strings.Join(parts, ", ")
}

func truncate(s string, n int) string {
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}
