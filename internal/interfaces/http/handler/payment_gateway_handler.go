// internal/interfaces/http/handler/payment_gateway_handler.go
package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/infrastructure/pg"

	// Register providers via side-effect imports.
	// Add/remove providers here to enable or disable them.
	_ "github.com/cureerel/cserver/internal/infrastructure/pg/dodo"
	_ "github.com/cureerel/cserver/internal/infrastructure/pg/razorpay"
	_ "github.com/cureerel/cserver/internal/infrastructure/pg/stripe"

	"github.com/gin-gonic/gin"
)

// planAmountsPaise maps membership plan → paise (₹ × 100)
var planAmountsPaise = map[string]int64{
	"basic": 49900,
	"pro":   149900,
}

// planEntityMap maps plan string → entity.MembershipPlan
var planEntityMap = map[string]entity.MembershipPlan{
	"free":  entity.PlanFree,
	"basic": entity.PlanBasic,
	"pro":   entity.PlanPro,
}

// ─────────────────────────────────────────────────────────────────────────────

type PaymentGatewayHandler struct {
	membershipSvc *service.MembershipService
	coinSvc       *service.CoinService
	paymentSvc    *service.PaymentService
	pgFactory     *pg.Factory
}

func NewPaymentGatewayHandler(
	membershipSvc *service.MembershipService,
	coinSvc *service.CoinService,
	paymentSvc *service.PaymentService,
) *PaymentGatewayHandler {
	return &PaymentGatewayHandler{
		membershipSvc: membershipSvc,
		coinSvc:       coinSvc,
		paymentSvc:    paymentSvc,
		pgFactory:     pg.FromEnv(),
	}
}


// ── POST /api/payments/razorpay/create-order ─────────────────────────────────
// Body: { "plan": "basic"|"pro"|"100", "purchase_type": "membership"|"coins" }
func (h *PaymentGatewayHandler) RazorpayCreateOrder(c *gin.Context) {
	gw, err := h.pgFactory.ResolveCheckout("razorpay")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	h.createCheckout(c, gw)
}

// ── POST /api/payments/razorpay/verify ───────────────────────────────────────
func (h *PaymentGatewayHandler) RazorpayVerify(c *gin.Context) {
	var body struct {
		Plan              string `json:"plan"                binding:"required"`
		PurchaseType      string `json:"purchase_type"`
		RazorpayOrderID   string `json:"razorpay_order_id"   binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	gw, err := h.pgFactory.ResolveCheckout("razorpay")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	verifier, ok := gw.(pg.PaymentVerifier)
	if !ok || !verifier.VerifyPaymentSignature(body.RazorpayOrderID, body.RazorpayPaymentID, body.RazorpaySignature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature mismatch"})
		return
	}

	pt := body.PurchaseType
	if pt == "" {
		pt = "membership"
	}
	uid, _ := getUID(c)
	h.fulfil(c, uid, body.Plan, pt)
}

// ── POST /api/payments/stripe/create-session ─────────────────────────────────
func (h *PaymentGatewayHandler) StripeCreateSession(c *gin.Context) {
	gw, err := h.pgFactory.ResolveCheckout("stripe")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	h.createCheckout(c, gw)
}

// ── POST /api/payments/stripe/webhook ────────────────────────────────────────
func (h *PaymentGatewayHandler) StripeWebhook(c *gin.Context) {
	h.handleWebhook(c, "stripe", c.GetHeader("Stripe-Signature"))
}

// ── POST /api/payments/razorpay/webhook ──────────────────────────────────────
func (h *PaymentGatewayHandler) RazorpayWebhook(c *gin.Context) {
	h.handleWebhook(c, "razorpay", c.GetHeader("X-Razorpay-Signature"))
}

// ─── shared helpers ───────────────────────────────────────────────────────────

// createCheckout is the shared checkout creation flow for any provider.
func (h *PaymentGatewayHandler) createCheckout(c *gin.Context, gw pg.Gateway) {
	var body struct {
		Plan         string `json:"plan"          binding:"required"`
		PurchaseType string `json:"purchase_type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan is required"})
		return
	}
	pt := body.PurchaseType
	if pt == "" {
		pt = "membership"
	}
	uid, _ := getUID(c)
	appURL := os.Getenv("APP_URL")

	in, err := buildCheckoutInput(body.Plan, pt, uid, appURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := gw.CreateCheckout(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("gateway error: %v", err)})
		return
	}

	resp := gin.H{"provider": gw.Name()}
	if result.URL != "" {
		resp["url"] = result.URL
	}
	if result.ProviderRefID != "" {
		resp["order_id"] = result.ProviderRefID
	}
	c.JSON(http.StatusOK, resp)
}

