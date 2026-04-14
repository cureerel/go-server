package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerCouponRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public: validate a coupon code at checkout
	rg.GET("/coupons/validate", d.CouponHandler.Validate)

	// Admin only
	coupons := rg.Group("/coupons")
	coupons.Use(middleware.AuthMiddleware(d.AuthService))
	coupons.Use(middleware.RoleMiddleware(entity.RoleAdmin))
	{
		coupons.GET("", d.CouponHandler.GetAll)
		coupons.POST("", d.CouponHandler.Create)
		coupons.GET("/:id", d.CouponHandler.GetByID)
		coupons.DELETE("/:id", d.CouponHandler.Delete)
	}
}
