package entity

import (
    "fmt"
    "time"
)

type Currency string
type ProductType string

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

const (
    ProductPhysical ProductType = "physical"
    ProductDigital  ProductType = "digital"
)

type Product struct {
    ID          uint        `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Type        ProductType `json:"type"`
    Price       int64       `json:"price"`
    Currency    Currency    `json:"currency"`
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