// internal/domain/entity/payout.go
package entity

import "time"

const (
	PayoutTypePartnerSale         = "partner_sale"
	PayoutTypeAffiliateCommission = "affiliate_commission"
	PayoutStatusPending           = "pending"
	PayoutStatusPaid              = "paid"
)

type Payout struct {
	ID          uint       `json:"id"`
	RecipientID uint       `json:"recipient_id"`
	Type        string     `json:"type"`
	AmountCents int64      `json:"amount_cents"`
	Status      string     `json:"status"`
	OrderID     *uint      `json:"order_id,omitempty"`
	Reference   string     `json:"reference,omitempty"`
	PaidBy      *uint      `json:"paid_by,omitempty"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (p *Payout) AmountUSD() float64 { return float64(p.AmountCents) / 100 }