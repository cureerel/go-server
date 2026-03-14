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
	userHandler          *handler.UserHandler,
	blogHandler          *handler.BlogHandler,
	authHandler          *handler.AuthHandler,
	authService          *service.AuthService,
	serviceHandler       *handler.ServiceHandler,
	orderHandler         *handler.OrderHandler,
	paymentHandler       *handler.PaymentHandler,
	couponHandler        *handler.CouponHandler,
	payoutHandler        *handler.PayoutHandler,
	ticketHandler        *handler.TicketHandler,
	dashboardHandler     *handler.DashboardHandler,
	superadminHandler    *handler.SuperAdminHandler,
	uploadHandler        *handler.UploadHandler,
	membershipHandler    *handler.MembershipHandler,      
	pgHandler            *handler.PaymentGatewayHandler,  
	log                  logger.Logger,
	allowedOrigins       []string,
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

	// ── Public ────────────────────────────────────────────────────
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

	api.GET("/blog",             blogHandler.GetAll)
	api.GET("/blog/slug/:slug",  blogHandler.GetBySlug)
	api.GET("/blog/:id",         blogHandler.GetByID)
	api.GET("/services",         serviceHandler.GetAll)
	api.GET("/services/:id",     serviceHandler.GetByID)
	api.GET("/coupons/validate", couponHandler.Validate)

	// Stripe webhook — NO auth middleware, Stripe sends Stripe-Signature header
	api.POST("/payments/stripe/webhook", pgHandler.StripeWebhook)

	// ── Protected (any authenticated user) ────────────────────────
	p := api.Group("")
	p.Use(middleware.AuthMiddleware(authService))
	{
		p.POST("/auth/logout", authHandler.Logout)
		p.GET("/users/me",     userHandler.GetMe)

		p.POST("/upload/image",   uploadHandler.UploadImage)
		p.DELETE("/upload/image", uploadHandler.DeleteImage)

		p.GET("/services/mine", serviceHandler.GetMine)

		// ── Memberships ───────────────────────────────────────────
		memberships := p.Group("/memberships")
		{
			memberships.GET("/me",        membershipHandler.GetMine)
			memberships.POST("/activate", membershipHandler.Activate)
			memberships.POST("/upgrade",  membershipHandler.Upgrade)
			memberships.DELETE("/cancel", membershipHandler.Cancel)
		}

		// ── Payment gateway (Razorpay + Stripe session creation) ──
		p.POST("/payments/razorpay/create-order", pgHandler.RazorpayCreateOrder)
		p.POST("/payments/razorpay/verify",       pgHandler.RazorpayVerify)
		p.POST("/payments/stripe/create-session", pgHandler.StripeCreateSession)

		// ── Orders (merged group — no duplicate path conflict) ────
		orders := p.Group("/orders")
		{
			orders.POST("",    orderHandler.Create)
			orders.GET("/me",  orderHandler.GetMyOrders)
			orders.GET("/:id", orderHandler.GetByID)

			orders.GET("",              middleware.RoleMiddleware(entity.RoleAdmin), orderHandler.GetAll)
			orders.PATCH("/:id/status", middleware.RoleMiddleware(entity.RoleAdmin), orderHandler.UpdateStatus)
		}

		// ── Payments (merged group) ───────────────────────────────
		payments := p.Group("/payments")
		{
			payments.GET("/:id", paymentHandler.GetByID)

			payments.GET("",               middleware.RoleMiddleware(entity.RoleAdmin), paymentHandler.GetAll)
			payments.POST("/:id/complete", middleware.RoleMiddleware(entity.RoleAdmin), paymentHandler.MarkCompleted)
			payments.POST("/:id/fail",     middleware.RoleMiddleware(entity.RoleAdmin), paymentHandler.MarkFailed)
			payments.POST("/:id/refund",   middleware.RoleMiddleware(entity.RoleAdmin), paymentHandler.Refund)
		}

		// ── Payouts (merged group) ────────────────────────────────
		payouts := p.Group("/payouts")
		{
			payouts.GET("/me", payoutHandler.GetMine)

			payouts.GET("",          middleware.RoleMiddleware(entity.RoleAdmin), payoutHandler.GetAll)
			payouts.POST("/:id/pay", middleware.RoleMiddleware(entity.RoleAdmin), payoutHandler.MarkPaid)
		}

		// ── Coupons (merged group) ────────────────────────────────
		coupons := p.Group("/coupons")
		{
			coupons.POST("",    middleware.RoleMiddleware(entity.RolePartner), couponHandler.Create)
			coupons.GET("/:id", middleware.RoleMiddleware(entity.RolePartner), couponHandler.GetByID)

			coupons.GET("",              middleware.RoleMiddleware(entity.RoleAdmin), couponHandler.GetAll)
			coupons.POST("/:id/approve", middleware.RoleMiddleware(entity.RoleAdmin), couponHandler.Approve)
			coupons.POST("/:id/reject",  middleware.RoleMiddleware(entity.RoleAdmin), couponHandler.Reject)
		}

		// ── Services (merged group) ───────────────────────────────
		services := p.Group("/services")
		{
			services.POST("",           middleware.RoleMiddleware(entity.RolePartner), serviceHandler.Create)
			services.PUT("/:id",        middleware.RoleMiddleware(entity.RolePartner), serviceHandler.Update)
			services.DELETE("/:id",     middleware.RoleMiddleware(entity.RolePartner), serviceHandler.Delete)
			services.POST("/:id/live",  middleware.RoleMiddleware(entity.RolePartner), serviceHandler.SetLive)
			services.POST("/:id/pause", middleware.RoleMiddleware(entity.RolePartner), serviceHandler.Pause)

			services.POST("/:id/approve", middleware.RoleMiddleware(entity.RoleAdmin), serviceHandler.Approve)
			services.POST("/:id/reject",  middleware.RoleMiddleware(entity.RoleAdmin), serviceHandler.Reject)
		}

		// ── Tickets (merged group) ────────────────────────────────
		tickets := p.Group("/tickets")
		{
			tickets.POST("",              ticketHandler.Create)
			tickets.GET("/me",            ticketHandler.GetMine)
			tickets.GET("/:id",           ticketHandler.GetByID)
			tickets.POST("/:id/close",    ticketHandler.Close)
			tickets.POST("/:id/messages", ticketHandler.SendMessage)
			tickets.GET("/:id/messages",  ticketHandler.GetMessages)

			tickets.GET("",              middleware.RoleMiddleware(entity.RoleWorker), ticketHandler.GetAll)
			tickets.POST("/:id/assign",  middleware.RoleMiddleware(entity.RoleWorker), ticketHandler.Assign)
			tickets.POST("/:id/resolve", middleware.RoleMiddleware(entity.RoleWorker), ticketHandler.Resolve)
		}

		// ── Dashboard (merged group) ──────────────────────────────
		dash := p.Group("/dashboard")
		{
			dash.GET("",        dashboardHandler.Get)
			dash.GET("/user",   dashboardHandler.UserView)
			dash.GET("/writer", middleware.RoleMiddleware(entity.RoleWriter),  dashboardHandler.WriterView)
			dash.GET("/partner", middleware.RoleMiddleware(entity.RolePartner), dashboardHandler.PartnerView)
			dash.GET("/worker", middleware.RoleMiddleware(entity.RoleWorker),  dashboardHandler.WorkerView)
			dash.GET("/admin",  middleware.RoleMiddleware(entity.RoleAdmin),   dashboardHandler.AdminView)
		}

		// ── Users (admin+) ────────────────────────────────────────
		users := p.Group("/users")
		{
			users.GET("",        middleware.RoleMiddleware(entity.RoleAdmin), userHandler.GetAllUsers)
			users.GET("/:id",    middleware.RoleMiddleware(entity.RoleAdmin), userHandler.GetUserByID)
			users.POST("",       middleware.RoleMiddleware(entity.RoleAdmin), userHandler.CreateUser)
			users.PUT("/:id",    middleware.RoleMiddleware(entity.RoleAdmin), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.RoleMiddleware(entity.RoleAdmin), userHandler.DeleteUser)
		}

		// ── Blogs (writer+) ───────────────────────────────────────
		blogs := p.Group("/blogs")
		{
			blogs.POST("",       middleware.RoleMiddleware(entity.RoleWriter), blogHandler.Create)
			blogs.PUT("/:id",    middleware.RoleMiddleware(entity.RoleWriter), blogHandler.Update)
			blogs.PATCH("/:id",  middleware.RoleMiddleware(entity.RoleWriter), blogHandler.Patch)
			blogs.DELETE("/:id", middleware.RoleMiddleware(entity.RoleWriter), blogHandler.Delete)
		}

		// ── Upgrade requests ──────────────────────────────────────
		upgrades := p.Group("/upgrade-requests")
		{
			upgrades.POST("",   superadminHandler.RequestUpgrade)
			upgrades.GET("/me", superadminHandler.GetMyUpgradeRequest)
		}

		// ── Superadmin ────────────────────────────────────────────
		sa := p.Group("/superadmin")
		sa.Use(middleware.RoleMiddleware(entity.RoleSuperAdmin))
		{
			sa.GET("/stats",                superadminHandler.Stats)
			sa.PATCH("/users/:id/role",     superadminHandler.SetRole)
			sa.GET("/upgrades",             superadminHandler.GetPendingUpgrades)
			sa.POST("/upgrades/:id/review", superadminHandler.ReviewUpgrade)
		}
	}

	return r
}