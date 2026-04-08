// internal/infrastructure/postgres/models/user_model.go
package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/cureerel/cserver/internal/domain/entity"
)

type User struct {
	ID                 uint           `gorm:"primaryKey"`
	Name               string         `gorm:"not null;size:100"`
	Email              string         `gorm:"uniqueIndex;not null;size:100"`
	PasswordHash       string         `gorm:"column:password_hash;size:255"`
	Role               string         `gorm:"default:'user';size:20"`
	IsActive           bool           `gorm:"default:true"`
	IsVerified         bool           `gorm:"default:false"`
	UpgradeRequestedAt *time.Time     `gorm:"column:upgrade_requested_at"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string { return "users" }

func (m *User) ToDomain() *entity.User {
	return &entity.User{
		ID:                 m.ID,
		Name:               m.Name,
		Email:              m.Email,
		PasswordHash:       m.PasswordHash,
		Role:               m.Role,
		IsActive:           m.IsActive,
		IsVerified:         m.IsVerified,
		UpgradeRequestedAt: m.UpgradeRequestedAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func UserFromDomain(e *entity.User) *User {
	return &User{
		ID:                 e.ID,
		Name:               e.Name,
		Email:              e.Email,
		PasswordHash:       e.PasswordHash,
		Role:               e.Role,
		IsActive:           e.IsActive,
		IsVerified:         e.IsVerified,
		UpgradeRequestedAt: e.UpgradeRequestedAt,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}