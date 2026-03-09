package router

import (

	"time"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/middleware"
	"github.com/cureerel/gotemplate/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	blogHandler *handler.BlogHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	webhookHandler *handler.WebhookHandler,
	productHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	membershipHandler *handler.MembershipHandler,
	log logger.Logger,
	allowedOrigins []string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Use env origins or fall back to defaults
	origins := allowedOrigins
	if len(origins) == 0 {
		origins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8080",
		}
	}

	// Debug log for CORS origins (remove in production if desired)
	log.Info("CORS configured", logger.Field{Key: "origins", Value: origins})

	r.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Request-ID",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().UTC()})
	})

	api := r.Group("/api")

	// ── Public routes ─────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Public blog read
	api.GET("/blog", blogHandler.GetAll)
    api.GET("/blog/slug/:slug", blogHandler.GetBySlug) 
    api.GET("/blog/:id", blogHandler.GetByID)  

	// Public product read
	api.GET("/products", productHandler.GetAll)
	api.GET("/products/:id", productHandler.GetByID)

	// ── Protected routes ──────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Auth
		protected.POST("/auth/logout", authHandler.Logout)

		// Users — self
		protected.GET("/users/me", userHandler.GetMe)

		// Users — admin only
		adminUsers := protected.Group("/users")
		adminUsers.Use(middleware.RoleMiddleware("admin"))
		{
			adminUsers.GET("", userHandler.GetAllUsers)
			adminUsers.GET("/:id", userHandler.GetUserByID)
			adminUsers.POST("", userHandler.CreateUser)
			adminUsers.PUT("/:id", userHandler.UpdateUser)
			adminUsers.DELETE("/:id", userHandler.DeleteUser)
		}

		// Blogs — authenticated write
		blogs := protected.Group("/blogs")
		{
			blogs.POST("", blogHandler.Create)
			blogs.PUT("/:id", blogHandler.Update)
			blogs.PATCH("/:id", blogHandler.Patch)
			blogs.DELETE("/:id", blogHandler.Delete)
		}

		// Products — admin write
		adminProducts := protected.Group("/products")
		adminProducts.Use(middleware.RoleMiddleware("admin"))
		{
			adminProducts.POST("", productHandler.Create)
			adminProducts.PUT("/:id", productHandler.Update)
			adminProducts.DELETE("/:id", productHandler.Delete)
		}

		// Orders — authenticated users
		orders := protected.Group("/orders")
		{
			orders.POST("", orderHandler.Create)
			orders.GET("/me", orderHandler.GetMyOrders)
			orders.GET("/:id", orderHandler.GetByID)
		}

		// Orders — admin status update
		adminOrders := protected.Group("/orders")
		adminOrders.Use(middleware.RoleMiddleware("admin"))
		{
			adminOrders.PATCH("/:id/status", orderHandler.UpdateStatus)
		}

		// Payments — owner or admin
		payments := protected.Group("/payments")
		{
			payments.GET("/:id", paymentHandler.GetByID)
		}

		// Payments — admin only
		adminPayments := protected.Group("/payments")
		adminPayments.Use(middleware.RoleMiddleware("admin"))
		{
			adminPayments.POST("/:id/refund", paymentHandler.Refund)
			adminPayments.POST("/:id/complete", paymentHandler.MarkCompleted)
			adminPayments.POST("/:id/fail", paymentHandler.MarkFailed)
		}

		// Memberships — authenticated users
		memberships := protected.Group("/memberships")
		{
			memberships.GET("/me", membershipHandler.GetMine)
			memberships.POST("/activate", membershipHandler.Activate)
			memberships.POST("/upgrade", membershipHandler.Upgrade)
			memberships.DELETE("/cancel", membershipHandler.Cancel)
		}
	}

	// ── Webhooks — public but signature verified ───────────────
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/stripe", webhookHandler.HandleStripe)
		webhooks.POST("/razorpay", webhookHandler.HandleRazorpay)
	}

	return r
}