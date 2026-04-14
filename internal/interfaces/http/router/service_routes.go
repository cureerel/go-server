package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerServiceRoutes(rg *gin.RouterGroup, d *Deps) {
	services := rg.Group("/services")

	// Public
	services.GET("", d.ServiceHandler.GetAll)

	// Admin only
	admin := services.Group("")
	admin.Use(middleware.AuthMiddleware(d.AuthService))
	admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		admin.POST("", d.ServiceHandler.Create)
		admin.DELETE("/:id", d.ServiceHandler.Delete)
		admin.POST("/:id/live", d.ServiceHandler.SetLive)
		admin.POST("/:id/pause", d.ServiceHandler.Pause)
	}
}
