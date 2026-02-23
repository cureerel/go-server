package models

import (
    "time"

    "gorm.io/gorm"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type Blog struct {
    ID        uint           `gorm:"primaryKey"`
    Title     string         `gorm:"not null;size:200"`
    Slug      string         `gorm:"uniqueIndex;not null;size:200"`
    Content   string         `gorm:"type:text"`
    AuthorID  uint           `gorm:"not null;index"`
    Status    string         `gorm:"default:'draft';size:20"`
    Tags      string         `gorm:"size:500"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Author    User           `gorm:"foreignKey:AuthorID"`
}

func (Blog) TableName() string {
    return "blogs"
}

func (m *Blog) ToDomain() *entity.Blog {
    return &entity.Blog{
        ID:        m.ID,
        Title:     m.Title,
        Slug:      m.Slug,
        Content:   m.Content,
        AuthorID:  m.AuthorID,
        Status:    m.Status,
        Tags:      m.Tags,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}

func BlogFromDomain(e *entity.Blog) *Blog {
    return &Blog{
        ID:        e.ID,
        Title:     e.Title,
        Slug:      e.Slug,
        Content:   e.Content,
        AuthorID:  e.AuthorID,
        Status:    e.Status,
        Tags:      e.Tags,
        CreatedAt: e.CreatedAt,
        UpdatedAt: e.UpdatedAt,
    }
}