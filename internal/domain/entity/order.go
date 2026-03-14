// internal/domain/entity/order.go
package entity

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderCancelled  OrderStatus = "cancelled"
	OrderCompleted  OrderStatus = "completed"
)

type Order struct {
	ID              uint        `json:"id"`
	UserID          uint        `json:"user_id"`
	ServiceID       *uint       `json:"service_id,omitempty"`
	Status          OrderStatus `json:"status"`
	TotalCents      int64       `json:"total_cents"`
	Currency        string      `json:"currency"`
	PaymentProvider string      `json:"payment_provider,omitempty"`
	CouponID        *uint       `json:"coupon_id,omitempty"`
	AffiliateID     *uint       `json:"affiliate_id,omitempty"`
	Items           []OrderItem `json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	ServiceID *uint  `json:"service_id,omitempty"`
	Title     string `json:"title"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
}

func (o *Order) TotalUSD() float64 { return float64(o.TotalCents) / 100 }
func (o *Order) IsPending() bool   { return o.Status == OrderPending }
func (o *Order) IsCompleted() bool { return o.Status == OrderCompleted }
func (o *Order) IsCancelled() bool { return o.Status == OrderCancelled }

func (o *Order) CalculateTotal() int64 {
	var t int64
	for _, item := range o.Items {
		t += item.UnitPrice * int64(item.Quantity)
	}
	return t
}