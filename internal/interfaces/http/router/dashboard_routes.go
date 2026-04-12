package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerDashboardRoutes(rg *gin.RouterGroup, d *Deps) {
	dash := rg.Group("/dashboard")
	dash.Use(middleware.AuthMiddleware(d.AuthService))
	{
		dash.GET("", d.DashboardHandler.Get)
		dash.GET("/user", d.DashboardHandler.UserView)
		
		dash.GET("/writer", middleware.WriterContent(), d.DashboardHandler.WriterView)
		dash.GET("/partner", middleware.RoleMiddleware(entity.RolePartner), d.DashboardHandler.PartnerView)
		dash.GET("/worker", middleware.RoleMiddleware(entity.RoleWorker), d.DashboardHandler.WorkerView)
		dash.GET("/admin", middleware.RoleMiddleware(entity.RoleAdmin), d.DashboardHandler.AdminView)
	}
}
