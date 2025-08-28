package utils

import "strings"

func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	return s
}

func NormalizeCategory(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
