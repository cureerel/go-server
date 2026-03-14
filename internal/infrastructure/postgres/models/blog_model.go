// internal/infrastructure/postgres/models/blog_model.go
package models

import (
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
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
	Tags          string         `gorm:"type:text"`
	CoverImageURL string         `gorm:"column:cover_image_url;type:text"`
	CoverImageKey string         `gorm:"column:cover_image_key;type:text"`
	ViewsTotal    int64          `gorm:"column:views_total;default:0"`
	PublishedAt   *time.Time     `gorm:"column:published_at"`
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
		Tags:          m.Tags,
		CoverImageURL: m.CoverImageURL,
		CoverImageKey: m.CoverImageKey,
		ViewsTotal:    m.ViewsTotal,
		PublishedAt:   m.PublishedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func BlogFromDomain(e *entity.Blog) *Blog {
	return &Blog{
		ID:            e.ID,
		Title:         e.Title,
		Slug:          e.Slug,
		Content:       e.Content,
		Excerpt:       e.Excerpt,
		AuthorID:      e.AuthorID,
		Status:        string(e.Status),
		Tags:          e.Tags,
		CoverImageURL: e.CoverImageURL,
		CoverImageKey: e.CoverImageKey,
		ViewsTotal:    e.ViewsTotal,
		PublishedAt:   e.PublishedAt,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}