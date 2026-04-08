package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

// MountRoutes registers all API routes (split by area; single mount keeps main.go stable).
func MountRoutes(r *gin.Engine, d *Deps) {
	api := r.Group("/api")

	// ── Public auth ─────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register/init", d.AuthHandler.RegisterInit)
		auth.POST("/register/verify", d.AuthHandler.RegisterVerify)
		auth.POST("/password/reset/init", d.AuthHandler.PasswordResetInit)
		auth.POST("/password/reset/verify", d.AuthHandler.PasswordResetVerify)
		auth.POST("/signup", d.AuthHandler.Signup)
		auth.POST("/login", d.AuthHandler.Login)
		auth.POST("/refresh", d.AuthHandler.Refresh)
	}

	api.GET("/coupons/validate", d.CouponHandler.Validate)
	api.POST("/payments/stripe/webhook", d.PGHandler.StripeWebhook)

	// Public catalog + blog reads (optional JWT for access-gated content)
	pub := api.Group("")
	pub.Use(middleware.OptionalAuth(d.AuthService))
	{
		pub.GET("/blog", d.BlogHandler.GetAll)
		pub.GET("/blog/slug/:slug", d.BlogHandler.GetBySlug)
		pub.GET("/blog/:id", d.BlogHandler.GetByID)
		pub.GET("/services", d.ServiceHandler.GetAll)
		pub.GET("/services/:id", d.ServiceHandler.GetByID)
	}

	// ── Authenticated ───────────────────────────────────────────
	p := api.Group("")
	p.Use(middleware.AuthMiddleware(d.AuthService))
	{
		p.POST("/auth/logout", d.AuthHandler.Logout)
		p.GET("/users/me", d.UserHandler.GetMe)

		p.POST("/upload/image", d.UploadHandler.UploadImage)
		p.DELETE("/upload/image", d.UploadHandler.DeleteImage)

		p.GET("/services/mine", d.ServiceHandler.GetMine)

		p.GET("/coins/balance", d.CoinHandler.GetBalance)

		memberships := p.Group("/memberships")
		{
			memberships.GET("/me", d.MembershipHandler.GetMine)
			memberships.POST("/activate", d.MembershipHandler.Activate)
			memberships.POST("/upgrade", d.MembershipHandler.Upgrade)
			memberships.DELETE("/cancel", d.MembershipHandler.Cancel)
		}

		p.POST("/payments/razorpay/create-order", d.PGHandler.RazorpayCreateOrder)
		p.POST("/payments/razorpay/verify", d.PGHandler.RazorpayVerify)
		p.POST("/payments/stripe/create-session", d.PGHandler.StripeCreateSession)

		orders := p.Group("/orders")
		{
			orders.POST("", d.OrderHandler.Create)
			orders.GET("/me", d.OrderHandler.GetMyOrders)
			orders.GET("/:id", d.OrderHandler.GetByID)
			orders.GET("", middleware.RoleMiddleware(entity.RoleAdmin), d.OrderHandler.GetAll)
			orders.PATCH("/:id/status", middleware.RoleMiddleware(entity.RoleAdmin), d.OrderHandler.UpdateStatus)
		}

		payments := p.Group("/payments")
		{
			payments.GET("/:id", d.PaymentHandler.GetByID)
			payments.GET("", middleware.RoleMiddleware(entity.RoleAdmin), d.PaymentHandler.GetAll)
			payments.POST("/:id/complete", middleware.RoleMiddleware(entity.RoleAdmin), d.PaymentHandler.MarkCompleted)
			payments.POST("/:id/fail", middleware.RoleMiddleware(entity.RoleAdmin), d.PaymentHandler.MarkFailed)
			payments.POST("/:id/refund", middleware.RoleMiddleware(entity.RoleAdmin), d.PaymentHandler.Refund)
		}

		payouts := p.Group("/payouts")
		{
			payouts.GET("/me", d.PayoutHandler.GetMine)
			payouts.GET("", middleware.RoleMiddleware(entity.RoleAdmin), d.PayoutHandler.GetAll)
			payouts.POST("/:id/pay", middleware.RoleMiddleware(entity.RoleAdmin), d.PayoutHandler.MarkPaid)
		}

		coupons := p.Group("/coupons")
		{
			coupons.POST("", middleware.RoleMiddleware(entity.RolePartner), d.CouponHandler.Create)
			coupons.GET("/:id", middleware.RoleMiddleware(entity.RolePartner), d.CouponHandler.GetByID)
			coupons.GET("", middleware.RoleMiddleware(entity.RoleAdmin), d.CouponHandler.GetAll)
			coupons.POST("/:id/approve", middleware.RoleMiddleware(entity.RoleAdmin), d.CouponHandler.Approve)
			coupons.POST("/:id/reject", middleware.RoleMiddleware(entity.RoleAdmin), d.CouponHandler.Reject)
		}

		services := p.Group("/services")
		{
			services.POST("", middleware.RoleMiddleware(entity.RolePartner), d.ServiceHandler.Create)
			services.PUT("/:id", middleware.RoleMiddleware(entity.RolePartner), d.ServiceHandler.Update)
			services.DELETE("/:id", middleware.RoleMiddleware(entity.RolePartner), d.ServiceHandler.Delete)
			services.POST("/:id/live", middleware.RoleMiddleware(entity.RolePartner), d.ServiceHandler.SetLive)
			services.POST("/:id/pause", middleware.RoleMiddleware(entity.RolePartner), d.ServiceHandler.Pause)
			services.POST("/:id/approve", middleware.RoleMiddleware(entity.RoleAdmin), d.ServiceHandler.Approve)
			services.POST("/:id/reject", middleware.RoleMiddleware(entity.RoleAdmin), d.ServiceHandler.Reject)
		}

		tickets := p.Group("/tickets")
		{
			tickets.POST("", d.TicketHandler.Create)
			tickets.GET("/me", d.TicketHandler.GetMine)
			tickets.GET("/:id", d.TicketHandler.GetByID)
			tickets.POST("/:id/close", d.TicketHandler.Close)
			tickets.POST("/:id/messages", d.TicketHandler.SendMessage)
			tickets.GET("/:id/messages", d.TicketHandler.GetMessages)
			tickets.GET("", middleware.RoleMiddleware(entity.RoleWorker), d.TicketHandler.GetAll)
			tickets.POST("/:id/assign", middleware.RoleMiddleware(entity.RoleWorker), d.TicketHandler.Assign)
			tickets.POST("/:id/resolve", middleware.RoleMiddleware(entity.RoleWorker), d.TicketHandler.Resolve)
		}

		dash := p.Group("/dashboard")
		{
			dash.GET("", d.DashboardHandler.Get)
			dash.GET("/user", d.DashboardHandler.UserView)
			dash.GET("/writer", middleware.WriterContent(), d.DashboardHandler.WriterView)
			dash.GET("/partner", middleware.RoleMiddleware(entity.RolePartner), d.DashboardHandler.PartnerView)
			dash.GET("/worker", middleware.RoleMiddleware(entity.RoleWorker), d.DashboardHandler.WorkerView)
			dash.GET("/admin", middleware.RoleMiddleware(entity.RoleAdmin), d.DashboardHandler.AdminView)
		}

		users := p.Group("/users")
		{
			users.GET("", middleware.RoleMiddleware(entity.RoleAdmin), d.UserHandler.GetAllUsers)
			users.GET("/:id", middleware.RoleMiddleware(entity.RoleAdmin), d.UserHandler.GetUserByID)
			users.POST("", middleware.RoleMiddleware(entity.RoleAdmin), d.UserHandler.CreateUser)
			users.PUT("/:id", middleware.RoleMiddleware(entity.RoleAdmin), d.UserHandler.UpdateUser)
			users.DELETE("/:id", middleware.RoleMiddleware(entity.RoleAdmin), d.UserHandler.DeleteUser)
		}

		p.POST("/blogs/:id/unlock", d.BlogHandler.UnlockPaidBlog)

		blogs := p.Group("/blogs")
		blogs.Use(middleware.WriterContent())
		{
			blogs.GET("/mine", d.BlogHandler.GetMine)
			blogs.GET("/:id/stats", d.BlogHandler.GetStats)
			blogs.POST("", d.BlogHandler.Create)
			blogs.PUT("/:id", d.BlogHandler.Update)
			blogs.PATCH("/:id", d.BlogHandler.Patch)
			blogs.DELETE("/:id", d.BlogHandler.Delete)
			blogs.POST("/:id/submit-review", d.BlogHandler.SubmitForReview)
		}

		rev := p.Group("/reviewer")
		rev.Use(middleware.RequireAnyRole(entity.RoleReviewer, entity.RoleAdmin, entity.RoleSuperAdmin))
		{
			rev.GET("/blogs/pending", d.BlogHandler.ListReviewQueue)
			rev.POST("/blogs/:id/approve", d.BlogHandler.ReviewApprove)
			rev.POST("/blogs/:id/reject", d.BlogHandler.ReviewReject)
		}

		upgrades := p.Group("/upgrade-requests")
		{
			upgrades.POST("", d.SuperAdminHandler.RequestUpgrade)
			upgrades.GET("/me", d.SuperAdminHandler.GetMyUpgradeRequest)
		}

		sa := p.Group("/superadmin")
		sa.Use(middleware.RoleMiddleware(entity.RoleSuperAdmin))
		{
			sa.GET("/stats", d.SuperAdminHandler.Stats)
			sa.PATCH("/users/:id/role", d.SuperAdminHandler.SetRole)
			sa.GET("/upgrades", d.SuperAdminHandler.GetPendingUpgrades)
			sa.POST("/upgrades/:id/review", d.SuperAdminHandler.ReviewUpgrade)
		}
	}
}
