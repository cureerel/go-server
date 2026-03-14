// internal/domain/repository/coupon_usage.go
package repository

import (
	"context"

	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type CouponUsageRepository interface {
	Create(ctx context.Context, usage *entity.CouponUsage) error
	GetByOrderID(ctx context.Context, orderID uint) (*entity.CouponUsage, error)
	GetByCouponID(ctx context.Context, couponID uint, page, limit int) ([]entity.CouponUsage, int64, error)
}