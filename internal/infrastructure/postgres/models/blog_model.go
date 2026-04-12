// internal/infrastructure/postgres/models/blog_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"gorm.io/gorm"
)

type Blog struct {
	ID            uint           `gorm:"primaryKey"`
	Title         string         `gorm:"not null;size:500"`
	Slug          string         `gorm:"uniqueIndex;not null;size:200"`
	Content       string         `gorm:"type:text"`
	Excerpt       string         `gorm:"type:text"`
	AuthorID      uint           `gorm:"not null;index"`
	Status        string         `gorm:"default:'draft';size:20;index"`
	AccessType    string         `gorm:"column:access_type;size:20;default:free"`
	CoinPrice     int64          `gorm:"column:coin_price;default:0"`
	Tags          string         `gorm:"type:text"`
	CoverImageURL string         `gorm:"column:cover_image_url;type:text"`
	CoverImageKey string         `gorm:"column:cover_image_key;type:text"`
	ViewsTotal    int64          `gorm:"column:views_total;default:0"`
	PublishedAt   *time.Time     `gorm:"column:published_at"`
	SubmittedForReviewAt *time.Time `gorm:"column:submitted_for_review_at"`
	ReviewedByID  *uint          `gorm:"column:reviewed_by_id"`
	ReviewNote    string         `gorm:"column:review_note;type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (Blog) TableName() string { return "blogs" }

func (m *Blog) ToDomain() *entity.Blog {
	return &entity.Blog{
		ID:            m.ID,
		Title:         m.Title,
		Slug:          m.Slug,
		Content:       m.Content,
		Excerpt:       m.Excerpt,
		AuthorID:      m.AuthorID,
		Status:        entity.BlogStatus(m.Status),
		AccessType:    entity.BlogAccessType(m.AccessType),
		CoinPrice:     m.CoinPrice,
		Tags:          m.Tags,
		CoverImageURL: m.CoverImageURL,
		CoverImageKey: m.CoverImageKey,
		ViewsTotal:    m.ViewsTotal,
		PublishedAt:   m.PublishedAt,
		SubmittedForReviewAt: m.SubmittedForReviewAt,
		ReviewedByID:  m.ReviewedByID,
		ReviewNote:    m.ReviewNote,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func BlogFromDomain(e *entity.Blog) *Blog {
	at := string(e.AccessType)
	if at == "" {
		at = "free"
	}
	return &Blog{
		ID:            e.ID,
		Title:         e.Title,
		Slug:          e.Slug,
		Content:       e.Content,
		Excerpt:       e.Excerpt,
		AuthorID:      e.AuthorID,
		Status:        string(e.Status),
		AccessType:    at,
		CoinPrice:     e.CoinPrice,
		Tags:          e.Tags,
		CoverImageURL: e.CoverImageURL,
		CoverImageKey: e.CoverImageKey,
		ViewsTotal:    e.ViewsTotal,
		PublishedAt:   e.PublishedAt,
		SubmittedForReviewAt: e.SubmittedForReviewAt,
		ReviewedByID:  e.ReviewedByID,
		ReviewNote:    e.ReviewNote,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}