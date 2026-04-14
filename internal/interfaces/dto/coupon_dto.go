// internal/interfaces/dto/coupon_dto.go
package dto

import "time"

type CreateCouponRequest struct {
	Code             string  `json:"code"              binding:"required,min=3,max=32"`
	Type             string  `json:"type"              binding:"required,oneof=discount affiliate both"`
	DiscountUSDCents int64   `json:"discount_usd_cents"`
	MaxDiscountCents int64   `json:"max_discount_cents"`
	CommissionPct    float64 `json:"commission_pct"`
	UsageLimit       *int    `json:"usage_limit"`
	ExpiresAt        *string `json:"expires_at"`
}

type ValidateCouponRequest struct {
	Code string `json:"code" binding:"required"`
}

type CouponResponse struct {
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
	CreatedAt        string     `json:"created_at"`
}

type CouponListResponse struct {
	Data  []CouponResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type ValidateCouponResponse struct {
	Valid            bool    `json:"valid"`
	Code             string  `json:"code"`
	Type             string  `json:"type"`
	DiscountUSDCents int64   `json:"discount_usd_cents"`
	CommissionPct    float64 `json:"commission_pct"`
}
