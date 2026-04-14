// internal/domain/entity/blog.go
package entity

import "time"

type BlogStatus string

const (
	BlogDraft     BlogStatus = "draft"
	BlogPublished BlogStatus = "published"
	BlogArchived  BlogStatus = "archived"
)

// Post access for readers (after publish).
type BlogAccessType string

const (
	AccessFree      BlogAccessType = "free"
	AccessMember    BlogAccessType = "member"
	AccessPaidCoins BlogAccessType = "paid_coins"
)

type Blog struct {
	ID           uint           `json:"id"`
	Title        string         `json:"title"`
	Slug         string         `json:"slug"`
	Content      string         `json:"content"`
	Keyword      string         `json:"keyword"`
	Tag          string         `json:"tag"`
	Excerpt      string         `json:"excerpt"`
	Thumbnail    string         `json:"thumbnail"`
	ThumbnailKey string         `json:"-"`
	Views        int64          `json:"views"`
	Status       BlogStatus     `json:"status"`
	AccessType   BlogAccessType `json:"access_type"`
	CoinPrice    int64          `json:"coin_price"`
	PublishedAt  *time.Time     `json:"published_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
