package handler

import (
    "net/http"
    "strconv"
    "fmt"

    "github.com/cureerel/gotemplate/internal/application/service"
    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}

type orderItemRequest struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type createOrderRequest struct {
    Currency entity.Currency    `json:"currency" binding:"required"`
    Items    []orderItemRequest `json:"items" binding:"required,min=1"`
}

type updateOrderStatusRequest struct {
    Status entity.OrderStatus `json:"status" binding:"required"`
}

func (h *OrderHandler) Create(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req createOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

    items := make([]service.OrderItemInput, len(req.Items))
    for i, item := range req.Items {
        items[i] = service.OrderItemInput{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
        }
    }

    order, err := h.orderService.Create(c.Request.Context(), service.CreateOrderInput{
        UserID:   uint(uid),
        Currency: req.Currency,
        Items:    items,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": order})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    order, err := h.orderService.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    uid, _ := strconv.ParseUint(fmt.Sprintf("%v", userID), 10, 32)

    orders, total, err := h.orderService.GetByUser(c.Request.Context(), uint(uid), page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": orders, "total": total, "page": page, "limit": limit})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req updateOrderStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.orderService.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
} 