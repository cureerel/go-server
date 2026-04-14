// Package dodo implements the pg.Gateway interface for Dodo Payments.
// Dodo is a crypto-capable, global payment rail alternative to Stripe.
// Import this package for its side-effects to register the provider:
//
//	import _ "github.com/cureerel/cserver/internal/infrastructure/pg/dodo"
package dodo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cureerel/cserver/internal/infrastructure/pg"
)

func init() { pg.Register(New()) }

const defaultBaseURL = "https://api.dodopayments.com/v1"

// client implements pg.Gateway for Dodo Payments.
type client struct {
	apiKey        string
	webhookSecret string
	baseURL       string
	appURL        string
	http          *http.Client
}

// New builds a Dodo client from environment variables.
// DODOPAYMENTS_API_KEY, DODOPAYMENTS_WEBHOOK_SECRET, DODOPAYMENTS_BASE_URL (optional), APP_URL
func New() pg.Gateway {
	base := os.Getenv("DODOPAYMENTS_BASE_URL")
	if base == "" {
		base = defaultBaseURL
	}
	return &client{
		apiKey:        os.Getenv("DODOPAYMENTS_API_KEY"),
		webhookSecret: os.Getenv("DODOPAYMENTS_WEBHOOK_SECRET"),
		baseURL:       base,
		appURL:        os.Getenv("APP_URL"),
		http:          &http.Client{},
	}
}

func (c *client) Name() string { return "dodo" }

// CreateCheckout creates a Dodo hosted checkout session.
func (c *client) CreateCheckout(ctx context.Context, in pg.CheckoutInput) (*pg.CheckoutResult, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("dodo: DODOPAYMENTS_API_KEY is not set")
	}

	body := map[string]any{
		"amount":      in.AmountMinor,
		"currency":    in.Currency,
		"description": in.Title,
		"receipt":     in.Receipt,
		"success_url": in.SuccessURL,
		"cancel_url":  in.CancelURL,
		"metadata":    in.Metadata,
	}
	raw, err := c.post(ctx, "/checkout/sessions", body)
	if err != nil {
		return nil, fmt.Errorf("dodo: CreateCheckout failed: %w", err)
	}

	result := &pg.CheckoutResult{}
	if u, ok := raw["url"].(string); ok {
		result.URL = u
	}
	if id, ok := raw["id"].(string); ok {
		result.ProviderRefID = id
	}
	return result, nil
}

// VerifyWebhook validates the Dodo-Signature header and parses the event.
// Dodo uses a simple bearer-style secret comparison (upgrade to HMAC when available).
func (c *client) VerifyWebhook(payload []byte, signature string) (*pg.WebhookEvent, error) {
	if c.webhookSecret != "" && signature != c.webhookSecret {
		return nil, fmt.Errorf("dodo: invalid webhook signature")
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID       string            `json:"id"`
				Status   string            `json:"status"`
				Metadata map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("dodo: failed to parse webhook payload: %w", err)
	}

	we := &pg.WebhookEvent{
		Provider:  "dodo",
		EventType: event.Type,
		Raw:       payload,
	}

	if event.Type == "checkout.completed" && event.Data.Object.Status == "paid" {
		meta := event.Data.Object.Metadata
		we.Plan = meta["plan"]
		we.PurchaseType = meta["purchase_type"]
		fmt.Sscanf(meta["user_id"], "%d", &we.UserID)
		if we.PurchaseType == "coins" {
			fmt.Sscanf(we.Plan, "%d", &we.Coins)
		}
	}
	return we, nil
}

// ─── internal HTTP helpers ────────────────────────────────────────────────────

func (c *client) post(ctx context.Context, path string, body map[string]any) (map[string]any, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("dodo: HTTP %d: %s", resp.StatusCode, string(raw))
	}
	var result map[string]any
	return result, json.Unmarshal(raw, &result)
}
