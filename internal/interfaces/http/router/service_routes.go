package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerServiceRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public service reads
	pub := rg.Group("/services")
	pub.Use(middleware.OptionalAuth(d.AuthService))
	{
		pub.GET("", d.ServiceHandler.GetAll)
		pub.GET("/:id", d.ServiceHandler.GetByID)
	}

	// Authenticated service actions
	auth := rg.Group("/services")
	auth.Use(middleware.AuthMiddleware(d.AuthService))
	{
		auth.GET("/mine", d.ServiceHandler.GetMine)

		// Partner actions
		partner := auth.Group("")
		partner.Use(middleware.RoleMiddleware(entity.RolePartner))
		{
			partner.POST("", d.ServiceHandler.Create)
			partner.PUT("/:id", d.ServiceHandler.Update)
			partner.DELETE("/:id", d.ServiceHandler.Delete)
			partner.POST("/:id/live", d.ServiceHandler.SetLive)
			partner.POST("/:id/pause", d.ServiceHandler.Pause)
		}

		// Admin actions
		admin := auth.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.POST("/:id/approve", d.ServiceHandler.Approve)
			admin.POST("/:id/reject", d.ServiceHandler.Reject)
		}
	}
}
