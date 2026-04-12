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
		tickets.POST("", d.TicketHandler.Create)
		tickets.GET("/me", d.TicketHandler.GetMine)
		tickets.GET("/:id", d.TicketHandler.GetByID)
		tickets.POST("/:id/close", d.TicketHandler.Close)
		tickets.POST("/:id/messages", d.TicketHandler.SendMessage)
		tickets.GET("/:id/messages", d.TicketHandler.GetMessages)

		// Worker operations
		worker := tickets.Group("")
		worker.Use(middleware.RoleMiddleware(entity.RoleWorker))
		{
			worker.GET("", d.TicketHandler.GetAll)
			worker.POST("/:id/assign", d.TicketHandler.Assign)
			worker.POST("/:id/resolve", d.TicketHandler.Resolve)
		}
	}
}
