package router

import (
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup, d *Deps) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register/init", d.AuthHandler.RegisterInit)
		auth.POST("/register/verify", d.AuthHandler.RegisterVerify)
		auth.POST("/password/reset/init", d.AuthHandler.PasswordResetInit)
		auth.POST("/password/reset/verify", d.AuthHandler.PasswordResetVerify)
		auth.POST("/signup", d.AuthHandler.Signup)
		auth.POST("/login", d.AuthHandler.Login)
		auth.POST("/refresh", d.AuthHandler.Refresh)
		
		// Authenticated logout
		auth.Use(middleware.AuthMiddleware(d.AuthService))
		auth.POST("/logout", d.AuthHandler.Logout)
	}
}
