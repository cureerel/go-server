package handler

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/cureerel/gotemplate/internal/application/service"
    "github.com/cureerel/gotemplate/internal/domain/entity"
    "github.com/cureerel/gotemplate/internal/interfaces/dto"
    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Create(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req dto.CreateOrderRequest
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
        Currency: entity.Currency(req.Currency),
        Items:    items,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"data": toOrderResponse(order)})
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
    c.JSON(http.StatusOK, gin.H{"data": toOrderResponse(order)})
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

    result := make([]dto.OrderResponse, len(orders))
    for i, o := range orders {
        result[i] = toOrderResponse(&o)
    }

    c.JSON(http.StatusOK, dto.OrderListResponse{
        Data:  result,
        Total: total,
        Page:  page,
        Limit: limit,
    })
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req dto.UpdateOrderStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.orderService.UpdateStatus(c.Request.Context(), uint(id), entity.OrderStatus(req.Status)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

func toOrderResponse(o *entity.Order) dto.OrderResponse {
    items := make([]dto.OrderItemResponse, len(o.Items))
    for i, item := range o.Items {
        items[i] = dto.OrderItemResponse{
            ID:        item.ID,
            ProductID: item.ProductID,
            Type:      string(item.Type),
            Quantity:  item.Quantity,
            UnitPrice: item.UnitPrice,
        }
    }
    return dto.OrderResponse{
        ID:          o.ID,
        UserID:      o.UserID,
        Status:      string(o.Status),
        TotalAmount: o.TotalAmount,
        Currency:    string(o.Currency),
        Items:       items,
        CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z"),
        UpdatedAt:   o.UpdatedAt.Format("2006-01-02T15:04:05Z"),
    }
}