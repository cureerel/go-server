// internal/domain/entity/blog.go
package entity

import "time"

type Blog struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Content        string    `json:"content"`
	AuthorID       uint      `json:"author_id"`
	Status         string    `json:"status"`
	Tags           string    `json:"tags"`
	CoverImageURL  string    `json:"cover_image_url,omitempty"`
	CoverImageKey  string    `json:"-"` // internal key for deletion, not exposed
	ViewsTotal     int64     `json:"views_total"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}