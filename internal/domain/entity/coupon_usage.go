// internal/domain/entity/coupon_usage.go
package entity

import "time"

type CouponUsage struct {
	ID                   uint      `json:"id"`
	CouponID             uint      `json:"coupon_id"`
	OrderID              uint      `json:"order_id"`
	UserID               uint      `json:"user_id"`
	DiscountAppliedCents int64     `json:"discount_applied_cents"`
	CommissionUSDCents   int64     `json:"commission_usd_cents"`
	CreatedAt            time.Time `json:"created_at"`
}