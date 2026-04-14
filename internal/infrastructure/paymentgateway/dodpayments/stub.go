// Package dodpayments is a placeholder for DOD Payments (or similar) global/crypto card rails.
package dodpayments

import (
	"errors"
	"os"
)

var ErrNotImplemented = errors.New("dodpayments: configure DODPAYMENTS_API_KEY and implement CreateCheckout in stub.go")

type Client struct {
	APIKey  string
	BaseURL string
}

func NewFromEnv() (*Client, error) {
	k := os.Getenv("DODPAYMENTS_API_KEY")
	if k == "" {
		return nil, ErrNotImplemented
	}
	return &Client{APIKey: k, BaseURL: os.Getenv("DODPAYMENTS_BASE_URL")}, nil
}

// CreateCheckout is a stub 
func (c *Client) CreateCheckout(amountMinor int64, currency, receipt string, metadata map[string]string) (map[string]any, error) {
	if c == nil || c.APIKey == "" {
		return nil, ErrNotImplemented
	}
	return nil, ErrNotImplemented
}
