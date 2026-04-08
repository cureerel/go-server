// internal/infrastructure/postgres/models/ticket_model.go
package models

import (
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type Ticket struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null;index"`
	Subject     string     `gorm:"not null;size:200"`
	Description string     `gorm:"type:text;not null"`
	Status      string     `gorm:"default:'open';size:20;index"`
	Priority    string     `gorm:"default:'medium';size:20"`
	AssignedTo  *uint      `gorm:"column:assigned_to"`
	ClosedAt    *time.Time `gorm:"column:closed_at"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Ticket) TableName() string { return "tickets" }

func (m *Ticket) ToDomain() *entity.Ticket {
	return &entity.Ticket{
		ID: m.ID, UserID: m.UserID, Subject: m.Subject,
		Description: m.Description, Status: m.Status, Priority: m.Priority,
		AssignedTo: m.AssignedTo, ClosedAt: m.ClosedAt,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func TicketFromDomain(e *entity.Ticket) *Ticket {
	return &Ticket{
		ID: e.ID, UserID: e.UserID, Subject: e.Subject,
		Description: e.Description, Status: e.Status, Priority: e.Priority,
		AssignedTo: e.AssignedTo, ClosedAt: e.ClosedAt,
		CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
}

// ── TicketMessage ─────────────────────────────────────────────

type TicketMessage struct {
	ID        uint      `gorm:"primaryKey"`
	TicketID  uint      `gorm:"not null;index"`
	SenderID  uint      `gorm:"not null;index"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}

func (TicketMessage) TableName() string { return "ticket_messages" }

func (m *TicketMessage) ToDomain() *entity.TicketMessage {
	return &entity.TicketMessage{
		ID: m.ID, TicketID: m.TicketID,
		SenderID: m.SenderID, Message: m.Message, CreatedAt: m.CreatedAt,
	}
}

func TicketMessageFromDomain(e *entity.TicketMessage) *TicketMessage {
	return &TicketMessage{
		ID: e.ID, TicketID: e.TicketID,
		SenderID: e.SenderID, Message: e.Message, CreatedAt: e.CreatedAt,
	}
}