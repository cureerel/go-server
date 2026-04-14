// internal/domain/entity/coupon.go
package entity

import "time"

const (
	CouponTypeDiscount  = "discount"
	CouponTypeAffiliate = "affiliate"
	CouponTypeBoth      = "both"

	CouponStatusPending  = "pending"
	CouponStatusApproved = "approved"
	CouponStatusRejected = "rejected"
)

type Coupon struct {
	ID               uint       `json:"id"`
	CreatorID        uint       `json:"creator_id"`
	Code             string     `json:"code"`
	Type             string     `json:"type"`
	DiscountUSDCents int64      `json:"discount_usd_cents"`
	MaxDiscountCents int64      `json:"max_discount_cents"`
	CommissionPct    float64    `json:"commission_pct"`
	Status           string     `json:"status"`
	UsageLimit       *int       `json:"usage_limit,omitempty"`
	UsedCount        int        `json:"used_count"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	ApprovedBy       *uint      `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}


type CouponUsage struct {
	ID                   uint      `json:"id"`
	CouponID             uint      `json:"coupon_id"`
	OrderID              uint      `json:"order_id"`
	UserID               uint      `json:"user_id"`
	DiscountAppliedCents int64     `json:"discount_applied_cents"`
	CommissionUSDCents   int64     `json:"commission_usd_cents"`
	CreatedAt            time.Time `json:"created_at"`
}



func (c *Coupon) IsValid() bool {
	if c.Status != CouponStatusApproved {
		return false
	}
	if c.ExpiresAt != nil && time.Now().After(*c.ExpiresAt) {
		return false
	}
	if c.UsageLimit != nil && c.UsedCount >= *c.UsageLimit {
		return false
	}
	return true
}

func (c *Coupon) ApplyDiscount(amountCents int64) int64 {
	d := c.DiscountUSDCents
	if d > c.MaxDiscountCents {
		d = c.MaxDiscountCents
	}
	if d > amountCents {
		d = amountCents
	}
	return d
}

