// internal/domain/repository/payment.go
package repository

import (
	"context"

	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type PaymentFilter struct {
	Page    int
	Limit   int
	UserID  *uint
	OrderID *uint
	Status  string
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id string) (*entity.Payment, error)
	GetByOrderID(ctx context.Context, orderID uint) (*entity.Payment, error)
	GetAll(ctx context.Context, filter PaymentFilter) ([]entity.Payment, int64, error)
	UpdateStatus(ctx context.Context, id string, status entity.PaymentStatus) error
	MarkRefunded(ctx context.Context, id string, refundID string) error
	SetProviderTxnID(ctx context.Context, id string, txnID string) error
}