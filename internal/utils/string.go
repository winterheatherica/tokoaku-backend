package utils

import (
	"regexp"
	"strings"
)

var spaceRegexp = regexp.MustCompile(`\s+`)

func ToSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	s = spaceRegexp.ReplaceAllString(s, "_")
	return strings.ToLower(s)
}
