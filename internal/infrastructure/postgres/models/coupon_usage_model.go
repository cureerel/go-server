// internal/infrastructure/postgres/models/coupon_usage_model.go
package models

import (
	"time"

	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type CouponUsage struct {
	ID                   uint      `gorm:"primaryKey"`
	CouponID             uint      `gorm:"not null;index"`
	OrderID              uint      `gorm:"not null;index"`
	UserID               uint      `gorm:"not null;index"`
	DiscountAppliedCents int64     `gorm:"column:discount_applied_cents;default:0"`
	CommissionUSDCents   int64     `gorm:"column:commission_usd_cents;default:0"`
	CreatedAt            time.Time
}

func (CouponUsage) TableName() string { return "coupon_usages" }

func (m *CouponUsage) ToDomain() *entity.CouponUsage {
	return &entity.CouponUsage{
		ID:                   m.ID,
		CouponID:             m.CouponID,
		OrderID:              m.OrderID,
		UserID:               m.UserID,
		DiscountAppliedCents: m.DiscountAppliedCents,
		CommissionUSDCents:   m.CommissionUSDCents,
		CreatedAt:            m.CreatedAt,
	}
}

func CouponUsageFromDomain(e *entity.CouponUsage) *CouponUsage {
	return &CouponUsage{
		ID:                   e.ID,
		CouponID:             e.CouponID,
		OrderID:              e.OrderID,
		UserID:               e.UserID,
		DiscountAppliedCents: e.DiscountAppliedCents,
		CommissionUSDCents:   e.CommissionUSDCents,
		CreatedAt:            e.CreatedAt,
	}
}