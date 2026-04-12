package router

import (
	"github.com/gin-gonic/gin"
)

func registerWebhookRoutes(rg *gin.RouterGroup, d *Deps) {
	webhooks := rg.Group("/webhooks")
	{
		webhooks.POST("/stripe", d.WebhookHandler.HandleStripe)
		webhooks.POST("/razorpay", d.WebhookHandler.HandleRazorpay)
	}
}
