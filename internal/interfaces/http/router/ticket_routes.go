package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerTicketRoutes(rg *gin.RouterGroup, d *Deps) {
	tickets := rg.Group("/tickets")
	tickets.Use(middleware.AuthMiddleware(d.AuthService))
	{
		// User: own ticket actions
		tickets.POST("", d.TicketHandler.Create)
		tickets.GET("/me", d.TicketHandler.GetMine)
		tickets.POST("/:id/close", d.TicketHandler.Close)
		tickets.POST("/:id/messages", d.TicketHandler.SendMessage)
		tickets.GET("/:id/messages", d.TicketHandler.GetMessages)

		// Admin: see all tickets and resolve
		admin := tickets.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.TicketHandler.GetAll)
			admin.POST("/:id/resolve", d.TicketHandler.Resolve)
		}
	}
}
