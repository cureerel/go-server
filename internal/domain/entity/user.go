package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null;size:100"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password     string         `json:"-" gorm:"-"` // Deprecated: use PasswordHash
	PasswordHash string         `json:"-" gorm:"column:password_hash;size:255"`
	Role         string         `json:"role" gorm:"default:'user';size:20"` // admin, user, editor
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

// GetID returns string ID for JWT claims
func (u *User) GetID() string {
	return fmt.Sprintf("%d", u.ID)
}

// MigratePassword moves Password to PasswordHash if needed
func (u *User) MigratePassword() {
	if u.Password != "" && u.PasswordHash == "" {
		u.PasswordHash = u.Password
		u.Password = ""
	}
}