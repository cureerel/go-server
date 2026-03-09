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

// GenerateSlug converts a title into a clean, lowercase URL slug.
//
//	"Hello World!"  → "hello-world"
//	"Café au Lait"  → "cafe-au-lait"
//	"  ---  "       → "untitled"
//
// No suffix is ever appended here. Use GenerateUniqueSlug when you need
// to guarantee uniqueness against the database.
func GenerateSlug(title string) string {
	// Strip Unicode combining marks (accents): é → e, ñ → n, etc.
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

// GenerateUniqueSlug returns a slug that passes the caller-supplied uniqueness
// check. A brand-new title returns a clean slug with NO suffix whatsoever.
// A counter (-2, -3 …) is only appended when the base slug is already taken.
//
// Usage (in a service):
//
//	slug := utils.GenerateUniqueSlug(input.Title, func(candidate string) bool {
//	    exists, _ := s.blogRepo.SlugExists(ctx, candidate)
//	    return exists
//	})
//
// When updating a blog, skip the current post's own slug to avoid false conflicts:
//
//	slug := utils.GenerateUniqueSlug(*input.Title, func(candidate string) bool {
//	    if candidate == currentSlug { return false }  // same post — not a conflict
//	    exists, _ := s.blogRepo.SlugExists(ctx, candidate)
//	    return exists
//	})
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