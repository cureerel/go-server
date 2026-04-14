// internal/infrastructure/postgres/repositories/order_repository.go
package repositories

import (
	"context"
	"errors"

	"github.com/cureerel/cserver/internal/domain/entity"
	"github.com/cureerel/cserver/internal/domain/repository"
	"github.com/cureerel/cserver/internal/infrastructure/postgres/models"
	"gorm.io/gorm"
)

type orderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	m := models.OrderFromDomain(order)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	order.ID = m.ID
	for i, item := range m.Items {
		order.Items[i].ID = item.ID
	}
	return nil
}

func (r *orderRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, order *entity.Order) error {
	m := models.OrderFromDomain(order)
	if err := tx.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	order.ID = m.ID
	for i, item := range m.Items {
		order.Items[i].ID = item.ID
	}
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*entity.Order, error) {
	var m models.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *orderRepository) GetAll(ctx context.Context, filter repository.OrderFilter) ([]entity.Order, int64, error) {
	var ms []models.Order
	var total int64
	offset := (filter.Page - 1) * filter.Limit

	q := r.db.WithContext(ctx).Model(&models.Order{})
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	q = q.Order("created_at DESC")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Preload("Items").Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	orders := make([]entity.Order, len(ms))
	for i, m := range ms {
		orders[i] = *m.ToDomain()
	}
	return orders, total, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status entity.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).Update("status", string(status)).Error
}

func (r *orderRepository) UpdateDeliveryStatus(ctx context.Context, id uint, status entity.DeliveryStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).Update("delivery_status", string(status)).Error
}

func (r *orderRepository) AttachPaymentID(ctx context.Context, id uint, paymentID string) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).
		Where("id = ?", id).Update("payment_id", paymentID).Error
}
