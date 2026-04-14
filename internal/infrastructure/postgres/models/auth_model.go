package models

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;size:255"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (m *RefreshToken) ToDomain() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		Revoked:   m.Revoked,
		CreatedAt: m.CreatedAt,
	}
}

func RefreshTokenFromDomain(e *entity.RefreshToken) *RefreshToken {
	return &RefreshToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		Revoked:   e.Revoked,
		CreatedAt: e.CreatedAt,
	}
}

type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey"`
	TokenHash string    `gorm:"not null;size:255;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}
