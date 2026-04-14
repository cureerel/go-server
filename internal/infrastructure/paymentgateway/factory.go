// Package paymentgateway configures checkout providers by region (India vs global).
package paymentgateway

import "os"

type Region string

const (
	RegionIN     Region = "IN"
	RegionGlobal Region = "GLOBAL"
)

type Factory struct {
	Region Region
}

func FromEnv() *Factory {
	r := RegionGlobal
	if v := os.Getenv("PAYMENT_REGION"); v == "IN" {
		r = RegionIN
	}
	return &Factory{Region: r}
}

// DefaultTopupProvider returns razorpay for India, stripe for global, unless PAYMENT_TOPUP_PROVIDER is set
// (razorpay|stripe|dodpayments).
func (f *Factory) DefaultTopupProvider() string {
	if v := os.Getenv("PAYMENT_TOPUP_PROVIDER"); v != "" {
		return v
	}
	if f.Region == RegionIN {
		return "razorpay"
	}
	return "stripe"
}
