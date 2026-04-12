package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerBlogRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public blog reads
	pub := rg.Group("/blog")
	pub.Use(middleware.OptionalAuth(d.AuthService))
	{
		pub.GET("", d.BlogHandler.GetAll)
		pub.GET("/slug/:slug", d.BlogHandler.GetBySlug)
		pub.GET("/:id", d.BlogHandler.GetByID)
	}

	// Authenticated blog actions
	auth := rg.Group("")
	auth.Use(middleware.AuthMiddleware(d.AuthService))
	{
		auth.POST("/blogs/:id/unlock", d.BlogHandler.UnlockPaidBlog)

		blogs := auth.Group("/blogs")
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

		rev := auth.Group("/reviewer")
		rev.Use(middleware.RequireAnyRole(entity.RoleReviewer, entity.RoleAdmin, entity.RoleSuperAdmin))
		{
			rev.GET("/blogs/pending", d.BlogHandler.ListReviewQueue)
			rev.POST("/blogs/:id/approve", d.BlogHandler.ReviewApprove)
			rev.POST("/blogs/:id/reject", d.BlogHandler.ReviewReject)
		}
	}
}
