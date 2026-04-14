// internal/interfaces/http/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/interfaces/http/handler"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/cureerel/cserver/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	blogHandler *handler.BlogHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	serviceHandler *handler.ServiceHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	couponHandler *handler.CouponHandler,

	ticketHandler *handler.TicketHandler,
	dashboardHandler *handler.DashboardHandler,
	superadminHandler *handler.AdminHandler,
	uploadHandler *handler.UploadHandler,
	membershipHandler *handler.MembershipHandler,
	pgHandler *handler.PaymentGatewayHandler,
	coinHandler *handler.CoinHandler,
	productHandler *handler.ProductHandler,
	webhookHandler *handler.WebhookHandler,
	log logger.Logger,
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

	d := &Deps{
		AuthService:       authService,
		UserHandler:       userHandler,
		BlogHandler:       blogHandler,
		AuthHandler:       authHandler,
		ServiceHandler:    serviceHandler,
		OrderHandler:      orderHandler,
		PaymentHandler:    paymentHandler,
		CouponHandler:     couponHandler,
		TicketHandler:     ticketHandler,
		DashboardHandler:  dashboardHandler,
		AdminHandler:      adminHandler,
		UploadHandler:     uploadHandler,
		MembershipHandler: membershipHandler,
		PGHandler:         pgHandler,
		CoinHandler:       coinHandler,
		ProductHandler:    productHandler,
		WebhookHandler:    webhookHandler,
		Log:               log,
		AllowedOrigins:    allowedOrigins,
	}
	MountRoutes(r, d)
	return r
}
