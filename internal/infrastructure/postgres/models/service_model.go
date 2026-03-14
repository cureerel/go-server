// internal/infrastructure/postgres/models/service_model.go
package models

import (
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
	"gorm.io/gorm"
)

type Service struct {
	ID            uint           `gorm:"primaryKey"`
	OwnerID       uint           `gorm:"not null;index"`
	Title         string         `gorm:"not null;size:200"`
	Description   string         `gorm:"type:text"`
	PriceUSDCents int64          `gorm:"column:price_usd_cents;not null"`
	Status        string         `gorm:"default:'pending';size:20;index"`
	CoverImageURL string         `gorm:"column:cover_image_url;type:text"`
	CoverImageKey string         `gorm:"column:cover_image_key;type:text"`
	ViewsTotal    int64          `gorm:"column:views_total;default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (Service) TableName() string { return "services" }

func (m *Service) ToDomain() *entity.Service {
	return &entity.Service{
		ID:            m.ID,
		OwnerID:       m.OwnerID,
		Title:         m.Title,
		Description:   m.Description,
		PriceUSDCents: m.PriceUSDCents,
		Status:        m.Status,
		CoverImageURL: m.CoverImageURL,
		CoverImageKey: m.CoverImageKey,
		ViewsTotal:    m.ViewsTotal,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func ServiceFromDomain(e *entity.Service) *Service {
	return &Service{
		ID:            e.ID,
		OwnerID:       e.OwnerID,
		Title:         e.Title,
		Description:   e.Description,
		PriceUSDCents: e.PriceUSDCents,
		Status:        e.Status,
		CoverImageURL: e.CoverImageURL,
		CoverImageKey: e.CoverImageKey,
		ViewsTotal:    e.ViewsTotal,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}