// handleWebhook is the shared webhook handler for any provider.
func (h *PaymentGatewayHandler) handleWebhook(c *gin.Context, providerName, signature string) {
	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 65536))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}
	gw, err := h.pgFactory.ResolveCheckout(providerName)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	event, err := gw.VerifyWebhook(payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if event.UserID > 0 {
		h.fulfil(c, event.UserID, event.Plan, event.PurchaseType)
		return
	}
	c.JSON(http.StatusOK, gin.H{"received": true, "event": event.EventType})
}

// fulfil credits coins or activates a membership after payment is confirmed.
func (h *PaymentGatewayHandler) fulfil(c *gin.Context, uid uint, plan, purchaseType string) {
	switch purchaseType {
	case "coins":
		var coins int64
		if _, err := fmt.Sscanf(plan, "%d", &coins); err != nil || coins <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coin pack"})
			return
		}
		if h.coinSvc == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "coin service unavailable"})
			return
		}
		if err := h.coinSvc.CreditTopUp(c.Request.Context(), uid, coins, "payment_gateway", nil); err != nil {
			respondErr(c, err)
			return
		}
		respond(c, gin.H{"credited_coins": coins})

	case "membership":
		ep, ok := planEntityMap[plan]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan: " + plan})
			return
		}
		membership, err := h.membershipSvc.Activate(c.Request.Context(), uid, ep)
		if err != nil {
			respondErr(c, err)
			return
		}
		respond(c, membership)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown purchase_type: " + purchaseType})
	}
}

// buildCheckoutInput constructs a provider-agnostic CheckoutInput.
func buildCheckoutInput(plan, purchaseType string, uid uint, appURL string) (pg.CheckoutInput, error) {
	var amount int64
	var title string

	switch purchaseType {
	case "coins":
		var coins int64
		if _, err := fmt.Sscanf(plan, "%d", &coins); err != nil || coins <= 0 {
			return pg.CheckoutInput{}, fmt.Errorf("invalid coin pack %q", plan)
		}
		amount = (coins / 10) * 100 // 10 coins = $1
		title = fmt.Sprintf("%d Platform Coins", coins)

	case "membership":
		paise, ok := planAmountsPaise[plan]
		if !ok {
			return pg.CheckoutInput{}, fmt.Errorf("invalid plan %q", plan)
		}
		amount = paise
		title = fmt.Sprintf("%s Membership", plan)

	default:
		return pg.CheckoutInput{}, fmt.Errorf("unknown purchase_type %q", purchaseType)
	}

	return pg.CheckoutInput{
		AmountMinor: amount,
		Currency:    "inr",
		Title:       title,
		SuccessURL:  appURL + "/dashboard?payment=success",
		CancelURL:   appURL + "/checkout?payment=cancelled",
		Receipt:     fmt.Sprintf("%s_%s_uid%d", purchaseType, plan, uid),
		Metadata: map[string]string{
			"user_id":       fmt.Sprintf("%d", uid),
			"plan":          plan,
			"purchase_type": purchaseType,
		},
		Notes: map[string]any{
			"user_id":       uid,
			"plan":          plan,
			"purchase_type": purchaseType,
		},
	}, nil
}
