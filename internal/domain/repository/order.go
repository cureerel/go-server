// internal/domain/repository/order.go
package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"gorm.io/gorm"
)

type OrderFilter struct {
	Page      int
	Limit     int
	UserID    *uint
	Status    string
	ServiceID *uint
}

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	CreateWithTx(ctx context.Context, tx *gorm.DB, order *entity.Order) error
	GetByID(ctx context.Context, id uint) (*entity.Order, error)
	GetAll(ctx context.Context, filter OrderFilter) ([]entity.Order, int64, error)
	UpdateStatus(ctx context.Context, id uint, status entity.OrderStatus) error
	AttachPaymentProvider(ctx context.Context, id uint, provider string) error
}