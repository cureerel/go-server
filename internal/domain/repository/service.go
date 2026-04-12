// internal/domain/repository/service.go
package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type ServiceFilter struct {
	Page    int
	Limit   int
	Status  string
	OwnerID *uint
	Search  string
}

type ServiceRepository interface {
	Create(ctx context.Context, svc *entity.Service) error
	GetByID(ctx context.Context, id uint) (*entity.Service, error)
	GetAll(ctx context.Context, filter ServiceFilter) ([]entity.Service, int64, error)
	Update(ctx context.Context, svc *entity.Service) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status string) error
	IncrementViews(ctx context.Context, id uint) error
	RecordView(ctx context.Context, serviceID uint, visitorHash string) (bool, error)
}
