package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(rg *gin.RouterGroup, d *Deps) {
	users := rg.Group("/users")
	users.Use(middleware.AuthMiddleware(d.AuthService))
	{
		users.GET("/me", d.UserHandler.GetMe)

		// Admin operations
		admin := users.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.UserHandler.GetAllUsers)
			admin.GET("/:id", d.UserHandler.GetUserByID)
			admin.POST("", d.UserHandler.CreateUser)
			admin.DELETE("/:id", d.UserHandler.DeleteUser)
		}
	}
}
