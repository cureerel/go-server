// internal/infrastructure/postgres/models/upgrade_request_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type UpgradeRequest struct {
	ID         uint       `gorm:"primaryKey"`
	UserID     uint       `gorm:"not null;index"`
	FromRole   string     `gorm:"column:from_role;not null;size:20"`
	ToRole     string     `gorm:"column:to_role;not null;size:20"`
	Status     string     `gorm:"default:'pending';size:20;index"`
	ReviewedBy *uint      `gorm:"column:reviewed_by"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at"`
	CreatedAt  time.Time
}

func (UpgradeRequest) TableName() string { return "upgrade_requests" }

func (m *UpgradeRequest) ToDomain() *entity.UpgradeRequest {
	return &entity.UpgradeRequest{
		ID: m.ID, UserID: m.UserID,
		FromRole: m.FromRole, ToRole: m.ToRole, Status: m.Status,
		ReviewedBy: m.ReviewedBy, ReviewedAt: m.ReviewedAt, CreatedAt: m.CreatedAt,
	}
}

func UpgradeRequestFromDomain(e *entity.UpgradeRequest) *UpgradeRequest {
	return &UpgradeRequest{
		ID: e.ID, UserID: e.UserID,
		FromRole: e.FromRole, ToRole: e.ToRole, Status: e.Status,
		ReviewedBy: e.ReviewedBy, ReviewedAt: e.ReviewedAt, CreatedAt: e.CreatedAt,
	}
}