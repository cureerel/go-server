package repository

import (
	"context"

	"github.com/cureerel/cserver/internal/domain/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	GetByID(ctx context.Context, id string) (*entity.Session, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)
	GetActiveByUserID(ctx context.Context, userID uint) ([]*entity.Session, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUserID(ctx context.Context, userID uint) error
	UpdateLastActive(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}
