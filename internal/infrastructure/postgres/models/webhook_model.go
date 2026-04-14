// internal/infrastructure/postgres/models/webhook_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type WebhookEvent struct {
	ID        string `gorm:"primaryKey;size:100"`
	Provider  string `gorm:"not null;size:50"`
	EventType string `gorm:"not null;size:100"`
	Payload   []byte `gorm:"type:bytea"`
	Signature string `gorm:"size:255"`
	Processed bool   `gorm:"default:false"`
	CreatedAt time.Time
}

func (WebhookEvent) TableName() string { return "webhook_events" }

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
