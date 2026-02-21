package entity

import "time"

type WebhookEvent struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"` // stripe, razorpay, etc.
	EventType string    `json:"event_type"`
	Payload   []byte    `json:"-"`
	Signature string    `json:"-"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"order_id"`
	Amount        int64     `json:"amount"` // cents
	Currency      string    `json:"currency"`
	Status        string    `json:"status"` // pending, completed, failed
	Provider      string    `json:"provider"`
	ProviderTxnID string    `json:"provider_txn_id"`
	CustomerEmail string    `json:"customer_email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}