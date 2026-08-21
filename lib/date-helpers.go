package lib

import (
	"fmt"
	"time"
)

var spanishMonthAbbr = [...]string{
	"ene", "feb", "mar", "abr", "may", "jun",
	"jul", "ago", "sep", "oct", "nov", "dic",
}

// DaysSince parses an RFC3339 date and returns how many whole days have
// elapsed since then, clamped to 0 for dates in the future.
func DaysSince(date string) (int, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, err
	}

	days := int(time.Since(t).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return days, nil
}

// RelativeSpanishLabel turns a day count into a short Spanish relative-time
// phrase: Hoy, Ayer, "Hace N días", "Hace N meses".
func RelativeSpanishLabel(days int) string {
	switch {
	case days == 0:
		return "Hoy"
	case days == 1:
		return "Ayer"
	case days < 30:
		return fmt.Sprintf("Hace %d días", days)
	default:
		months := days / 30
		if months == 1 {
			return "Hace 1 mes"
		}
		return fmt.Sprintf("Hace %d meses", months)
	}
}

// FormatSpanishDate parses an RFC3339 date and formats it like "20 ago 2026".
func FormatSpanishDate(date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d %s %d", t.Day(), spanishMonthAbbr[t.Month()-1], t.Year()), nil
}
