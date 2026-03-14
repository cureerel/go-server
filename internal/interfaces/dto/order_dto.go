// internal/interfaces/dto/order_dto.go
package dto

// ── Requests ──────────────────────────────────────────────────

type CreateOrderRequest struct {
	ServiceID   uint   `json:"service_id"   binding:"required,min=1"`
	Provider    string `json:"provider"     binding:"required,oneof=stripe razorpay"`
	CouponCode  string `json:"coupon_code"`  
	AffiliateID *uint  `json:"affiliate_id"` 
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
}

// ── Responses ─────────────────────────────────────────────────

type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ServiceID *uint   `json:"service_id,omitempty"`
	Title     string  `json:"title"`
	Quantity  int     `json:"quantity"`
	UnitPrice int64   `json:"unit_price"`
	UnitUSD   float64 `json:"unit_usd"`
}

type OrderResponse struct {
	ID              uint                `json:"id"`
	UserID          uint                `json:"user_id"`
	ServiceID       *uint               `json:"service_id,omitempty"`
	Status          string              `json:"status"`
	TotalCents      int64               `json:"total_cents"`
	TotalUSD        float64             `json:"total_usd"`
	Currency        string              `json:"currency"`
	PaymentProvider string              `json:"payment_provider,omitempty"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

type OrderListResponse struct {
	Data  []OrderResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}