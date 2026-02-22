package handler

import (
    "net/http"

    "github.com/cureerel/gotemplate/internal/application/service"
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
    c.JSON(http.StatusOK, gin.H{"data": payment})
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
    c.JSON(http.StatusOK, gin.H{"message": "payment refunded"})
}

// POST /payments/:id/complete — internal/webhook use
func (h *PaymentHandler) MarkCompleted(c *gin.Context) {
    id := c.Param("id")
    if err := h.paymentService.MarkCompleted(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "payment marked completed"})
}

// POST /payments/:id/fail — internal/webhook use
func (h *PaymentHandler) MarkFailed(c *gin.Context) {
    id := c.Param("id")
    if err := h.paymentService.MarkFailed(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "payment marked failed"})
}