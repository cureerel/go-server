package router

import (
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerCoinRoutes(rg *gin.RouterGroup, d *Deps) {
	coins := rg.Group("/coins")
	coins.Use(middleware.AuthMiddleware(d.AuthService))
	{
		coins.GET("/balance", d.CoinHandler.GetBalance)
	}
}
