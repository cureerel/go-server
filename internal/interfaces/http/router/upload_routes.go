package router

import (
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerUploadRoutes(rg *gin.RouterGroup, d *Deps) {
	upload := rg.Group("/upload")
	upload.Use(middleware.AuthMiddleware(d.AuthService))
	{
		upload.POST("/image", d.UploadHandler.UploadImage)
		upload.DELETE("/image", d.UploadHandler.DeleteImage)
	}
}
