// Package pg is the unified payment gateway layer.
// Any provider (Stripe, Razorpay, Dodo, …) implements the Gateway interface
// and self-registers in its init(). The handler layer uses the registry — it
// never imports a specific provider directly.
package pg

import (
	"context"
	"fmt"
	"os"
)

// ─── Core types ──────────────────────────────────────────────────────────────

// CheckoutInput is provider-agnostic checkout data.
type CheckoutInput struct {
	AmountMinor int64             // smallest currency unit (cents, paise …)
	Currency    string            // ISO 4217 e.g. "usd", "inr"
	Title       string            // line-item label shown on the checkout page
	SuccessURL  string            // full URL to redirect on success
	CancelURL   string            // full URL to redirect on cancel
	Receipt     string            // unique receipt / order ref for the provider
	Metadata    map[string]string // arbitrary KV forwarded to the provider
	Notes       map[string]any    // Razorpay-style notes (optional)
}

// CheckoutResult is what every provider returns after creating a checkout.
type CheckoutResult struct {
	URL           string // redirect URL (Stripe, Dodo)
	ProviderRefID string // provider order/session ID (Razorpay order_id, Stripe cs_…)
}

// WebhookEvent is a normalised representation of a confirmed payment event.
type WebhookEvent struct {
	Provider     string
	EventType    string
	UserID       uint
	Plan         string // "basic" | "pro" | "100" (coin pack) …
	PurchaseType string // "membership" | "coins"
	Coins        int64  // populated when PurchaseType == "coins"
	Raw          []byte
}

// ─── Gateway interface ────────────────────────────────────────────────────────

// Gateway is the contract every payment provider must satisfy.
type Gateway interface {
	// Name returns the canonical provider key used in the registry ("stripe", "razorpay", "dodo" …).
	Name() string

	// CreateCheckout initiates a payment session or order and returns a result
	// the API can forward to the client.
	CreateCheckout(ctx context.Context, in CheckoutInput) (*CheckoutResult, error)

	// VerifyWebhook validates the provider's webhook signature and parses the
	// event into a normalised WebhookEvent.  Returns an error if the signature
	// is invalid or the payload is malformed.
	VerifyWebhook(payload []byte, signature string) (*WebhookEvent, error)
}

// PaymentVerifier is an optional interface for two-step flows (e.g. Razorpay)
// where the frontend confirms payment before the server can consider it done.
type PaymentVerifier interface {
	VerifyPaymentSignature(providerOrderID, paymentID, signature string) bool
}

// ─── Registry ─────────────────────────────────────────────────────────────────

var registry = map[string]Gateway{}

// Register adds a provider to the global registry.
// Call this from each provider's init() function.
func Register(g Gateway) {
	registry[g.Name()] = g
}

// Get returns a registered Gateway by name.
func Get(name string) (Gateway, bool) {
	g, ok := registry[name]
	return g, ok
}

// MustGet panics if the provider is not registered.
func MustGet(name string) Gateway {
	g, ok := registry[name]
	if !ok {
		panic(fmt.Sprintf("pg: provider %q not registered", name))
	}
	return g
}

// All returns every registered gateway keyed by name.
func All() map[string]Gateway { return registry }

// ─── Factory / region-based selection ────────────────────────────────────────

type Region string

const (
	RegionIN     Region = "IN"
	RegionGlobal Region = "GLOBAL"
)

// Factory selects the appropriate provider based on deployment region and
// environment overrides.
type Factory struct {
	Region Region
	AppURL string
}

// FromEnv builds a Factory from PAYMENT_REGION and APP_URL env vars.
func FromEnv() *Factory {
	r := RegionGlobal
	if os.Getenv("PAYMENT_REGION") == "IN" {
		r = RegionIN
	}
	return &Factory{Region: r, AppURL: os.Getenv("APP_URL")}
}

// DefaultProvider returns the recommended provider name for one-time checkout:
//   - PAYMENT_PROVIDER env var (highest priority)
//   - IN region → "razorpay"
//   - Global   → "stripe"
func (f *Factory) DefaultProvider() string {
	if v := os.Getenv("PAYMENT_PROVIDER"); v != "" {
		return v
	}
	if f.Region == RegionIN {
		return "razorpay"
	}
	return "stripe"
}

// DefaultTopupProvider returns the preferred provider for coin top-ups.
// Can be overridden with PAYMENT_TOPUP_PROVIDER.
func (f *Factory) DefaultTopupProvider() string {
	if v := os.Getenv("PAYMENT_TOPUP_PROVIDER"); v != "" {
		return v
	}
	return f.DefaultProvider()
}

// ResolveCheckout picks the Gateway for a checkout, falling back to the default.
func (f *Factory) ResolveCheckout(preferredProvider string) (Gateway, error) {
	name := preferredProvider
	if name == "" {
		name = f.DefaultProvider()
	}
	g, ok := Get(name)
	if !ok {
		return nil, fmt.Errorf("pg: provider %q not registered — import its package for side-effects", name)
	}
	return g, nil
}
