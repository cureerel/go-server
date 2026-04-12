package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerPaymentRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public webhooks
	rg.POST("/payments/stripe/webhook", d.PGHandler.StripeWebhook)
	rg.POST("/payments/razorpay/webhook", d.PGHandler.RazorpayWebhook)

	payments := rg.Group("/payments")
	payments.Use(middleware.AuthMiddleware(d.AuthService))
	{
		payments.GET("/:id", d.PaymentHandler.GetByID)
		
		// Payment Gateway operations
		payments.POST("/razorpay/create-order", d.PGHandler.RazorpayCreateOrder)
		payments.POST("/razorpay/verify", d.PGHandler.RazorpayVerify)
		payments.POST("/stripe/create-session", d.PGHandler.StripeCreateSession)

		// Admin operations
		admin := payments.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.PaymentHandler.GetAll)
			admin.POST("/:id/complete", d.PaymentHandler.MarkCompleted)
			admin.POST("/:id/fail", d.PaymentHandler.MarkFailed)
			admin.POST("/:id/refund", d.PaymentHandler.Refund)
		}
	}
}
