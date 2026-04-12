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
		upgrades := auth.Group("/upgrade-requests")
		{
			upgrades.POST("", d.SuperAdminHandler.RequestUpgrade)
			upgrades.GET("/me", d.SuperAdminHandler.GetMyUpgradeRequest)
		}

		sa := auth.Group("/superadmin")
		sa.Use(middleware.RoleMiddleware(entity.RoleSuperAdmin))
		{
			sa.GET("/stats", d.SuperAdminHandler.Stats)
			sa.PATCH("/users/:id/role", d.SuperAdminHandler.SetRole)
			sa.GET("/upgrades", d.SuperAdminHandler.GetPendingUpgrades)
			sa.POST("/upgrades/:id/review", d.SuperAdminHandler.ReviewUpgrade)
		}
	}
}
