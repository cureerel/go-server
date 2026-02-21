package persistence

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
	"github.com/cureerel/gotemplate/internal/domain/repository"
	"gorm.io/gorm"
)

type webhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) StoreEvent(ctx context.Context, event *entity.WebhookEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *webhookRepository) GetEventByID(ctx context.Context, id string) (*entity.WebhookEvent, error) {
	var event entity.WebhookEvent
	err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error
	return &event, err
}

func (r *webhookRepository) MarkProcessed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entity.WebhookEvent{}).
		Where("id = ?", id).Update("processed", true).Error
}

func (r *webhookRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *webhookRepository) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Payment{}).
		Where("id = ?", id).Update("status", status).Error
}