package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerAdminRoutes(rg *gin.RouterGroup, d *Deps) {
	auth := rg.Group("")
	auth.Use(middleware.AuthMiddleware(d.AuthService))
	{
		admin := auth.Group("/admin")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("/stats", d.AdminHandler.Stats)
		}
	}
}
