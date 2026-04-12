package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerPayoutRoutes(rg *gin.RouterGroup, d *Deps) {
	payouts := rg.Group("/payouts")
	payouts.Use(middleware.AuthMiddleware(d.AuthService))
	{
		payouts.GET("/me", d.PayoutHandler.GetMine)

		// Admin operations
		admin := payouts.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.PayoutHandler.GetAll)
			admin.POST("/:id/pay", d.PayoutHandler.MarkPaid)
		}
	}
}
