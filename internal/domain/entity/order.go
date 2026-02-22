package entity

import "time"

type OrderStatus string

const (
    OrderPending    OrderStatus = "pending"
    OrderConfirmed  OrderStatus = "confirmed"
    OrderDispatched OrderStatus = "dispatched"
    OrderDelivered  OrderStatus = "delivered"
    OrderCancelled  OrderStatus = "cancelled"
    OrderCompleted  OrderStatus = "completed"
)

type Order struct {
    ID          uint        `json:"id"`
    UserID      uint        `json:"user_id"`
    Status      OrderStatus `json:"status"`
    TotalAmount int64       `json:"total_amount"`
    Currency    Currency    `json:"currency"`
    Items       []OrderItem `json:"items"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
    ID        uint        `json:"id"`
    OrderID   uint        `json:"order_id"`
    ProductID uint        `json:"product_id"`
    Type      ProductType `json:"type"`
    Quantity  int         `json:"quantity"`
    UnitPrice int64       `json:"unit_price"`
}

func (o *Order) CalculateTotal() int64 {
    var total int64
    for _, item := range o.Items {
        total += item.UnitPrice * int64(item.Quantity)
    }
    return total
}