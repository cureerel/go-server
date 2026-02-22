package models

import (
    "time"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type WebhookEvent struct {
    ID        string    `gorm:"primaryKey;size:100"`
    Provider  string    `gorm:"not null;size:50"`
    EventType string    `gorm:"not null;size:100"`
    Payload   []byte    `gorm:"type:bytea"`
    Signature string    `gorm:"size:255"`
    Processed bool      `gorm:"default:false"`
    CreatedAt time.Time
}

func (WebhookEvent) TableName() string {
    return "webhook_events"
}

func (m *WebhookEvent) ToDomain() *entity.WebhookEvent {
    return &entity.WebhookEvent{
        ID:        m.ID,
        Provider:  m.Provider,
        EventType: m.EventType,
        Payload:   m.Payload,
        Signature: m.Signature,
        Processed: m.Processed,
        CreatedAt: m.CreatedAt,
    }
}

func WebhookEventFromDomain(e *entity.WebhookEvent) *WebhookEvent {
    return &WebhookEvent{
        ID:        e.ID,
        Provider:  e.Provider,
        EventType: e.EventType,
        Payload:   e.Payload,
        Signature: e.Signature,
        Processed: e.Processed,
        CreatedAt: e.CreatedAt,
    }
}

type Payment struct {
    ID            string    `gorm:"primaryKey;size:100"`
    OrderID       string    `gorm:"not null;size:100;index"`
    Amount        int64     `gorm:"not null"`
    Currency      string    `gorm:"not null;size:10"`
    Status        string    `gorm:"not null;size:20;default:'pending'"`
    Provider      string    `gorm:"not null;size:50"`
    ProviderTxnID string    `gorm:"size:255"`
    CustomerEmail string    `gorm:"not null;size:100"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

func (Payment) TableName() string {
    return "payments"
}

func (m *Payment) ToDomain() *entity.Payment {
    return &entity.Payment{
        ID:            m.ID,
        OrderID:       m.OrderID,
        Amount:        m.Amount,
        Currency:      entity.Currency(m.Currency),
        Status:        entity.PaymentStatus(m.Status),
        Provider:      entity.PaymentProvider(m.Provider),
        ProviderTxnID: m.ProviderTxnID,
        CustomerEmail: m.CustomerEmail,
        CreatedAt:     m.CreatedAt,
        UpdatedAt:     m.UpdatedAt,
    }
}

func PaymentFromDomain(e *entity.Payment) *Payment {
    return &Payment{
        ID:            e.ID,
        OrderID:       e.OrderID,
        Amount:        e.Amount,
        Currency:      string(e.Currency),
        Status:        string(e.Status),
        Provider:      string(e.Provider),
        ProviderTxnID: e.ProviderTxnID,
        CustomerEmail: e.CustomerEmail,
        CreatedAt:     e.CreatedAt,
        UpdatedAt:     e.UpdatedAt,
    }
}