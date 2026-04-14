// internal/interfaces/dto/payment_dto.go
package dto

// ── Requests

type RefundRequest struct {
	RefundID string `json:"refund_id" binding:"required"`
}

// ── Responses

type PaymentResponse struct {
	ID            string  `json:"id"`
	OrderID       uint    `json:"order_id"`
	UserID        uint    `json:"user_id"`
	AmountCents   int64   `json:"amount_cents"`
	AmountUSD     float64 `json:"amount_usd"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	Provider      string  `json:"provider"`
	ProviderTxnID string  `json:"provider_txn_id,omitempty"`
	CustomerEmail string  `json:"customer_email,omitempty"`
	Description   string  `json:"description,omitempty"`
	RefundID      string  `json:"refund_id,omitempty"`
	RefundedAt    *string `json:"refunded_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type PaymentListResponse struct {
	Data  []PaymentResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}
