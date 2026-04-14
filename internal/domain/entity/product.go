// internal/domain/entity/product.go
package entity

import (
	"fmt"
	"time"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
	CurrencyINR Currency = "INR"
	CurrencyJPY Currency = "JPY"
	CurrencyAUD Currency = "AUD"
	CurrencyCAD Currency = "CAD"
	CurrencySGD Currency = "SGD"
	CurrencyAED Currency = "AED"
)

type ProductType string

const (
	ProductPhysical ProductType = "physical"
	ProductDigital  ProductType = "digital"
)

type Product struct {
	ID          uint        `json:"id"`
	SKU         string      `json:"sku"` // unique human-readable slug (e.g. prd_lx4k9mz)
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        ProductType `json:"type"`
	Price       int64       `json:"price"` // smallest currency unit (cents/paise)
	Currency    Currency    `json:"currency"`
	Stock       int         `json:"stock"` // -1 = unlimited (digital)
	ImageURL    string      `json:"image_url,omitempty"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (p *Product) FormattedPrice() string {
	switch p.Currency {
	case CurrencyJPY:
		return fmt.Sprintf("%s %.0f", string(p.Currency), float64(p.Price))
	default:
		return fmt.Sprintf("%s %.2f", string(p.Currency), float64(p.Price)/100)
	}
}

func (p *Product) IsDigital() bool   { return p.Type == ProductDigital }
func (p *Product) InStock() bool     { return p.Stock == -1 || p.Stock > 0 }
func (p *Product) PriceUSD() float64 { return float64(p.Price) / 100 }
