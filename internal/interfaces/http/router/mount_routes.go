package router

import (
	"github.com/gin-gonic/gin"
)

// MountRoutes registers
func MountRoutes(r *gin.Engine, d *Deps) {
	api := r.Group("/api")

	// Register modular routes
	registerAuthRoutes(api, d)
	registerBlogRoutes(api, d)
	registerServiceRoutes(api, d)
	registerOrderRoutes(api, d)
	registerPaymentRoutes(api, d)
	registerCouponRoutes(api, d)
	registerTicketRoutes(api, d)
	registerUserRoutes(api, d)
	registerDashboardRoutes(api, d)
	registerAdminRoutes(api, d)
	registerMembershipRoutes(api, d)
	registerCoinRoutes(api, d)
	registerUploadRoutes(api, d)
	registerProductRoutes(api, d)
	registerWebhookRoutes(api, d)
}
