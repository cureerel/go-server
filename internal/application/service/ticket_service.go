// internal/application/service/ticket_service.go
package service

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/pkg/apperror"
)

type TicketService struct {
	ticketRepo repository.TicketRepository
}

func NewTicketService(ticketRepo repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

type CreateTicketInput struct {
	UserID      uint
	Subject     string
	Description string
	Priority    string
}

func (s *TicketService) Create(ctx context.Context, in CreateTicketInput) (*entity.Ticket, error) {
	priority := in.Priority
	if !validPriority(priority) {
		priority = entity.TicketPriorityMedium
	}
	t := &entity.Ticket{
		UserID:      in.UserID,
		Subject:     in.Subject,
		Description: in.Description,
		Status:      entity.TicketOpen,
		Priority:    priority,
	}
	if err := s.ticketRepo.Create(ctx, t); err != nil {
		return nil, apperror.NewInternal(err, "failed to create ticket")
	}
	return t, nil
}

func (s *TicketService) GetByID(ctx context.Context, id uint) (*entity.Ticket, error) {
	t, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch ticket")
	}
	if t == nil {
		return nil, apperror.NewNotFound("ticket not found")
	}
	return t, nil
}

func (s *TicketService) GetMine(ctx context.Context, userID uint, page, limit int, status string) ([]entity.Ticket, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.ticketRepo.GetAll(ctx, repository.TicketFilter{
		Page: page, Limit: limit, UserID: &userID, Status: status,
	})
}

func (s *TicketService) GetAll(ctx context.Context, page, limit int, status, priority string) ([]entity.Ticket, int64, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	return s.ticketRepo.GetAll(ctx, repository.TicketFilter{
		Page: page, Limit: limit, Status: status, Priority: priority,
	})
}



func (s *TicketService) Resolve(ctx context.Context, ticketID uint) error {
	t, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil || t == nil {
		return apperror.NewNotFound("ticket not found")
	}
	if t.IsClosed() {
		return apperror.NewBadRequest("ticket is already closed")
	}
	t.Status = entity.TicketResolved
	return s.ticketRepo.Update(ctx, t)
}

func (s *TicketService) Close(ctx context.Context, ticketID, callerID uint, callerRole string) error {
	t, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil || t == nil {
		return apperror.NewNotFound("ticket not found")
	}
	if t.IsClosed() {
		return nil
	}
	u := &entity.User{Role: callerRole}
	if t.UserID != callerID && !u.HasRole(entity.RoleAdmin) {
		return apperror.NewForbidden("you don't own this ticket")
	}
	return s.ticketRepo.Close(ctx, ticketID)
}

func (s *TicketService) SendMessage(ctx context.Context, ticketID, senderID uint, message string) (*entity.TicketMessage, error) {
	t, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil || t == nil {
		return nil, apperror.NewNotFound("ticket not found")
	}
	if t.IsClosed() {
		return nil, apperror.NewBadRequest("cannot message on a closed ticket")
	}
	msg := &entity.TicketMessage{
		TicketID: ticketID,
		SenderID: senderID,
		Message:  message,
	}
	if err := s.ticketRepo.CreateMessage(ctx, msg); err != nil {
		return nil, apperror.NewInternal(err, "failed to send message")
	}
	return msg, nil
}

func (s *TicketService) GetMessages(ctx context.Context, ticketID, callerID uint, callerRole string) ([]entity.TicketMessage, error) {
	t, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil || t == nil {
		return nil, apperror.NewNotFound("ticket not found")
	}
	u := &entity.User{Role: callerRole}
	if t.UserID != callerID && !u.HasRole(entity.RoleWorker) {
		return nil, apperror.NewForbidden("access denied")
	}
	msgs, err := s.ticketRepo.GetMessages(ctx, ticketID)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch messages")
	}
	return msgs, nil
}

func validPriority(p string) bool {
	switch p {
	case entity.TicketPriorityLow, entity.TicketPriorityMedium,
		entity.TicketPriorityHigh, entity.TicketPriorityUrgent:
		return true
	}
	return false
}