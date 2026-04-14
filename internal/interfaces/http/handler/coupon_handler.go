// internal/interfaces/http/handler/coupon_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc *service.CouponService
}

func NewCouponHandler(svc *service.CouponService) *CouponHandler {
	return &CouponHandler{svc: svc}
}






// GET /api/coupons —
func (h *CouponHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.Query("status")
	coupons, total, err := h.svc.GetAll(c.Request.Context(), page, limit, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.CouponResponse, len(coupons))
	for i := range coupons {
		list[i] = toCouponResponse(&coupons[i])
	}
	c.JSON(http.StatusOK, dto.CouponListResponse{Data: list, Total: total, Page: page, Limit: limit})




func t\oCouponResponse(c *entity.Coupon) dto.CouponResponse {
	return dto.CouponResponse{
		ID:               c.ID,
		CreatorID:        c.CreatorID,
		Code:             c.Code,
		Type:             c.Type,
		DiscountUSDCents: c.DiscountUSDCents,
		MaxDiscountCents: c.MaxDiscountCents,

		Status:           c.Status,
		UsageLimit:       c.UsageLimit,
		UsedCount:        c.UsedCount,
		ExpiresAt:        c.ExpiresAt,
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}