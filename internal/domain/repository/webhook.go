package repository

import (
    "context"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type WebhookRepository interface {
    SaveEvent(ctx context.Context, event *entity.WebhookEvent) error
    GetEventByID(ctx context.Context, id string) (*entity.WebhookEvent, error)
    MarkProcessed(ctx context.Context, id string) error
    SavePayment(ctx context.Context, payment *entity.Payment) error
    GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error)
    UpdatePaymentStatus(ctx context.Context, id string, status string) error
}