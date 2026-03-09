// internal/interfaces/http/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/middleware"
	"github.com/cureerel/gotemplate/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler    *handler.UserHandler,
	blogHandler    *handler.BlogHandler,
	authHandler    *handler.AuthHandler,
	authService    *service.AuthService,
	serviceHandler *handler.ServiceHandler,
	orderHandler   *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	uploadHandler  *handler.UploadHandler,
	log            logger.Logger,
	allowedOrigins []string,
) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())

	origins := allowedOrigins
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "ts": time.Now().UTC()})
	})

	api := r.Group("/api")

	// ── Public ────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register/init",         authHandler.RegisterInit)
		auth.POST("/register/verify",       authHandler.RegisterVerify)
		auth.POST("/password/reset/init",   authHandler.PasswordResetInit)
		auth.POST("/password/reset/verify", authHandler.PasswordResetVerify)
		auth.POST("/signup",                authHandler.Signup)
		auth.POST("/login",                 authHandler.Login)
		auth.POST("/refresh",               authHandler.Refresh)
	}

	api.GET("/blog",            blogHandler.GetAll)
	api.GET("/blog/slug/:slug", blogHandler.GetBySlug)
	api.GET("/blog/:id",        blogHandler.GetByID)

	api.GET("/services",     serviceHandler.GetAll)
	api.GET("/services/:id", serviceHandler.GetByID)

	// ── Protected ─────────────────────────────────────────────
	p := api.Group("")
	p.Use(middleware.AuthMiddleware(authService))
	{
		p.POST("/auth/logout", authHandler.Logout)
		p.GET("/users/me",     userHandler.GetMe)

		p.POST("/upload/image",   uploadHandler.UploadImage)
		p.DELETE("/upload/image", uploadHandler.DeleteImage)

		p.GET("/services/mine", serviceHandler.GetMine)

		// Orders
		orders := p.Group("/orders")
		{
			orders.POST("",    orderHandler.Create)
			orders.GET("/me",  orderHandler.GetMyOrders)
			orders.GET("/:id", orderHandler.GetByID)
		}

		// Payment — own
		p.GET("/payments/:id", paymentHandler.GetByID)
	}

	// ── writer+ ───────────────────────────────────────────────
	blogs := p.Group("/blogs")
	blogs.Use(middleware.RoleMiddleware(entity.RoleWriter))
	{
		blogs.POST("",       blogHandler.Create)
		blogs.PUT("/:id",    blogHandler.Update)
		blogs.PATCH("/:id",  blogHandler.Patch)
		blogs.DELETE("/:id", blogHandler.Delete)
	}

	// ── partner+ ──────────────────────────────────────────────
	myServices := p.Group("/services")
	myServices.Use(middleware.RoleMiddleware(entity.RolePartner))
	{
		myServices.POST("",           serviceHandler.Create)
		myServices.PUT("/:id",        serviceHandler.Update)
		myServices.DELETE("/:id",     serviceHandler.Delete)
		myServices.POST("/:id/live",  serviceHandler.SetLive)
		myServices.POST("/:id/pause", serviceHandler.Pause)
	}

	// ── admin+ ────────────────────────────────────────────────
	adminUsers := p.Group("/users")
	adminUsers.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		adminUsers.GET("",        userHandler.GetAllUsers)
		adminUsers.GET("/:id",    userHandler.GetUserByID)
		adminUsers.POST("",       userHandler.CreateUser)
		adminUsers.PUT("/:id",    userHandler.UpdateUser)
		adminUsers.DELETE("/:id", userHandler.DeleteUser)
	}

	adminServices := p.Group("/services")
	adminServices.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		adminServices.POST("/:id/approve", serviceHandler.Approve)
		adminServices.POST("/:id/reject",  serviceHandler.Reject)
		adminServices.POST("/:id/pause",   serviceHandler.Pause)
	}

	adminOrders := p.Group("/orders")
	adminOrders.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		adminOrders.GET("",              orderHandler.GetAll)
		adminOrders.PATCH("/:id/status", orderHandler.UpdateStatus)
	}

	adminPayments := p.Group("/payments")
	adminPayments.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		adminPayments.GET("",               paymentHandler.GetAll)
		adminPayments.POST("/:id/complete", paymentHandler.MarkCompleted)
		adminPayments.POST("/:id/fail",     paymentHandler.MarkFailed)
		adminPayments.POST("/:id/refund",   paymentHandler.Refund)
	}

	return r
}