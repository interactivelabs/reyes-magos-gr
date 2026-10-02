// Command scrape-toys reads every toy from the database, fetches its Amazon
// SourceURL and writes one JSON file per toy with the product title, feature
// bullets and the product details tables.
//
// Usage:
//
//	go run ./cmd/scrape-toys -out ./scraped -delay 3s
//
// It uses the same TURSO_* environment variables as the app.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"

	"golang.org/x/net/html"
)

type ToyFile struct {
	ToyID     int64     `json:"toy_id"`
	ToyName   string    `json:"toy_name"`
	SourceURL string    `json:"source_url"`
	ScrapedAt time.Time `json:"scraped_at"`
	ProductDetails
}

var errBlocked = errors.New("amazon returned a captcha/robot check page")

// blockedMarker shows up on the "continue shopping" interstitial Amazon
// serves when it thinks we're a bot.
const blockedMarker = "Haz clic en el botón de abajo para continuar comprando"

func main() {
	outDir := flag.String("out", "scraped", "directory where JSON files are written")
	delay := flag.Duration("delay", 3*time.Second, "pause between requests")
	force := flag.Bool("force", false, "re-scrape toys that already have a JSON file")
	saveHTML := flag.Bool("save-html", false, "also save the raw HTML next to each JSON (for debugging)")
	flag.Parse()

	db, connector, dir, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := connector.Sync(); err != nil {
		log.Fatal("Error syncing database: ", err)
	}
	defer os.RemoveAll(dir)
	defer connector.Close()
	defer db.Close()

	toys, err := store.NewToysStore(db).GetToys()
	if err != nil {
		log.Fatal("Error reading toys: ", err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	var ok, skipped, failed int

	for i, toy := range toys {
		jsonPath := filepath.Join(*outDir, fmt.Sprintf("toy-%d.json", toy.ToyID))
		prefix := fmt.Sprintf("[%d/%d] toy %d", i+1, len(toys), toy.ToyID)

		if strings.TrimSpace(toy.SourceURL) == "" {
			log.Printf("%s: no source_url, skipping", prefix)
			skipped++
			continue
		}
		if _, err := os.Stat(jsonPath); err == nil && !*force {
			log.Printf("%s: %s exists, skipping", prefix, jsonPath)
			skipped++
			continue
		}

		err := scrapeToy(client, toy, jsonPath, *saveHTML)
		if errors.Is(err, errBlocked) {
			// Back off for the rest of the run and retry this toy once.
			*delay += time.Second
			log.Printf("%s: blocked, retrying once with delay %s", prefix, *delay)
			time.Sleep(*delay)
			err = scrapeToy(client, toy, jsonPath, *saveHTML)
		}

		if err != nil {
			log.Printf("%s: %v", prefix, err)
			failed++
		} else {
			log.Printf("%s: wrote %s", prefix, jsonPath)
			ok++
		}

		if i < len(toys)-1 {
			time.Sleep(*delay)
		}
	}

	log.Printf("Done: %d written, %d skipped, %d failed (final delay %s)", ok, skipped, failed, *delay)
}

func scrapeToy(client *http.Client, toy models.Toy, jsonPath string, saveHTML bool) error {
	body, err := fetch(client, toy.SourceURL)
	if err != nil {
		return err
	}

	if saveHTML {
		htmlPath := strings.TrimSuffix(jsonPath, ".json") + ".html"
		if err := os.WriteFile(htmlPath, body, 0o644); err != nil {
			return err
		}
	}

	content := string(body)
	if strings.Contains(content, blockedMarker) {
		return errBlocked
	}

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("parsing html: %w", err)
	}

	details := parseProduct(doc)
	if details.Title == "" {
		if findByID(doc, "captchacharacters") != nil || strings.Contains(content, "/errors/validateCaptcha") {
			return errBlocked
		}
		return errors.New("productTitle not found, page layout may differ (use -save-html to inspect)")
	}

	out, err := json.MarshalIndent(ToyFile{
		ToyID:          toy.ToyID,
		ToyName:        toy.ToyName,
		SourceURL:      toy.SourceURL,
		ScrapedAt:      time.Now().UTC(),
		ProductDetails: details,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonPath, out, 0o644)
}

func fetch(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Amazon serves a robot-check page to obvious bots, so look like a browser.
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-MX,es;q=0.9,en;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusServiceUnavailable {
		return nil, errBlocked
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s", res.Status)
	}
	return io.ReadAll(res.Body)
}
