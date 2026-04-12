// internal/infrastructure/postgres/models/payment_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Payment struct {
	ID                 string     `gorm:"primaryKey;size:64"`
	OrderID            uint       `gorm:"not null;index"`
	UserID             uint       `gorm:"not null;index"`
	AmountCents        int64      `gorm:"column:amount_cents;not null"`
	Currency           string     `gorm:"size:10;default:'USD'"`
	Status             string     `gorm:"default:'pending';size:20;index"`
	Provider           string     `gorm:"not null;size:20"`
	ProviderTxnID      string     `gorm:"column:provider_txn_id;size:128"`
	ProviderPaymentID  string     `gorm:"column:provider_payment_id;size:128"`
	CustomerEmail      string     `gorm:"column:customer_email;size:100"`
	Description        string     `gorm:"type:text"`
	RefundID           string     `gorm:"column:refund_id;size:128"`
	RefundedAt         *time.Time `gorm:"column:refunded_at"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (Payment) TableName() string { return "payments" }

func (m *Payment) ToDomain() *entity.Payment {
	return &entity.Payment{
		ID:                m.ID,
		OrderID:           m.OrderID,
		UserID:            m.UserID,
		AmountCents:       m.AmountCents,
		Currency:          m.Currency,
		Status:            entity.PaymentStatus(m.Status),
		Provider:          entity.PaymentProvider(m.Provider),
		ProviderTxnID:     m.ProviderTxnID,
		ProviderPaymentID: m.ProviderPaymentID,
		CustomerEmail:     m.CustomerEmail,
		Description:       m.Description,
		RefundID:          m.RefundID,
		RefundedAt:        m.RefundedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func PaymentFromDomain(e *entity.Payment) *Payment {
	return &Payment{
		ID:                e.ID,
		OrderID:           e.OrderID,
		UserID:            e.UserID,
		AmountCents:       e.AmountCents,
		Currency:          e.Currency,
		Status:            string(e.Status),
		Provider:          string(e.Provider),
		ProviderTxnID:     e.ProviderTxnID,
		ProviderPaymentID: e.ProviderPaymentID,
		CustomerEmail:     e.CustomerEmail,
		Description:       e.Description,
		RefundID:          e.RefundID,
		RefundedAt:        e.RefundedAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}