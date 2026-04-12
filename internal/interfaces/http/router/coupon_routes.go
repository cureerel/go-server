package router

import (
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func registerCouponRoutes(rg *gin.RouterGroup, d *Deps) {
	// Public validation
	rg.GET("/coupons/validate", d.CouponHandler.Validate)

	coupons := rg.Group("/coupons")
	coupons.Use(middleware.AuthMiddleware(d.AuthService))
	{
		// Partner operations
		partner := coupons.Group("")
		partner.Use(middleware.RoleMiddleware(entity.RolePartner))
		{
			partner.POST("", d.CouponHandler.Create)
			partner.GET("/:id", d.CouponHandler.GetByID)
		}

		// Admin operations
		admin := coupons.Group("")
		admin.Use(middleware.RoleMiddleware(entity.RoleAdmin))
		{
			admin.GET("", d.CouponHandler.GetAll)
			admin.POST("/:id/approve", d.CouponHandler.Approve)
			admin.POST("/:id/reject", d.CouponHandler.Reject)
		}
	}
}
