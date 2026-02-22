package handler

import (
	"net/http"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// GET /payments/:id — admin or owner only
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	payment, err := h.paymentService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toPaymentResponse(payment)})
}

// POST /payments/:id/refund — admin only
func (h *PaymentHandler) Refund(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.paymentService.Refund(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, _ := h.paymentService.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": toPaymentResponse(payment), "message": "payment refunded"})
}

// POST /payments/:id/complete — internal/webhook use
func (h *PaymentHandler) MarkCompleted(c *gin.Context) {
	id := c.Param("id")
	if err := h.paymentService.MarkCompleted(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, _ := h.paymentService.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": toPaymentResponse(payment), "message": "payment marked completed"})
}

// POST /payments/:id/fail — internal/webhook use
func (h *PaymentHandler) MarkFailed(c *gin.Context) {
	id := c.Param("id")
	if err := h.paymentService.MarkFailed(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, _ := h.paymentService.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"data": toPaymentResponse(payment), "message": "payment marked failed"})
}

// helper to map entity.Payment -> dto.PaymentResponse
func toPaymentResponse(p *entity.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:            p.ID,
		UserID:        p.UserID,
		OrderID:       p.OrderID,
		Amount:        p.Amount,
		Currency:      string(p.Currency),
		Status:        string(p.Status),
		Provider:      string(p.Provider),
		ProviderTxnID: p.ProviderTxnID,
		CustomerEmail: p.CustomerEmail,
		Description:   p.Description,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}