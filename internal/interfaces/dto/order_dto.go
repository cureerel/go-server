package dto

type OrderItemRequest struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
    Currency string             `json:"currency" binding:"required"`
    Items    []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=pending confirmed dispatched delivered completed cancelled"`
}

type OrderItemResponse struct {
    ID        uint   `json:"id"`
    ProductID uint   `json:"product_id"`
    Type      string `json:"type"`
    Quantity  int    `json:"quantity"`
    UnitPrice int64  `json:"unit_price"`
}

type OrderResponse struct {
    ID          uint                `json:"id"`
    UserID      uint                `json:"user_id"`
    Status      string              `json:"status"`
    TotalAmount int64               `json:"total_amount"`
    Currency    string              `json:"currency"`
    Items       []OrderItemResponse `json:"items"`
    CreatedAt   string              `json:"created_at"`
    UpdatedAt   string              `json:"updated_at"`
}

type OrderListResponse struct {
    Data  []OrderResponse `json:"data"`
    Total int64           `json:"total"`
    Page  int             `json:"page"`
    Limit int             `json:"limit"`
}