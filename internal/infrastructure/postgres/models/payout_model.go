// internal/infrastructure/postgres/models/payout_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Payout struct {
	ID          uint       `gorm:"primaryKey"`
	RecipientID uint       `gorm:"not null;index"`
	Type        string     `gorm:"not null;size:30;index"`
	AmountCents int64      `gorm:"column:amount_cents;not null"`
	Status      string     `gorm:"default:'pending';size:20;index"`
	OrderID     *uint      `gorm:"column:order_id"`
	Reference   string     `gorm:"type:text"`
	PaidBy      *uint      `gorm:"column:paid_by"`
	PaidAt      *time.Time `gorm:"column:paid_at"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Payout) TableName() string { return "payouts" }

func (m *Payout) ToDomain() *entity.Payout {
	return &entity.Payout{
		ID: m.ID, RecipientID: m.RecipientID, Type: m.Type,
		AmountCents: m.AmountCents, Status: m.Status, OrderID: m.OrderID,
		Reference: m.Reference, PaidBy: m.PaidBy, PaidAt: m.PaidAt,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func PayoutFromDomain(e *entity.Payout) *Payout {
	return &Payout{
		ID: e.ID, RecipientID: e.RecipientID, Type: e.Type,
		AmountCents: e.AmountCents, Status: e.Status, OrderID: e.OrderID,
		Reference: e.Reference, PaidBy: e.PaidBy, PaidAt: e.PaidAt,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}