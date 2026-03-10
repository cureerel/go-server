// internal/domain/entity/upgrade_request.go
package entity

import "time"

const (
	UpgradeStatusPending  = "pending"
	UpgradeStatusApproved = "approved"
	UpgradeStatusRejected = "rejected"
)

type UpgradeRequest struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	FromRole   string     `json:"from_role"`
	ToRole     string     `json:"to_role"`
	Status     string     `json:"status"`
	ReviewedBy *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}