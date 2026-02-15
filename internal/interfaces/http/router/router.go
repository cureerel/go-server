package router

import (
	"time"

	"github.com/cureerel/gotemplate/internal/interfaces/http/handler"
	"github.com/cureerel/gotemplate/internal/interfaces/http/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, blogHandler *handler.BlogHandler) *gin.Engine {
	r := gin.Default()
	
	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://cureerel.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Global Middleware
	r.Use(middleware.Logger())

	// API Routes
	api := r.Group("/api")
	{
		// User routes - only using existing methods
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.POST("", userHandler.CreateUser)
			// Add these methods to UserHandler if you need them:
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Blog routes
		blogs := api.Group("/blogs")
		{
			blogs.POST("", blogHandler.Create)           // POST /api/blogs
			blogs.GET("", blogHandler.GetAll)            // GET /api/blogs (list)
			blogs.GET("/:id", blogHandler.GetByID)       // GET /api/blogs/:id
			blogs.GET("/slug/:slug", blogHandler.GetBySlug) // GET /api/blogs/slug/:slug
			blogs.PUT("/:id", blogHandler.Update)        // PUT /api/blogs/:id (full update)
			blogs.PATCH("/:id", blogHandler.Patch)       // PATCH /api/blogs/:id (partial update)
			blogs.DELETE("/:id", blogHandler.Delete)     // DELETE /api/blogs/:id
		}
	}

	return r
}