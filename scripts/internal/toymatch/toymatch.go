// Package toymatch finds the database toy a classification result refers to,
// shared by the update scripts. Results carry the IDs of the database they
// were classified from, so another database (prod, classified from dev) can
// only be matched by name.
package toymatch

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"reyes-magos-gr/store/models"
)

type Matcher struct {
	byName  bool
	ids     map[int64]models.Toy
	names   map[string][]models.Toy
	matched map[int64]string
}

// New matches by ID, or by name only when byName is set: an ID from another
// database could belong to an unrelated toy.
func New(toys []models.Toy, byName bool) *Matcher {
	m := &Matcher{
		byName:  byName,
		ids:     make(map[int64]models.Toy, len(toys)),
		names:   make(map[string][]models.Toy, len(toys)),
		matched: map[int64]string{},
	}
	for _, t := range toys {
		m.ids[t.ToyID] = t
		m.names[normalize(t.ToyName)] = append(m.names[normalize(t.ToyName)], t)
	}
	return m
}

// Find returns the toy for a result's ID and name. Names shared by several
// toys, and toys already matched by an earlier result, are errors rather
// than guesses.
func (m *Matcher) Find(id int64, name string) (models.Toy, error) {
	var toy models.Toy
	if m.byName {
		candidates := m.names[normalize(name)]
		switch len(candidates) {
		case 0:
			return toy, fmt.Errorf("no toy named %q in the database", name)
		case 1:
			toy = candidates[0]
		default:
			return toy, fmt.Errorf("%d toys named %q in the database", len(candidates), name)
		}
	} else {
		var ok bool
		if toy, ok = m.ids[id]; !ok {
			return toy, fmt.Errorf("toy %d is not in the database", id)
		}
	}
	if earlier, ok := m.matched[toy.ToyID]; ok {
		return toy, fmt.Errorf("toy %d was already matched by %q", toy.ToyID, earlier)
	}
	m.matched[toy.ToyID] = name
	return toy, nil
}

// normalize ignores case and spacing, which scraping can vary.
func normalize(name string) string {
	return strings.ToLower(strings.Join(strings.Fields(name), " "))
}

// Unmatched returns the database toys no result was matched to, by ID.
func (m *Matcher) Unmatched() []models.Toy {
	var toys []models.Toy
	for _, t := range m.ids {
		if _, ok := m.matched[t.ToyID]; !ok {
			toys = append(toys, t)
		}
	}
	slices.SortFunc(toys, func(a, b models.Toy) int { return cmp.Compare(a.ToyID, b.ToyID) })
	return toys
}
