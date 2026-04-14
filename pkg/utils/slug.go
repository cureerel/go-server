// pkg/utils/slug.go
package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	reBadChars  = regexp.MustCompile(`[^a-z0-9]+`)
	reMultiDash = regexp.MustCompile(`-{2,}`)
)

func GenerateSlug(title string) string {

	t := transform.Chain(
		norm.NFD,
		transform.RemoveFunc(func(r rune) bool { return unicode.Is(unicode.Mn, r) }),
		norm.NFC,
	)
	out, _, _ := transform.String(t, title)

	s := strings.ToLower(out)
	s = reBadChars.ReplaceAllString(s, "-")
	s = reMultiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "untitled"
	}
	return s
}

func GenerateUniqueSlug(title string, exists func(slug string) bool) string {
	base := GenerateSlug(title)

	// Happy path: fresh slug — return it clean, no suffix.
	if !exists(base) {
		return base
	}

	// Conflict: try base-2, base-3 … base-999
	for i := 2; i <= 999; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !exists(candidate) {
			return candidate
		}
	}

	// Practically unreachable (999 duplicate titles).
	return base
}
