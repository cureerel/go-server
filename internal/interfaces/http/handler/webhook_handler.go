package handler

import (
	"io"
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/pkg/logger"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
	log            logger.Logger
}

func NewWebhookHandler(svc *service.WebhookService, log logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		webhookService: svc,
		log:            log,
	}
}

// HandleStripe handles Stripe webhooks
func (h *WebhookHandler) HandleStripe(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
		return
	}

	if err := h.webhookService.ProcessStripeWebhook(c.Request.Context(), payload, signature); err != nil {
		// FIX: Use proper Field struct
		h.log.Error("Stripe webhook processing failed",
			logger.Field{Key: "error", Value: err.Error()},
			logger.Field{Key: "path", Value: c.Request.URL.Path},
		)
		// Return 200 to prevent Stripe retries for invalid signatures
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// HandleRazorpay handles Razorpay webhooks
func (h *WebhookHandler) HandleRazorpay(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
		return
	}

	if err := h.webhookService.ProcessRazorpayWebhook(c.Request.Context(), payload, signature); err != nil {
		h.log.Error("Razorpay webhook processing failed",
			logger.Field{Key: "error", Value: err.Error()},
			logger.Field{Key: "path", Value: c.Request.URL.Path},
		)
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}
