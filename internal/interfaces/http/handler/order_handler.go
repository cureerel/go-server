// internal/interfaces/http/handler/order_handler.go
package handler

import (
	"net/http"

	"github.com/cureerel/gotemplate/internal/application/service"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderSvc   *service.OrderService
	paymentSvc *service.PaymentService
}

func NewOrderHandler(orderSvc *service.OrderService, paymentSvc *service.PaymentService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc, paymentSvc: paymentSvc}
}

// POST /api/orders — any authenticated user
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, ok := getUID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	order, err := h.orderSvc.Create(c.Request.Context(), service.CreateOrderInput{
		UserID:      uid,
		ServiceID:   req.ServiceID,
		Provider:    req.Provider,
		AffiliateID: req.AffiliateID,
		// CouponID resolved in a later phase (phase 7 coupons)
	})
	if err != nil {
		respondErr(c, err)
		return
	}

	// Create pending payment record immediately after order
	userEmail, _ := c.Get("email")
	email, _ := userEmail.(string)
	payment, err := h.paymentSvc.InitPayment(
		c.Request.Context(),
		order,
		entity.PaymentProvider(req.Provider),
		email,
	)
	if err != nil {
		respondErr(c, err)
		return
	}

	respondCreated(c, gin.H{
		"order":   toOrderResponse(order),
		"payment": toPaymentResponse(payment),
	})
}

// GET /api/orders/me — any authenticated user
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	uid, _ := getUID(c)
	page, limit := paginate(c)
	orders, total, err := h.orderSvc.GetMyOrders(c.Request.Context(), uid, page, limit)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		list[i] = toOrderResponse(&orders[i])
	}
	c.JSON(http.StatusOK, dto.OrderListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// GET /api/orders/:id — owner or admin
func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	order, err := h.orderSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		respondErr(c, err)
		return
	}
	uid, _ := getUID(c)
	if order.UserID != uid && !hasRole(c, entity.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	respond(c, toOrderResponse(order))
}

// GET /api/orders — admin+
func (h *OrderHandler) GetAll(c *gin.Context) {
	page, limit := paginate(c)
	status := c.Query("status")
	orders, total, err := h.orderSvc.GetAll(c.Request.Context(), page, limit, status)
	if err != nil {
		respondErr(c, err)
		return
	}
	list := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		list[i] = toOrderResponse(&orders[i])
	}
	c.JSON(http.StatusOK, dto.OrderListResponse{Data: list, Total: total, Page: page, Limit: limit})
}

// PATCH /api/orders/:id/status — admin+
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondErr(c, err)
		return
	}
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.orderSvc.UpdateStatus(c.Request.Context(), id, entity.OrderStatus(req.Status)); err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

// ── mapper ────────────────────────────────────────────────────

func toOrderResponse(o *entity.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = dto.OrderItemResponse{
			ID:        item.ID,
			ServiceID: item.ServiceID,
			Title:     item.Title,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			UnitUSD:   float64(item.UnitPrice) / 100,
		}
	}
	return dto.OrderResponse{
		ID:              o.ID,
		UserID:          o.UserID,
		ServiceID:       o.ServiceID,
		Status:          string(o.Status),
		TotalCents:      o.TotalCents,
		TotalUSD:        o.TotalUSD(),
		Currency:        o.Currency,
		PaymentProvider: o.PaymentProvider,
		Items:           items,
		CreatedAt:       o.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}