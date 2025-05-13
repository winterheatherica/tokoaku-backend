package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func SlugifyText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9\\-]`)
	return re.ReplaceAllString(text, "")
}

func GenerateUniqueSlug(baseText string) (string, error) {
	baseSlug := SlugifyText(baseText)
	suffix := uuid.NewString()[:6]
	return fmt.Sprintf("%s-%s", baseSlug, suffix), nil
}
