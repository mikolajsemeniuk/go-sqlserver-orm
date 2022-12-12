package yachts

import (
	"regexp"
	"strings"
)

func ToSnakeCase(word string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	word = re.ReplaceAllString(word, "${1}_${2}")
	return strings.ToLower(word)
}
