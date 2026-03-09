// internal/domain/entity/payment.go
package entity

import "time"

type PaymentStatus   string
type PaymentProvider string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

const (
	ProviderStripe   PaymentProvider = "stripe"
	ProviderRazorpay PaymentProvider = "razorpay"
)

type Payment struct {
	ID                 string          `json:"id"`
	OrderID            uint            `json:"order_id"`
	UserID             uint            `json:"user_id"`
	AmountCents        int64           `json:"amount_cents"`
	Currency           string          `json:"currency"`
	Status             PaymentStatus   `json:"status"`
	Provider           PaymentProvider `json:"provider"`
	ProviderTxnID      string          `json:"provider_txn_id,omitempty"`
	ProviderPaymentID  string          `json:"provider_payment_id,omitempty"`
	CustomerEmail      string          `json:"customer_email,omitempty"`
	Description        string          `json:"description,omitempty"`
	RefundID           string          `json:"refund_id,omitempty"`
	RefundedAt         *time.Time      `json:"refunded_at,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func (p *Payment) AmountUSD() float64  { return float64(p.AmountCents) / 100 }
func (p *Payment) IsPending() bool     { return p.Status == PaymentPending }
func (p *Payment) IsCompleted() bool   { return p.Status == PaymentCompleted }
func (p *Payment) IsRefunded() bool    { return p.Status == PaymentRefunded }