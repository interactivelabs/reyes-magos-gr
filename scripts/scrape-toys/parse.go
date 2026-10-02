package main

import (
	"strings"

	"golang.org/x/net/html"
)

// ProductDetails is what we extract from an Amazon product page. Attribute
// values are `any` on purpose: most are plain strings, but some cells hold
// lists (best-seller ranks) or small objects (customer reviews).
type ProductDetails struct {
	Title                string           `json:"title"`
	TitleDifferentiators string           `json:"title_differentiators,omitempty"`
	Features             []string         `json:"features"`
	Details              []DetailsSection `json:"details"`
}

type DetailsSection struct {
	Name       string         `json:"name"`
	Attributes map[string]any `json:"attributes"`
}

func parseProduct(doc *html.Node) ProductDetails {
	return ProductDetails{
		Title:                text(findByID(doc, "productTitle")),
		TitleDifferentiators: parseTitleDifferentiators(findByID(doc, "titleSection")),
		Features:             parseFeatures(findByID(doc, "feature-bullets")),
		Details:              parseProdDetails(findByID(doc, "prodDetails")),
	}
}

func parseTitleDifferentiators(titleSection *html.Node) string {
	return text(findFirst(titleSection, func(n *html.Node) bool {
		return hasClass(n, "dp-title-differentiators")
	}))
}

func parseFeatures(bullets *html.Node) []string {
	features := []string{}
	for _, li := range findAll(bullets, isElement("li")) {
		if t := text(li); t != "" {
			features = append(features, t)
		}
	}
	return features
}

// parseProdDetails walks #prodDetails keeping track of the closest expander
// heading ("Detalles del producto", "Estilo", ...) and collects every
// th/td row of the tables found under it. Tables outside an expander end up
// in a section named after the #prodDetails heading.
func parseProdDetails(prodDetails *html.Node) []DetailsSection {
	sections := []DetailsSection{}
	if prodDetails == nil {
		return sections
	}

	byName := map[string]int{}
	add := func(section, key string, value any) {
		i, ok := byName[section]
		if !ok {
			i = len(sections)
			byName[section] = i
			sections = append(sections, DetailsSection{Name: section, Attributes: map[string]any{}})
		}
		sections[i].Attributes[key] = value
	}

	defaultSection := text(findFirst(prodDetails, isElement("h1")))
	if defaultSection == "" {
		defaultSection = "General"
	}

	var walk func(n *html.Node, section string)
	walk = func(n *html.Node, section string) {
		if hasClass(n, "a-expander-container") {
			if prompt := text(findFirst(n, func(c *html.Node) bool {
				return hasClass(c, "a-expander-prompt")
			})); prompt != "" {
				section = prompt
			}
		}

		if isElement("tr")(n) {
			th := findFirst(n, isElement("th"))
			td := findFirst(n, isElement("td"))
			if key := text(th); key != "" && td != nil {
				add(section, key, cellValue(td))
			}
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, section)
		}
	}
	walk(prodDetails, defaultSection)

	return sections
}

func cellValue(td *html.Node) any {
	if reviews := findByID(td, "averageCustomerReviews"); reviews != nil {
		rating := attr(findByID(reviews, "acrPopover"), "title")
		count := strings.Trim(text(findByID(reviews, "acrCustomerReviewText")), "()")
		return map[string]string{"rating": rating, "count": count}
	}

	if lis := findAll(td, isElement("li")); len(lis) > 0 {
		items := make([]string, 0, len(lis))
		for _, li := range lis {
			if t := text(li); t != "" {
				items = append(items, t)
			}
		}
		return items
	}

	return text(td)
}

// --- small html.Node helpers ---

func isElement(tag string) func(*html.Node) bool {
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tag
	}
}

func attr(n *html.Node, key string) string {
	if n == nil {
		return ""
	}
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func hasClass(n *html.Node, class string) bool {
	if n == nil || n.Type != html.ElementNode {
		return false
	}
	for c := range strings.FieldsSeq(attr(n, "class")) {
		if c == class {
			return true
		}
	}
	return false
}

func findFirst(n *html.Node, match func(*html.Node) bool) *html.Node {
	if n == nil {
		return nil
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if match(c) {
			return c
		}
		if found := findFirst(c, match); found != nil {
			return found
		}
	}
	return nil
}

func findAll(n *html.Node, match func(*html.Node) bool) []*html.Node {
	var out []*html.Node
	if n == nil {
		return out
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if match(c) {
			out = append(out, c)
		}
		out = append(out, findAll(c, match)...)
	}
	return out
}

func findByID(n *html.Node, id string) *html.Node {
	return findFirst(n, func(c *html.Node) bool {
		return c.Type == html.ElementNode && attr(c, "id") == id
	})
}

// text returns the visible text under n with whitespace collapsed,
// skipping script and style contents.
func text(n *html.Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			b.WriteByte(' ')
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.Join(strings.Fields(b.String()), " ")
}
