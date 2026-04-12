package router

import (
	"github.com/gin-gonic/gin"
)

// MountRoutes registers all API routes (split by area; single mount keeps main.go stable).
func MountRoutes(r *gin.Engine, d *Deps) {
	api := r.Group("/api")

	// Register modular routes
	registerAuthRoutes(api, d)
	registerBlogRoutes(api, d)
	registerServiceRoutes(api, d)
	registerOrderRoutes(api, d)
	registerPaymentRoutes(api, d)
	registerCouponRoutes(api, d)
	registerPayoutRoutes(api, d)
	registerTicketRoutes(api, d)
	registerUserRoutes(api, d)
	registerDashboardRoutes(api, d)
	registerSuperAdminRoutes(api, d)
	registerMembershipRoutes(api, d)
	registerCoinRoutes(api, d)
	registerUploadRoutes(api, d)
	registerProductRoutes(api, d)
	registerWebhookRoutes(api, d)
}

