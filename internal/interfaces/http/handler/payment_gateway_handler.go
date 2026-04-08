// internal/interfaces/http/handler/payment_gateway_handler.go
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/infrastructure/paymentgateway"
	"github.com/gin-gonic/gin"
	rzpsdk "github.com/razorpay/razorpay-go"
	stripe "github.com/stripe/stripe-go/v79"
	stripesession "github.com/stripe/stripe-go/v79/checkout/session"
	stripewebhook "github.com/stripe/stripe-go/v79/webhook"
)

// planAmountsPaise maps plan → amount in paise (INR smallest unit)
var planAmountsPaise = map[string]int64{
	"basic": 49900,  // ₹499
	"pro":   149900, // ₹1499
}

// coinPackPaise / coinPackCoins — fiat top-up packs (coins are platform credits; 1 coin ≈ ₹1 in these packs).
var coinPackPaise = map[string]int64{
	"100": 10000, // ₹100 → 100 coins
	"500": 50000,
}
var coinPackCoins = map[string]int64{
	"100": 100,
	"500": 500,
}

// planEntityMap maps string → entity.MembershipPlan
var planEntityMap = map[string]entity.MembershipPlan{
	"free":       entity.PlanFree,
	"basic":      entity.PlanBasic,
	"pro":        entity.PlanPro,
	"enterprise": entity.PlanEnterprise,
}

// ─────────────────────────────────────────────────────────────────
type PaymentGatewayHandler struct {
	membershipSvc *service.MembershipService
	coinSvc       *service.CoinService
	pgFactory     *paymentgateway.Factory
	rzpClient     *rzpsdk.Client
}

func NewPaymentGatewayHandler(membershipSvc *service.MembershipService, coinSvc *service.CoinService) *PaymentGatewayHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	rzpClient := rzpsdk.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	return &PaymentGatewayHandler{
		membershipSvc: membershipSvc,
		coinSvc:       coinSvc,
		pgFactory:     paymentgateway.FromEnv(),
		rzpClient:     rzpClient,
	}
}

// ── POST /api/payments/razorpay/create-order ─────────────────────
// Body:    { "plan": "...", "purchase_type": "membership"|"coins" }
// Returns: { "order_id": "order_xxx", "key_id": "rzp_live_xxx", "amount": 49900 }
func (h *PaymentGatewayHandler) RazorpayCreateOrder(c *gin.Context) {
	var body struct {
		Plan          string `json:"plan" binding:"required"`
		PurchaseType  string `json:"purchase_type"` // default membership
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan is required"})
		return
	}

	var amount int64
	var ok bool
	pt := body.PurchaseType
	if pt == "" {
		pt = "membership"
	}
	switch pt {
	case "coins":
		amount, ok = coinPackPaise[body.Plan]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coin pack"})
			return
		}
	case "membership":
		amount, ok = planAmountsPaise[body.Plan]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "purchase_type must be membership or coins"})
		return
	}

	uid, _ := getUID(c)

	receipt := fmt.Sprintf("%s_%s_user%d", pt, body.Plan, uid)
	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  receipt,
		"notes": map[string]interface{}{
			"user_id":        uid,
			"plan":           body.Plan,
			"purchase_type":  pt,
			"default_provider": h.pgFactory.DefaultTopupProvider(),
		},
	}

	order, err := h.rzpClient.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Razorpay order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id": order["id"],
		"key_id":   os.Getenv("RAZORPAY_KEY_ID"),
		"amount":   amount,
	})
}

