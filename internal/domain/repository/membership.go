package repository

import (
	"context"
	"github.com/cureerel/cserver/internal/domain/entity"
)

type MembershipRepository interface {
	Create(ctx context.Context, membership *entity.Membership) error
	GetByUserID(ctx context.Context, userID uint) (*entity.Membership, error)
	Update(ctx context.Context, membership *entity.Membership) error
	Cancel(ctx context.Context, id uint) error
}
