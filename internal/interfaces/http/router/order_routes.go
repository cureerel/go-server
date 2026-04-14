package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerOrderRoutes(rg *gin.RouterGroup, d *Deps) {
	orders := rg.Group("/orders")
	orders.Use(middleware.AuthMiddleware(d.AuthService))
	{
		orders.POST("", d.OrderHandler.Create)
		orders.GET("/me", d.OrderHandler.GetMyOrders)
		orders.GET("/:id", d.OrderHandler.GetByID)

		// Admin
		admin := orders.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.OrderHandler.GetAll)
			admin.PATCH("/:id/delivery", d.OrderHandler.UpdateDeliveryStatus)
		}
	}
}
