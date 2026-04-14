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

// IsValid returns true when the coupon can be applied at checkout.
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

// ApplyDiscount returns the discount amount (capped at order total and MaxDiscountCents).
func (c *Coupon) ApplyDiscount(amountCents int64) int64 {
	if c.Type == CouponTypeAffiliate {
		return 0 // affiliate coupons give commission, not buyer discount
	}
	d := c.DiscountUSDCents
	if c.MaxDiscountCents > 0 && d > c.MaxDiscountCents {
		d = c.MaxDiscountCents
	}
	if d > amountCents {
		d = amountCents
	}
	return d
}

// CommissionAmount returns how many cents the coupon creator earns on this order.
func (c *Coupon) CommissionAmount(orderCents int64) int64 {
	if c.CommissionPct <= 0 || (c.Type == CouponTypeDiscount) {
		return 0
	}
	return int64(float64(orderCents) * c.CommissionPct / 100)
}
