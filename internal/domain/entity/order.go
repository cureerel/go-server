// internal/domain/entity/order.go
package entity

import "time"

// Order lifecycle
type OrderStatus string

const (
	OrderInCart   OrderStatus = "in_cart"
	OrderPaid     OrderStatus = "paid"
	OrderRefunded OrderStatus = "refunded"
)

// Delivery lifecycle
type DeliveryStatus string

const (
	DeliveryCreated    DeliveryStatus = "created"
	DeliveryInProgress DeliveryStatus = "in_progress"
	DeliveryPending    DeliveryStatus = "pending"
	DeliveryCompleted  DeliveryStatus = "completed"
	DeliveryReview     DeliveryStatus = "review" // customer left review
)

type Order struct {
	ID             uint           `json:"id"`
	UserID         uint           `json:"user_id"`
	Status         OrderStatus    `json:"status"`
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
	TotalCents     int64          `json:"total_cents"`
	Currency       string         `json:"currency"`
	CouponID       *uint          `json:"coupon_id,omitempty"`
	AffiliateID    *uint          `json:"affiliate_id,omitempty"`
	PaymentID      string         `json:"payment_id,omitempty"`
	Items          []OrderItem    `json:"items"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type OrderItem struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	ProductID *uint  `json:"product_id,omitempty"`
	ServiceID *uint  `json:"service_id,omitempty"`
	Title     string `json:"title"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}

func (o *Order) TotalUSD() float64 { return float64(o.TotalCents) / 100 }
func (o *Order) IsInCart() bool    { return o.Status == OrderInCart }
func (o *Order) IsPaid() bool      { return o.Status == OrderPaid }
func (o *Order) IsRefunded() bool  { return o.Status == OrderRefunded }

func (o *Order) CalculateTotal() int64 {
	var t int64
	for _, item := range o.Items {
		t += item.UnitPrice * int64(item.Quantity)
	}
	return t
}
