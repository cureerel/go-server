// internal/infrastructure/postgres/models/coupon_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Coupon struct {
	ID               uint       `gorm:"primaryKey"`
	CreatorID        uint       `gorm:"not null;index"`
	Code             string     `gorm:"uniqueIndex;not null;size:32"`
	Type             string     `gorm:"not null;size:20"`
	DiscountUSDCents int64      `gorm:"column:discount_usd_cents;default:0"`
	MaxDiscountCents int64      `gorm:"column:max_discount_cents;default:10000"`
	CommissionPct    float64    `gorm:"column:commission_pct;default:0"`
	Status           string     `gorm:"default:'pending';size:20;index"`
	UsageLimit       *int       `gorm:"column:usage_limit"`
	UsedCount        int        `gorm:"column:used_count;default:0"`
	ExpiresAt        *time.Time `gorm:"column:expires_at"`
	ApprovedBy       *uint      `gorm:"column:approved_by"`
	ApprovedAt       *time.Time `gorm:"column:approved_at"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (Coupon) TableName() string { return "coupons" }

func (m *Coupon) ToDomain() *entity.Coupon {
	return &entity.Coupon{
		ID: m.ID, CreatorID: m.CreatorID, Code: m.Code,
		Type: m.Type, DiscountUSDCents: m.DiscountUSDCents,
		MaxDiscountCents: m.MaxDiscountCents, CommissionPct: m.CommissionPct,
		Status: m.Status, UsageLimit: m.UsageLimit, UsedCount: m.UsedCount,
		ExpiresAt: m.ExpiresAt, ApprovedBy: m.ApprovedBy, ApprovedAt: m.ApprovedAt,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func CouponFromDomain(e *entity.Coupon) *Coupon {
	return &Coupon{
		ID: e.ID, CreatorID: e.CreatorID, Code: e.Code,
		Type: e.Type, DiscountUSDCents: e.DiscountUSDCents,
		MaxDiscountCents: e.MaxDiscountCents, CommissionPct: e.CommissionPct,
		Status: e.Status, UsageLimit: e.UsageLimit, UsedCount: e.UsedCount,
		ExpiresAt: e.ExpiresAt, ApprovedBy: e.ApprovedBy, ApprovedAt: e.ApprovedAt,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

type CouponUsage struct {
	ID                   uint  `gorm:"primaryKey"`
	CouponID             uint  `gorm:"not null;index"`
	OrderID              uint  `gorm:"not null;index"`
	UserID               uint  `gorm:"not null;index"`
	DiscountAppliedCents int64 `gorm:"column:discount_applied_cents;default:0"`
	CommissionUSDCents   int64 `gorm:"column:commission_usd_cents;default:0"`
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
