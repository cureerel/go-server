// internal/domain/repository/repository.go
package repository

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
)





// Coupon
type CouponRepository interface {
	Create(ctx context.Context, coupon *entity.Coupon) error
	GetByID(ctx context.Context, id uint) (*entity.Coupon, error)
	GetByCode(ctx context.Context, code string) (*entity.Coupon, error)
	GetAll(ctx context.Context, page, limit int, status string) ([]entity.Coupon, int64, error)
	UpdateStatus(ctx context.Context, id uint, status string, approvedBy uint) error
	IncrementUsed(ctx context.Context, id uint) error
}

// Coupon usage
type CouponUsageRepository interface {
	Create(ctx context.Context, usage *entity.CouponUsage) error
	GetByOrderID(ctx context.Context, orderID uint) (*entity.CouponUsage, error)
	GetByCouponID(ctx context.Context, couponID uint, page, limit int) ([]entity.CouponUsage, int64, error)
}
