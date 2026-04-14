// internal/domain/repository/ticket.go
package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type TicketFilter struct {
	Page     int
	Limit    int
	Status   string
	UserID   *uint
	Priority string
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *entity.Ticket) error
	GetByID(ctx context.Context, id uint) (*entity.Ticket, error)
	GetAll(ctx context.Context, filter TicketFilter) ([]entity.Ticket, int64, error)
	Update(ctx context.Context, ticket *entity.Ticket) error
	Close(ctx context.Context, id uint) error
	Assign(ctx context.Context, id uint, workerID uint) error
	CreateMessage(ctx context.Context, msg *entity.TicketMessage) error
	GetMessages(ctx context.Context, ticketID uint) ([]entity.TicketMessage, error)
}
