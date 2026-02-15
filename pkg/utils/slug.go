package utils

import (
	"regexp"
	"strconv"  // ADDED THIS IMPORT
	"strings"
	"time"
)

func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)
	
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim hyphens from ends
	slug = strings.Trim(slug, "-")
	
	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	return slug + "-" + strconv.FormatInt(timestamp, 10)
}