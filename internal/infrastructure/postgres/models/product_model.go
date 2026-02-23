package models

import (
    "time"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null;size:200"`
    Description string    `gorm:"type:text"`
    Price       int64     `gorm:"not null"`
    Currency    string    `gorm:"not null;size:10"`
    IsActive    bool      `gorm:"default:true"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (Product) TableName() string {
    return "products"
}

func (m *Product) ToDomain() *entity.Product {
    return &entity.Product{
        ID:          m.ID,
        Name:        m.Name,
        Description: m.Description,
        Price:       m.Price,
        Currency:    entity.Currency(m.Currency),
        IsActive:    m.IsActive,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
    }
}

func ProductFromDomain(e *entity.Product) *Product {
    return &Product{
        ID:          e.ID,
        Name:        e.Name,
        Description: e.Description,
        Price:       e.Price,
        Currency:    string(e.Currency),
        IsActive:    e.IsActive,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,
    }
}