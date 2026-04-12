package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerProductRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public product reads
	pub := rg.Group("/products")
	{
		pub.GET("", d.ProductHandler.GetAll)
		pub.GET("/:id", d.ProductHandler.GetByID)
	}

	// Admin product management
	admin := rg.Group("/products")
	admin.Use(middleware.AuthMiddleware(d.AuthService))
	admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		admin.POST("", d.ProductHandler.Create)
		admin.PUT("/:id", d.ProductHandler.Update)
		admin.DELETE("/:id", d.ProductHandler.Delete)
	}
}
