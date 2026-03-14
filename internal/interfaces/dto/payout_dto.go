// internal/interfaces/dto/payout_dto.go
package dto

type MarkPaidRequest struct {
	Reference string `json:"reference" binding:"required"`
}

type PayoutResponse struct {
	ID          uint    `json:"id"`
	RecipientID uint    `json:"recipient_id"`
	Type        string  `json:"type"`
	AmountCents int64   `json:"amount_cents"`
	AmountUSD   float64 `json:"amount_usd"`
	Status      string  `json:"status"`
	OrderID     *uint   `json:"order_id,omitempty"`
	Reference   string  `json:"reference,omitempty"`
	PaidBy      *uint   `json:"paid_by,omitempty"`
	PaidAt      *string `json:"paid_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type PayoutListResponse struct {
	Data  []PayoutResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}
