// internal/infrastructure/postgres/repositories/service_repository.go
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

type serviceRepository struct{ db *gorm.DB }

func NewServiceRepository(db *gorm.DB) repository.ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, svc *entity.Service) error {
	m := models.ServiceFromDomain(svc)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	svc.ID = m.ID
	return nil
}

func (r *serviceRepository) GetByID(ctx context.Context, id uint) (*entity.Service, error) {
	var m models.Service
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *serviceRepository) GetAll(ctx context.Context, filter repository.ServiceFilter) ([]entity.Service, int64, error) {
	var ms []models.Service
	var total int64
	offset := (filter.Page - 1) * filter.Limit

	q := r.db.WithContext(ctx).Model(&models.Service{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.OwnerID != nil {
		q = q.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.Search != "" {
		q = q.Where("title ILIKE ?", "%"+filter.Search+"%")
	}
	q = q.Order("created_at DESC")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	svcs := make([]entity.Service, len(ms))
	for i, m := range ms {
		svcs[i] = *m.ToDomain()
	}
	return svcs, total, nil
}

func (r *serviceRepository) Update(ctx context.Context, svc *entity.Service) error {
	return r.db.WithContext(ctx).Save(models.ServiceFromDomain(svc)).Error
}

func (r *serviceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Service{}, id).Error
}

func (r *serviceRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).
		Where("id = ?", id).Update("status", status).Error
}

func (r *serviceRepository) IncrementViews(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Service{}).
		Where("id = ?", id).
		UpdateColumn("views_total", gorm.Expr("views_total + 1")).Error
}

func (r *serviceRepository) RecordView(ctx context.Context, serviceID uint, visitorHash string) (bool, error) {
	today := time.Now().UTC().Format("2006-01-02")
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO service_views (service_id, visitor_hash, viewed_date, created_at)
		VALUES (?, ?, ?, NOW())
		ON CONFLICT (service_id, visitor_hash, viewed_date) DO NOTHING`,
		serviceID, visitorHash, today,
	)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected > 0 {
		_ = r.IncrementViews(ctx, serviceID)
		return true, nil
	}
	return false, nil
}
