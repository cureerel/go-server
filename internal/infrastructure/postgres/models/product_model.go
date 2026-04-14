// internal/infrastructure/postgres/models/product_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	SKU         string `gorm:"uniqueIndex;not null;size:30"`
	Name        string `gorm:"not null;size:200"`
	Description string `gorm:"type:text"`
	Type        string `gorm:"not null;size:20;default:'digital'"`
	Price       int64  `gorm:"not null"`
	Currency    string `gorm:"not null;size:10;default:'USD'"`
	Stock       int    `gorm:"not null;default:-1"` // -1 = unlimited
	ImageURL    string `gorm:"column:image_url;type:text"`
	IsActive    bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Product) TableName() string { return "products" }

func (m *Product) ToDomain() *entity.Product {
	return &entity.Product{
		ID:          m.ID,
		SKU:         m.SKU,
		Name:        m.Name,
		Description: m.Description,
		Type:        entity.ProductType(m.Type),
		Price:       m.Price,
		Currency:    entity.Currency(m.Currency),
		Stock:       m.Stock,
		ImageURL:    m.ImageURL,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ProductFromDomain(e *entity.Product) *Product {
	return &Product{
		ID:          e.ID,
		SKU:         e.SKU,
		Name:        e.Name,
		Description: e.Description,
		Type:        string(e.Type),
		Price:       e.Price,
		Currency:    string(e.Currency),
		Stock:       e.Stock,
		ImageURL:    e.ImageURL,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