// ── POST /api/payments/razorpay/verify ───────────────────────────
// Body: { "plan", "purchase_type", "razorpay_order_id", "razorpay_payment_id", "razorpay_signature" }
// Verifies HMAC signature then activates membership or credits coins.
func (h *PaymentGatewayHandler) RazorpayVerify(c *gin.Context) {
	var body struct {
		Plan              string `json:"plan" binding:"required"`
		PurchaseType      string `json:"purchase_type"`
		RazorpayOrderID   string `json:"razorpay_order_id"   binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify HMAC-SHA256: sign(order_id + "|" + payment_id)
	mac := hmac.New(sha256.New, []byte(os.Getenv("RAZORPAY_KEY_SECRET")))
	mac.Write([]byte(body.RazorpayOrderID + "|" + body.RazorpayPaymentID))
	expected := hex.EncodeToString(mac.Sum(nil))

	if expected != body.RazorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature mismatch"})
		return
	}

	pt := body.PurchaseType
	if pt == "" {
		pt = "membership"
	}
	uid, _ := getUID(c)

	switch pt {
	case "coins":
		coins, ok := coinPackCoins[body.Plan]
		if !ok || h.coinSvc == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coin pack"})
			return
		}
		if err := h.coinSvc.CreditTopUp(c.Request.Context(), uid, coins, "razorpay", nil); err != nil {
			respondErr(c, err)
			return
		}
		respond(c, gin.H{"credited_coins": coins})
		return
	case "membership":
		plan, ok := planEntityMap[body.Plan]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
		membership, err := h.membershipSvc.Activate(c.Request.Context(), uid, plan)
		if err != nil {
			respondErr(c, err)
			return
		}
		respond(c, membership)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase_type"})
	}
}

// ── POST /api/payments/stripe/create-session ─────────────────────
// Body:    { "plan": "basic" | "pro" }
// Returns: { "url": "https://checkout.stripe.com/pay/cs_xxx" }
func (h *PaymentGatewayHandler) StripeCreateSession(c *gin.Context) {
	var body struct {
		Plan         string `json:"plan" binding:"required"`
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
	var amount int64
	var ok bool
	var title string
	switch pt {
	case "coins":
		amount, ok = coinPackPaise[body.Plan]
		title = fmt.Sprintf("Coin pack %s", body.Plan)
	case "membership":
		amount, ok = planAmountsPaise[body.Plan]
		title = fmt.Sprintf("%s Membership", body.Plan)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase_type"})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan or pack"})
		return
	}

	uid, _ := getUID(c)
	appURL := os.Getenv("APP_URL") // e.g. http://localhost:3000

	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("inr"),
					UnitAmount: stripe.Int64(amount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(title),
					},
				},
			},
		},
		SuccessURL: stripe.String(fmt.Sprintf("%s/payment/success?plan=%s&purchase_type=%s", appURL, body.Plan, pt)),
		CancelURL:  stripe.String(fmt.Sprintf("%s?payment=cancelled", appURL)),
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", uid))
	params.AddMetadata("plan", body.Plan)
	params.AddMetadata("purchase_type", pt)

	sess, err := stripesession.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create Stripe session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": sess.URL})
}

// ── POST /api/payments/stripe/webhook ────────────────────────────
// No auth middleware — Stripe sends Stripe-Signature header.
// Register URL in Stripe Dashboard → Webhooks.
func (h *PaymentGatewayHandler) StripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(io.LimitReader(c.Request.Body, 65536))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	event, err := stripewebhook.ConstructEvent(
		payload,
		c.GetHeader("Stripe-Signature"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "webhook signature verification failed"})
		return
	}

	if event.Type == "checkout.session.completed" {
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err == nil &&
			(sess.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid ||
				sess.Status == stripe.CheckoutSessionStatusComplete) {

			planStr := sess.Metadata["plan"]
			pt := sess.Metadata["purchase_type"]
			if pt == "" {
				pt = "membership"
			}
			var uid uint
			fmt.Sscanf(sess.Metadata["user_id"], "%d", &uid)
			if uid > 0 {
				if pt == "coins" && h.coinSvc != nil {
					if coins, ok := coinPackCoins[planStr]; ok {
						_ = h.coinSvc.CreditTopUp(c.Request.Context(), uid, coins, "stripe", nil) //nolint:errcheck
					}
				} else if plan, ok := planEntityMap[planStr]; ok {
					h.membershipSvc.Activate(c.Request.Context(), uid, plan) //nolint:errcheck
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}