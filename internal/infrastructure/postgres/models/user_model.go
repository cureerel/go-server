// internal/infrastructure/postgres/models/user_model.go
package models

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                 uint           `gorm:"primaryKey"`
	Username           string         `gorm:"not null;size:100;column:username"`
	Email              string         `gorm:"uniqueIndex;not null;size:100"`
	PasswordHash       string         `gorm:"column:password_hash;size:255"`
	Role               string         `gorm:"default:'user';size:20"`
	FirstName          string         `gorm:"size:100;column:first_name"`
	LastName           string         `gorm:"size:100;column:last_name"`
	Country            string         `gorm:"size:100"`
	PhoneNumber        string         `gorm:"size:50;column:phone_number"`
	Address            string         `gorm:"type:text"`
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
		Username:           m.Username,
		Email:              m.Email,
		PasswordHash:       m.PasswordHash,
		Role:               m.Role,
		FirstName:          m.FirstName,
		LastName:           m.LastName,
		Country:            m.Country,
		PhoneNumber:        m.PhoneNumber,
		Address:            m.Address,
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
		Username:           e.Username,
		Email:              e.Email,
		PasswordHash:       e.PasswordHash,
		Role:               e.Role,
		FirstName:          e.FirstName,
		LastName:           e.LastName,
		Country:            e.Country,
		PhoneNumber:        e.PhoneNumber,
		Address:            e.Address,
		IsActive:           e.IsActive,
		IsVerified:         e.IsVerified,
		UpgradeRequestedAt: e.UpgradeRequestedAt,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}
