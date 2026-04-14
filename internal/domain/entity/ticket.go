// internal/domain/entity/ticket.go
package entity

import "time"

const (
	TicketOpen       = "open"
	TicketInProgress = "in_progress"
	TicketResolved   = "resolved"
	TicketClosed     = "closed"

	TicketPriorityLow    = "low"
	TicketPriorityMedium = "medium"
	TicketPriorityHigh   = "high"
	TicketPriorityUrgent = "urgent"
)

type Ticket struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssignedTo  *uint      `json:"assigned_to,omitempty"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *Ticket) IsClosed() bool { return t.Status == TicketClosed }

type TicketMessage struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	SenderID  uint      `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
