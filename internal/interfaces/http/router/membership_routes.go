package router

import (
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerMembershipRoutes(rg *gin.RouterGroup, d *Deps) {
	memberships := rg.Group("/memberships")
	memberships.Use(middleware.AuthMiddleware(d.AuthService))
	{
		memberships.GET("/me", d.MembershipHandler.GetMine)
		memberships.POST("/activate", d.MembershipHandler.Activate)
		memberships.POST("/upgrade", d.MembershipHandler.Upgrade)
		memberships.DELETE("/cancel", d.MembershipHandler.Cancel)
	}
}
