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

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
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
	rzpClient     *rzpsdk.Client
}

func NewPaymentGatewayHandler(membershipSvc *service.MembershipService) *PaymentGatewayHandler {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	rzpClient := rzpsdk.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	return &PaymentGatewayHandler{
		membershipSvc: membershipSvc,
		rzpClient:     rzpClient,
	}
}

// ── POST /api/payments/razorpay/create-order ─────────────────────
// Body:    { "plan": "basic" | "pro" }
// Returns: { "order_id": "order_xxx", "key_id": "rzp_live_xxx", "amount": 49900 }
func (h *PaymentGatewayHandler) RazorpayCreateOrder(c *gin.Context) {
	var body struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan is required"})
		return
	}

	amount, ok := planAmountsPaise[body.Plan]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
		return
	}

	uid, _ := getUID(c)

	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  fmt.Sprintf("membership_%s_user%d", body.Plan, uid),
		"notes": map[string]interface{}{
			"user_id": uid,
			"plan":    body.Plan,
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
// Body: { "plan", "razorpay_order_id", "razorpay_payment_id", "razorpay_signature" }
// Verifies HMAC signature then activates membership.
func (h *PaymentGatewayHandler) RazorpayVerify(c *gin.Context) {
	var body struct {
		Plan               string `json:"plan"                binding:"required"`
		RazorpayOrderID    string `json:"razorpay_order_id"   binding:"required"`
		RazorpayPaymentID  string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature  string `json:"razorpay_signature"  binding:"required"`
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

	plan, ok := planEntityMap[body.Plan]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
		return
	}

	uid, _ := getUID(c)
	membership, err := h.membershipSvc.Activate(c.Request.Context(), uid, plan)
	if err != nil {
		respondErr(c, err)
		return
	}

	respond(c, membership)
}

// ── POST /api/payments/stripe/create-session ─────────────────────
// Body:    { "plan": "basic" | "pro" }
// Returns: { "url": "https://checkout.stripe.com/pay/cs_xxx" }
func (h *PaymentGatewayHandler) StripeCreateSession(c *gin.Context) {
	var body struct {
		Plan string `json:"plan" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan is required"})
		return
	}

	amount, ok := planAmountsPaise[body.Plan]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
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
						Name: stripe.String(fmt.Sprintf("%s Membership", body.Plan)),
					},
				},
			},
		},
		SuccessURL: stripe.String(fmt.Sprintf("%s/payment/success?plan=%s", appURL, body.Plan)),
		CancelURL:  stripe.String(fmt.Sprintf("%s?payment=cancelled", appURL)),
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", uid))
	params.AddMetadata("plan", body.Plan)

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
			var uid uint
			fmt.Sscanf(sess.Metadata["user_id"], "%d", &uid)

			if plan, ok := planEntityMap[planStr]; ok && uid > 0 {
				// best-effort — success page may have already activated
				h.membershipSvc.Activate(c.Request.Context(), uid, plan) //nolint:errcheck
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}