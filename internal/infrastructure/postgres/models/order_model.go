package models

import (
    "time"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type Order struct {
    ID          uint        `gorm:"primaryKey"`
    UserID      uint        `gorm:"not null;index"`
    Status      string      `gorm:"not null;size:20;default:'pending'"`
    TotalAmount int64       `gorm:"not null"`
    Currency    string      `gorm:"not null;size:10"`
    Items       []OrderItem `gorm:"foreignKey:OrderID"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (Order) TableName() string {
    return "orders"
}

type OrderItem struct {
    ID        uint  `gorm:"primaryKey"`
    OrderID   uint  `gorm:"not null;index"`
    ProductID uint  `gorm:"not null;index"`
    Quantity  int   `gorm:"not null"`
    UnitPrice int64 `gorm:"not null"`
}

func (OrderItem) TableName() string {
    return "order_items"
}

func (m *Order) ToDomain() *entity.Order {
    items := make([]entity.OrderItem, len(m.Items))
    for i, item := range m.Items {
        items[i] = entity.OrderItem{
            ID:        item.ID,
            OrderID:   item.OrderID,
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            UnitPrice: item.UnitPrice,
        }
    }
    return &entity.Order{
        ID:          m.ID,
        UserID:      m.UserID,
        Status:      entity.OrderStatus(m.Status),
        TotalAmount: m.TotalAmount,
        Currency:    entity.Currency(m.Currency),
        Items:       items,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
    }
}

func OrderFromDomain(e *entity.Order) *Order {
    items := make([]OrderItem, len(e.Items))
    for i, item := range e.Items {
        items[i] = OrderItem{
            ID:        item.ID,
            OrderID:   item.OrderID,
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            UnitPrice: item.UnitPrice,
        }
    }
    return &Order{
        ID:          e.ID,
        UserID:      e.UserID,
        Status:      string(e.Status),
        TotalAmount: e.TotalAmount,
        Currency:    string(e.Currency),
        Items:       items,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,
    }
}