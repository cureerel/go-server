// internal/infrastructure/postgres/repositories/ticket_repository.go
package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type ticketRepository struct{ db *gorm.DB }

func NewTicketRepository(db *gorm.DB) repository.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, t *entity.Ticket) error {
	m := models.TicketFromDomain(t)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id uint) (*entity.Ticket, error) {
	var m models.Ticket
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *ticketRepository) GetAll(ctx context.Context, filter repository.TicketFilter) ([]entity.Ticket, int64, error) {
	var ms []models.Ticket
	var total int64
	offset := (filter.Page - 1) * filter.Limit
	q := r.db.WithContext(ctx).Model(&models.Ticket{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.Priority != "" {
		q = q.Where("priority = ?", filter.Priority)
	}
	q = q.Order("created_at DESC")
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	tickets := make([]entity.Ticket, len(ms))
	for i, m := range ms {
		tickets[i] = *m.ToDomain()
	}
	return tickets, total, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *entity.Ticket) error {
	return r.db.WithContext(ctx).Save(models.TicketFromDomain(ticket)).Error
}

func (r *ticketRepository) Close(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", id).Updates(map[string]any{
		"status":    entity.TicketClosed,
		"closed_at": now,
	}).Error
}

func (r *ticketRepository) Assign(ctx context.Context, id uint, workerID uint) error {
	return r.db.WithContext(ctx).Model(&models.Ticket{}).Where("id = ?", id).Updates(map[string]any{
		"assigned_to": workerID,
		"status":      entity.TicketInProgress,
	}).Error
}

func (r *ticketRepository) CreateMessage(ctx context.Context, msg *entity.TicketMessage) error {
	m := models.TicketMessageFromDomain(msg)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	return nil
}

func (r *ticketRepository) GetMessages(ctx context.Context, ticketID uint) ([]entity.TicketMessage, error) {
	var ms []models.TicketMessage
	if err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		Order("created_at ASC").
		Find(&ms).Error; err != nil {
		return nil, err
	}
	msgs := make([]entity.TicketMessage, len(ms))
	for i, m := range ms {
		msgs[i] = *m.ToDomain()
	}
	return msgs, nil
}