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
	log logger.Logger,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://cureerel.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.ErrorHandler()) // FIX: Remove log parameter

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().UTC()})
	})

	// Public API Routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			// TODO: auth.POST("/logout", authHandler.Logout) // Add later
		}

		// Public blog access (read-only)
		blogsPublic := api.Group("/blogs")
		{
			blogsPublic.GET("", blogHandler.GetAll)
			blogsPublic.GET("/:id", blogHandler.GetByID)
			blogsPublic.GET("/slug/:slug", blogHandler.GetBySlug)
		}
	}

	// Protected API Routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			// TODO: users.GET("/me", userHandler.GetCurrentUser) // Add later
			users.GET("/:id", userHandler.GetUserByID)

			// Admin only
			admin := users.Group("")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.POST("", userHandler.CreateUser)
				admin.PUT("/:id", userHandler.UpdateUser)
				admin.DELETE("/:id", userHandler.DeleteUser)
			}
		}

		// Blog routes (protected write operations)
		blogsProtected := protected.Group("/blogs")
		{
			blogsProtected.POST("", blogHandler.Create)
			blogsProtected.PUT("/:id", blogHandler.Update)
			blogsProtected.PATCH("/:id", blogHandler.Patch)
			blogsProtected.DELETE("/:id", blogHandler.Delete)
			// TODO: blogsProtected.PATCH("/:id/publish", blogHandler.Publish) // Add later
		}
	}

	// Webhooks (public but signature verified)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/stripe", webhookHandler.HandleStripe)
		webhooks.POST("/razorpay", webhookHandler.HandleRazorpay)
	}

	return r
}