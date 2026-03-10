// internal/interfaces/http/handler/payout_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type PayoutHandler struct {
	svc *service.PayoutService
}

func NewPayoutHandler(svc *service.PayoutService) *PayoutHandler {
	return &PayoutHandler{svc: svc}
}

// GET /api/payouts/me — partner+ (own payouts)
func (h *PayoutHandler) GetMine(c *gin.Context) {
	uid, _ := getUID(c)
	page, limit := paginate(c)
	payouts, total, err := h.svc.GetMine(c.Request.Context(), uid, page, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.PayoutResponse, len(payouts))
	for i := range payouts {
		list[i] = toPayoutResponse(&payouts[i])
	}
	c.JSON(http.StatusOK, dto.PayoutListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// GET /api/payouts — admin+
func (h *PayoutHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.Query("status")
	payouts, total, err := h.svc.GetAll(c.Request.Context(), page, limit, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.PayoutResponse, len(payouts))
	for i := range payouts {
		list[i] = toPayoutResponse(&payouts[i])
	}
	c.JSON(http.StatusOK, dto.PayoutListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// POST /api/payouts/:id/pay — admin+
func (h *PayoutHandler) MarkPaid(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.MarkPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, _ := getUID(c)
	if err := h.svc.MarkPaid(c.Request.Context(), id, uid, req.Reference); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payout marked as paid"})
}

func toPayoutResponse(p *entity.Payout) dto.PayoutResponse {
	r := dto.PayoutResponse{
		ID:          p.ID,
		RecipientID: p.RecipientID,
		Type:        p.Type,
		AmountCents: p.AmountCents,
		AmountUSD:   p.AmountUSD(),
		Status:      p.Status,
		OrderID:     p.OrderID,
		Reference:   p.Reference,
		PaidBy:      p.PaidBy,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if p.PaidAt != nil {
		s := p.PaidAt.Format("2006-01-02T15:04:05Z")
		r.PaidAt = &s
	}
	return r
}