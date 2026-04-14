package router

import (
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerBlogRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public blog reads
	pub := rg.Group("/blog")
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
		{
			blogs.GET("/mine", d.BlogHandler.GetMine)
			blogs.GET("/:id/stats", d.BlogHandler.GetStats)
			blogs.POST("", d.BlogHandler.Create)
			blogs.PUT("/:id", d.BlogHandler.Update)
			blogs.PATCH("/:id", d.BlogHandler.Patch)
			blogs.DELETE("/:id", d.BlogHandler.Delete)

		}

	}
}
