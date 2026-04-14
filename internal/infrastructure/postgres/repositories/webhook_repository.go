package repositories

import (
	"context"
	"errors"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type webhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) SaveEvent(ctx context.Context, event *entity.WebhookEvent) error {
	m := models.WebhookEventFromDomain(event)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *webhookRepository) GetEventByID(ctx context.Context, id string) (*entity.WebhookEvent, error) {
	var m models.WebhookEvent
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *webhookRepository) MarkProcessed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&models.WebhookEvent{}).
		Where("id = ?", id).
		Update("processed", true).Error
}

func (r *webhookRepository) SavePayment(ctx context.Context, payment *entity.Payment) error {
	m := models.PaymentFromDomain(payment)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *webhookRepository) GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error) {
	var m models.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *webhookRepository) UpdatePaymentStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("id = ?", id).
		Update("status", status).Error
}
