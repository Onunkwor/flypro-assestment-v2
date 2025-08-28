package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func ExpensesCacheKey(filters map[string]interface{}, offset, limit int) string {

	var parts []string
	for k, v := range filters {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(parts)

	rawKey := fmt.Sprintf("expenses:%s:offset=%d:limit=%d", strings.Join(parts, ":"), offset, limit)

	h := sha1.Sum([]byte(rawKey))
	return "expenses:" + hex.EncodeToString(h[:])
}
