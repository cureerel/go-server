// Package stripe implements the pg.Gateway interface for Stripe Checkout.
// Import this package for its side-effects to register the provider:
//
//	import _ "github.com/cureerel/cserver/internal/infrastructure/pg/stripe"
package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cureerel/cserver/internal/infrastructure/pg"
	gosdk "github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/webhook"
)

func init() { pg.Register(New()) }

// client implements pg.Gateway for Stripe.
type client struct {
	secretKey     string
	webhookSecret string
	appURL        string
}

// New builds a Stripe client from environment variables.
// STRIPE_SECRET_KEY, STRIPE_WEBHOOK_SECRET, APP_URL
func New() pg.Gateway {
	key := os.Getenv("STRIPE_SECRET_KEY")
	gosdk.Key = key
	return &client{
		secretKey:     key,
		webhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		appURL:        os.Getenv("APP_URL"),
	}
}

func (c *client) Name() string { return "stripe" }

func (c *client) CreateCheckout(_ context.Context, in pg.CheckoutInput) (*pg.CheckoutResult, error) {
	if c.secretKey == "" {
		return nil, fmt.Errorf("stripe: STRIPE_SECRET_KEY is not set")
	}
	params := &gosdk.CheckoutSessionParams{
		Mode: gosdk.String(string(gosdk.CheckoutSessionModePayment)),
		LineItems: []*gosdk.CheckoutSessionLineItemParams{
			{
				Quantity: gosdk.Int64(1),
				PriceData: &gosdk.CheckoutSessionLineItemPriceDataParams{
					Currency:   gosdk.String(in.Currency),
					UnitAmount: gosdk.Int64(in.AmountMinor),
					ProductData: &gosdk.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: gosdk.String(in.Title),
					},
				},
			},
		},
		SuccessURL: gosdk.String(in.SuccessURL),
		CancelURL:  gosdk.String(in.CancelURL),
	}
	for k, v := range in.Metadata {
		params.AddMetadata(k, v)
	}
	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}
	return &pg.CheckoutResult{URL: sess.URL, ProviderRefID: sess.ID}, nil
}

func (c *client) VerifyWebhook(payload []byte, signature string) (*pg.WebhookEvent, error) {
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		c.webhookSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		return nil, fmt.Errorf("stripe: invalid webhook signature: %w", err)
	}

	if event.Type != "checkout.session.completed" {
		return &pg.WebhookEvent{Provider: "stripe", EventType: string(event.Type), Raw: payload}, nil
	}

	var sess struct {
		PaymentStatus string            `json:"payment_status"`
		Status        string            `json:"status"`
		Metadata      map[string]string `json:"metadata"`
	}
	if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
		return nil, fmt.Errorf("stripe: failed to parse session: %w", err)
	}
	if sess.PaymentStatus != "paid" && sess.Status != "complete" {
			return &pg.WebhookEvent{Provider: "stripe", EventType: string(event.Type), Raw: payload}, nil
	}

	we := &pg.WebhookEvent{
		Provider:     "stripe",
		EventType:    string(event.Type),
		Raw:          payload,
		Plan:         sess.Metadata["plan"],
		PurchaseType: sess.Metadata["purchase_type"],
	}
	fmt.Sscanf(sess.Metadata["user_id"], "%d", &we.UserID)
	if we.PurchaseType == "coins" {
		fmt.Sscanf(we.Plan, "%d", &we.Coins)
	}
	return we, nil
}
