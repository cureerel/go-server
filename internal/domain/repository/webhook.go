package repository

import (
	"context"
	"github.com/cureerel/gotemplate/internal/domain/entity"
)

type WebhookRepository interface {
	StoreEvent(ctx context.Context, event *entity.WebhookEvent) error
	GetEventByID(ctx context.Context, id string) (*entity.WebhookEvent, error)
	MarkProcessed(ctx context.Context, id string) error
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	UpdatePaymentStatus(ctx context.Context, id, status string) error
}