// internal/interfaces/dto/order_dto.go
package dto

// Requests

type AddToCartRequest struct {
	ProductID *uint `json:"product_id"`
	ServiceID *uint `json:"service_id"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type CheckoutRequest struct {
	CouponCode  string `json:"coupon_code"`
	AffiliateID *uint  `json:"affiliate_id"`
}

type UpdateDeliveryStatusRequest struct {
	DeliveryStatus string `json:"delivery_status" binding:"required,oneof=created in_progress pending completed review"`
}

// Responses

type OrderItemResponse struct {
	ID        uint    `json:"id"`
	ProductID *uint   `json:"product_id,omitempty"`
	ServiceID *uint   `json:"service_id,omitempty"`
	Title     string  `json:"title"`
	Quantity  int     `json:"quantity"`
	UnitPrice int64   `json:"unit_price"`
	UnitUSD   float64 `json:"unit_usd"`
}

type OrderResponse struct {
	ID             uint                `json:"id"`
	UserID         uint                `json:"user_id"`
	Status         string              `json:"status"`
	DeliveryStatus string              `json:"delivery_status"`
	TotalCents     int64               `json:"total_cents"`
	TotalUSD       float64             `json:"total_usd"`
	Currency       string              `json:"currency"`
	PaymentID      string              `json:"payment_id,omitempty"`
	Items          []OrderItemResponse `json:"items"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

type OrderListResponse struct {
	Data  []OrderResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}
