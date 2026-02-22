package repository

import (
    "context"
    "time"
    "github.com/cureerel/gotemplate/internal/domain/entity"
)

type AuthRepository interface {
    SaveRefreshToken(ctx context.Context, token *entity.RefreshToken) error
    GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
    RevokeRefreshToken(ctx context.Context, tokenHash string) error
    BlacklistToken(ctx context.Context, tokenHash string, expiresAt time.Time) error
    IsTokenBlacklisted(ctx context.Context, tokenHash string) (bool, error)
}