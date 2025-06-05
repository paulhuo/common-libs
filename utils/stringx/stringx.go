package stringx

import (
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// Random returns a random string of the given length consisting of letters and digits.
func Random(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ToSnake converts a camelCase or spaced string to snake_case.
func ToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if r == ' ' || r == '-' {
			result = append(result, '_')
			continue
		}
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
			continue
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// ToCamel converts underscore/space/dash separated strings to CamelCase.
func ToCamel(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	parts := strings.Fields(s)
	for i, p := range parts {
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}

// Truncate shortens a string to max length. If truncated, "..." is appended.
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		if max <= 0 {
			return ""
		}
		return s[:max]
	}
	return s[:max-3] + "..."
}
