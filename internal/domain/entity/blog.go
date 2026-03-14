// internal/domain/entity/blog.go
package entity

import "time"

type BlogStatus string

const (
	BlogDraft     BlogStatus = "draft"
	BlogPublished BlogStatus = "published"
	BlogArchived  BlogStatus = "archived"
)

type Blog struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt"`   
	AuthorID      uint       `json:"author_id"`
	Status        BlogStatus `json:"status"`
	Tags          string     `json:"tags"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	CoverImageKey string     `json:"-"`
	ViewsTotal    int64      `json:"views_total"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}