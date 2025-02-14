package lib

import "strings"

func RepeatStringWithDevider(s string, count int, d string) string {
	if count <= 0 {
		return ""
	}

	var builder strings.Builder
	for i := 0; i < count; i++ {
		builder.WriteString(s)
		if i < count-1 {
			builder.WriteString(d)
		}
	}
	return builder.String()
}
