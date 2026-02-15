package entity

import (
	"time"
	"gorm.io/gorm"

)

type Blog struct{
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null;size:200"`
	Content     string         `json:"content" gorm:"type:text"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	Status      string         `json:"status" gorm:"default:'draft';size:20"` // draft, published, archived
	Tags        string         `json:"tags" gorm:"size:500"` // comma-separated
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}


func (Blog) TableName() string {
	return "blog"
}