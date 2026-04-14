// internal/infrastructure/postgres/models/order_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Order struct {
	ID             uint        `gorm:"primaryKey"`
	UserID         uint        `gorm:"not null;index"`
	Status         string      `gorm:"default:'in_cart';size:20;index"`
	DeliveryStatus string      `gorm:"column:delivery_status;size:20;default:'created'"`
	TotalCents     int64       `gorm:"column:total_cents;not null"`
	Currency       string      `gorm:"size:10;default:'USD'"`
	CouponID       *uint       `gorm:"column:coupon_id"`
	AffiliateID    *uint       `gorm:"column:affiliate_id"`
	PaymentID      string      `gorm:"column:payment_id;size:100"`
	Items          []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Order) TableName() string { return "orders" }

func (m *Order) ToDomain() *entity.Order {
	o := &entity.Order{
		ID:             m.ID,
		UserID:         m.UserID,
		Status:         entity.OrderStatus(m.Status),
		DeliveryStatus: entity.DeliveryStatus(m.DeliveryStatus),
		TotalCents:     m.TotalCents,
		Currency:       m.Currency,
		CouponID:       m.CouponID,
		AffiliateID:    m.AffiliateID,
		PaymentID:      m.PaymentID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	for _, item := range m.Items {
		o.Items = append(o.Items, *item.ToDomain())
	}
	return o
}

func OrderFromDomain(e *entity.Order) *Order {
	m := &Order{
		ID:             e.ID,
		UserID:         e.UserID,
		Status:         string(e.Status),
		DeliveryStatus: string(e.DeliveryStatus),
		TotalCents:     e.TotalCents,
		Currency:       e.Currency,
		CouponID:       e.CouponID,
		AffiliateID:    e.AffiliateID,
		PaymentID:      e.PaymentID,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
	for _, item := range e.Items {
		m.Items = append(m.Items, *OrderItemFromDomain(&item))
	}
	return m
}

// OrderItem

type OrderItem struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   uint   `gorm:"not null;index"`
	ProductID *uint  `gorm:"column:product_id"`
	ServiceID *uint  `gorm:"column:service_id"`
	Title     string `gorm:"not null;size:200"`
	Quantity  int    `gorm:"default:1"`
	UnitPrice int64  `gorm:"column:unit_price;not null"`
	CreatedAt time.Time
}

func (OrderItem) TableName() string { return "order_items" }

func (m *OrderItem) ToDomain() *entity.OrderItem {
	return &entity.OrderItem{
		ID:        m.ID,
		OrderID:   m.OrderID,
		ProductID: m.ProductID,
		ServiceID: m.ServiceID,
		Title:     m.Title,
		Quantity:  m.Quantity,
		UnitPrice: m.UnitPrice,
	}
}

func OrderItemFromDomain(e *entity.OrderItem) *OrderItem {
	return &OrderItem{
		ID:        e.ID,
		OrderID:   e.OrderID,
		ProductID: e.ProductID,
		ServiceID: e.ServiceID,
		Title:     e.Title,
		Quantity:  e.Quantity,
		UnitPrice: e.UnitPrice,
	}
}
