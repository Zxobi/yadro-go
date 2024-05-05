package util

import "strings"

func GeneratePlaceholders(count int) string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		placeholders = append(placeholders, "?")
	}
	return strings.Join(placeholders, ",")
}
