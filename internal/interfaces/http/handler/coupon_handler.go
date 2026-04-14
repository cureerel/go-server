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

// GET /api/coupons — admin: list all
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
}

// GET /api/coupons/:id — admin
func (h *CouponHandler) GetByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	coupon, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	respond(c, toCouponResponse(coupon))
}

// POST /api/coupons — admin creates coupon
func (h *CouponHandler) Create(c *gin.Context) {
	var req dto.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	coupon, err := h.svc.Create(c.Request.Context(), service.CreateCouponInput{
		CreatorID:        uid,
		Code:             req.Code,
		Type:             req.Type,
		DiscountUSDCents: req.DiscountUSDCents,
		MaxDiscountCents: req.MaxDiscountCents,
		CommissionPct:    req.CommissionPct,
		UsageLimit:       req.UsageLimit,
		ExpiresAt:        req.ExpiresAt,
	})
	if err != nil {
		respondErr(c, err)
		return
	}
	respondCreated(c, toCouponResponse(coupon))
}

// DELETE /api/coupons/:id — admin
func (h *CouponHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted"})
}

// GET /api/coupons/validate?code=XXX — public checkout validation
func (h *CouponHandler) Validate(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}
	coupon, err := h.svc.Validate(c.Request.Context(), code)
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ValidateCouponResponse{
		Valid:            true,
		Code:             coupon.Code,
		Type:             coupon.Type,
		DiscountUSDCents: coupon.DiscountUSDCents,
		CommissionPct:    coupon.CommissionPct,
	})
}

// mapper

func toCouponResponse(c *entity.Coupon) dto.CouponResponse {
	return dto.CouponResponse{
		ID:               c.ID,
		CreatorID:        c.CreatorID,
		Code:             c.Code,
		Type:             c.Type,
		DiscountUSDCents: c.DiscountUSDCents,
		MaxDiscountCents: c.MaxDiscountCents,
		CommissionPct:    c.CommissionPct,
		Status:           c.Status,
		UsageLimit:       c.UsageLimit,
		UsedCount:        c.UsedCount,
		ExpiresAt:        c.ExpiresAt,
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}