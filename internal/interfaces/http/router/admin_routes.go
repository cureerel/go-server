package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerSuperAdminRoutes(rg *gin.RouterGroup, d *Deps) {
	auth := rg.Group("")
	auth.Use(middleware.AuthMiddleware(d.AuthService))
	{
	

		sa := auth.Group("/superadmin")
		sa.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			sa.GET("/stats", d.SuperAdminHandler.Stats)
			
		}
	}
}
