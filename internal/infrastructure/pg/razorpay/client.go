// Package razorpay implements the pg.Gateway interface for Razorpay.
// Also implements pg.PaymentVerifier for the 2-step frontend confirmation flow.
// Import this package for its side-effects to register the provider:
//
//	import _ "github.com/cureerel/cserver/internal/infrastructure/pg/razorpay"
package razorpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cureerel/cserver/internal/infrastructure/pg"
	rzpsdk "github.com/razorpay/razorpay-go"
)

func init() { pg.Register(New()) }

// client implements pg.Gateway + pg.PaymentVerifier for Razorpay.
type client struct {
	keyID         string
	keySecret     string
	webhookSecret string
	sdk           *rzpsdk.Client
}

// New builds a Razorpay client from environment variables.
// RAZORPAY_KEY_ID, RAZORPAY_KEY_SECRET, RAZORPAY_WEBHOOK_SECRET
func New() pg.Gateway {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	return &client{
		keyID:         keyID,
		keySecret:     keySecret,
		webhookSecret: os.Getenv("RAZORPAY_WEBHOOK_SECRET"),
		sdk:           rzpsdk.NewClient(keyID, keySecret),
	}
}

func (c *client) Name() string { return "razorpay" }

// CreateCheckout creates a Razorpay order (step 1 of the 2-step flow).
// The frontend uses ProviderRefID (order_id) to open the Razorpay checkout modal.
func (c *client) CreateCheckout(_ context.Context, in pg.CheckoutInput) (*pg.CheckoutResult, error) {
	if c.keyID == "" || c.keySecret == "" {
		return nil, fmt.Errorf("razorpay: RAZORPAY_KEY_ID / RAZORPAY_KEY_SECRET not set")
	}
	data := map[string]interface{}{
		"amount":   in.AmountMinor,
		"currency": in.Currency,
		"receipt":  in.Receipt,
		"notes":    in.Notes,
	}
	order, err := c.sdk.Order.Create(data, nil)
	if err != nil {
		return nil, fmt.Errorf("razorpay: failed to create order: %w", err)
	}
	id, _ := order["id"].(string)
	return &pg.CheckoutResult{ProviderRefID: id}, nil
}

// VerifyWebhook validates X-Razorpay-Signature and parses a payment.captured event.
func (c *client) VerifyWebhook(payload []byte, signature string) (*pg.WebhookEvent, error) {
	if !c.hmacMatch(c.webhookSecret, payload, signature) {
		return nil, fmt.Errorf("razorpay: invalid webhook signature")
	}

	var event struct {
		Event   string `json:"event"`
		Payload struct {
			Payment struct {
				Entity struct {
					Notes struct {
						UserID       uint   `json:"user_id"`
						Plan         string `json:"plan"`
						PurchaseType string `json:"purchase_type"`
					} `json:"notes"`
				} `json:"entity"`
			} `json:"payment"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("razorpay: failed to parse webhook payload: %w", err)
	}

	we := &pg.WebhookEvent{
		Provider:  "razorpay",
		EventType: event.Event,
		Raw:       payload,
	}
	if event.Event == "payment.captured" {
		n := event.Payload.Payment.Entity.Notes
		we.UserID = n.UserID
		we.Plan = n.Plan
		we.PurchaseType = n.PurchaseType
		if we.PurchaseType == "coins" {
			fmt.Sscanf(we.Plan, "%d", &we.Coins)
		}
	}
	return we, nil
}

// VerifyPaymentSignature verifies the HMAC returned after a successful Razorpay payment.
// Implements pg.PaymentVerifier.
func (c *client) VerifyPaymentSignature(providerOrderID, paymentID, signature string) bool {
	return c.hmacMatch(c.keySecret, []byte(providerOrderID+"|"+paymentID), signature)
}

// KeyID exposes the public key for the frontend (needed to open the checkout modal).
func (c *client) KeyID() string { return c.keyID }

func (c *client) hmacMatch(secret string, payload []byte, expected string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil)) == expected
}
