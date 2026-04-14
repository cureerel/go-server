// internal/infrastructure/postgres/models/blog_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"gorm.io/gorm"
)

type Blog struct {
	ID           uint       `gorm:"primaryKey"`
	Title        string     `gorm:"not null;size:500"`
	Slug         string     `gorm:"uniqueIndex;not null;size:200"`
	Content      string     `gorm:"type:text"`
	Keyword      string     `gorm:"size:500"`
	Tag          string     `gorm:"size:500"`
	Excerpt      string     `gorm:"type:text"`
	Thumbnail    string     `gorm:"column:thumbnail;type:text"`
	ThumbnailKey string     `gorm:"column:thumbnail_key;type:text"`
	Views        int64      `gorm:"column:views;default:0"`
	Status       string     `gorm:"default:'draft';size:20;index"`
	AccessType   string     `gorm:"column:access_type;size:20;default:free"`
	CoinPrice    int64      `gorm:"column:coin_price;default:0"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Blog) TableName() string { return "blogs" }

func (m *Blog) ToDomain() *entity.Blog {
	return &entity.Blog{
		ID:           m.ID,
		Title:        m.Title,
		Slug:         m.Slug,
		Content:      m.Content,
		Keyword:      m.Keyword,
		Tag:          m.Tag,
		Excerpt:      m.Excerpt,
		Thumbnail:    m.Thumbnail,
		ThumbnailKey: m.ThumbnailKey,
		Views:        m.Views,
		Status:       entity.BlogStatus(m.Status),
		AccessType:   entity.BlogAccessType(m.AccessType),
		CoinPrice:    m.CoinPrice,
		PublishedAt:  m.PublishedAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func BlogFromDomain(e *entity.Blog) *Blog {
	at := string(e.AccessType)
	if at == "" {
		at = "free"
	}
	return &Blog{
		ID:           e.ID,
		Title:        e.Title,
		Slug:         e.Slug,
		Content:      e.Content,
		Keyword:      e.Keyword,
		Tag:          e.Tag,
		Excerpt:      e.Excerpt,
		Thumbnail:    e.Thumbnail,
		ThumbnailKey: e.ThumbnailKey,
		Views:        e.Views,
		Status:       string(e.Status),
		AccessType:   at,
		CoinPrice:    e.CoinPrice,
		PublishedAt:  e.PublishedAt,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
