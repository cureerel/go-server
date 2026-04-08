// internal/infrastructure/postgres/models/session.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Session struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uint      `gorm:"not null;index"`
	TokenHash  string    `gorm:"not null;size:255;uniqueIndex"`
	UserAgent  string    `gorm:"size:512"`
	IPAddress  string    `gorm:"size:45"` // supports IPv6
	ExpiresAt  time.Time `gorm:"not null;index"`
	Revoked    bool      `gorm:"default:false;index"`
	CreatedAt  time.Time
	LastActive time.Time `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserID"`
}

func (Session) TableName() string {
	return "sessions"
}

func (m *Session) ToDomain() *entity.Session {
	return &entity.Session{
		ID:         m.ID,
		UserID:     m.UserID,
		TokenHash:  m.TokenHash,
		UserAgent:  m.UserAgent,
		IPAddress:  m.IPAddress,
		ExpiresAt:  m.ExpiresAt,
		Revoked:    m.Revoked,
		CreatedAt:  m.CreatedAt,
		LastActive: m.LastActive,
	}
}

func SessionFromDomain(e *entity.Session) *Session {
	return &Session{
		ID:         e.ID,
		UserID:     e.UserID,
		TokenHash:  e.TokenHash,
		UserAgent:  e.UserAgent,
		IPAddress:  e.IPAddress,
		ExpiresAt:  e.ExpiresAt,
		Revoked:    e.Revoked,
		CreatedAt:  e.CreatedAt,
		LastActive: e.LastActive,
	}
}