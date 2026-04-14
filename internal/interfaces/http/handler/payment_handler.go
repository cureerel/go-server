// internal/interfaces/http/handler/payment_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/cserver/internal/application/service"
	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// GET /api/payments/:id —  admin
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment id required"})
		return
	}
	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if p.UserID != uid && !hasRole(c, entity.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	respond(c, toPaymentResponse(p))
}

// GET /api/payments — admin
func (h *PaymentHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.Query("status")
	payments, total, err := h.svc.GetAll(c.Request.Context(), page, limit, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.PaymentResponse, len(payments))
	for i := range payments {
		list[i] = toPaymentResponse(&payments[i])
	}
	c.JSON(http.StatusOK, dto.PaymentListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// POST /api/payments/:id/complete — admin
func (h *PaymentHandler) MarkCompleted(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.MarkCompleted(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment marked completed"})
}

// POST /api/payments/:id/fail — admin
func (h *PaymentHandler) MarkFailed(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.MarkFailed(c.Request.Context(), id); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment marked failed"})
}

// POST /api/payments/:id/refund — admin
func (h *PaymentHandler) Refund(c *gin.Context) {
	id := c.Param("id")
	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Refund(c.Request.Context(), id, req.RefundID); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "payment refunded"})
}

//  mapper

func toPaymentResponse(p *entity.Payment) dto.PaymentResponse {
	r := dto.PaymentResponse{
		ID:            p.ID,
		OrderID:       p.OrderID,
		UserID:        p.UserID,
		AmountCents:   p.AmountCents,
		AmountUSD:     p.AmountUSD(),
		Currency:      p.Currency,
		Status:        string(p.Status),
		Provider:      string(p.Provider),
		ProviderTxnID: p.ProviderTxnID,
		CustomerEmail: p.CustomerEmail,
		Description:   p.Description,
		RefundID:      p.RefundID,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if p.RefundedAt != nil {
		s := p.RefundedAt.Format("2006-01-02T15:04:05Z")
		r.RefundedAt = &s
	}
	return r
}
