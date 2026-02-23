package entity

import "time"

type PaymentStatus string
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
    ID            string          `json:"id"`
    UserID        uint            `json:"user_id"`
    OrderID       string          `json:"order_id"`
    Amount        int64           `json:"amount"`
    Currency      Currency        `json:"currency"`
    Status        PaymentStatus   `json:"status"`
    Provider      PaymentProvider `json:"provider"`
    ProviderTxnID string          `json:"provider_txn_id"`
    CustomerEmail string          `json:"customer_email"`
    Description   string          `json:"description"`
    RefundedAt    *time.Time      `json:"refunded_at,omitempty"`
    CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
}