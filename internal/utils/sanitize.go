package utils

import "strings"

func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	return s
}
