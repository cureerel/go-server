// internal/domain/entity/blog.go
package entity

import "time"

type BlogStatus string

const (
	BlogDraft      BlogStatus = "draft"
	BlogInReview   BlogStatus = "in_review"
	BlogPublished  BlogStatus = "published"
	BlogRejected   BlogStatus = "rejected"
	BlogArchived   BlogStatus = "archived"
)

// Post access for readers (after publish).
type BlogAccessType string

const (
	AccessFree     BlogAccessType = "free"
	AccessMember   BlogAccessType = "member"
	AccessPaidCoins BlogAccessType = "paid_coins"
)

type Blog struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt"`
	AuthorID      uint       `json:"author_id"`
	Status        BlogStatus `json:"status"`
	AccessType    BlogAccessType `json:"access_type"`
	CoinPrice     int64      `json:"coin_price"`
	Tags          string     `json:"tags"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	CoverImageKey string     `json:"-"`
	ViewsTotal    int64      `json:"views_total"`
	PublishedAt   *time.Time `json:"published_at,omitempty"`
	SubmittedForReviewAt *time.Time `json:"submitted_for_review_at,omitempty"`
	ReviewedByID  *uint      `json:"reviewed_by_id,omitempty"`
	ReviewNote    string     `json:"review_note,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type BlogAuthor struct {
	BlogID uint `json:"blog_id"`
	UserID uint `json:"user_id"`
	Role   string `json:"role"`
}