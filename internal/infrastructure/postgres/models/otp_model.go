// internal/infrastructure/postgres/models/otp_model.go
package models

import (
	"time"
	"github.com/cureerel/cserver/internal/domain/entity"
)

type OTP struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `gorm:"index"`
	Email     string    `gorm:"not null;size:100;index"`
	Code      string    `gorm:"not null;size:6"`
	Type      string    `gorm:"not null;size:20;index"`
	Used      bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (OTP) TableName() string { return "otps" }

func (m *OTP) ToDomain() *entity.OTP {
	return &entity.OTP{
		ID:        m.ID,
		UserID:    m.UserID,
		Email:     m.Email,
		Code:      m.Code,
		Type:      m.Type,
		Used:      m.Used,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

func OTPFromDomain(e *entity.OTP) *OTP {
	return &OTP{
		ID:        e.ID,
		UserID:    e.UserID,
		Email:     e.Email,
		Code:      e.Code,
		Type:      e.Type,
		Used:      e.Used,
		ExpiresAt: e.ExpiresAt,
		CreatedAt: e.CreatedAt,
	}
}