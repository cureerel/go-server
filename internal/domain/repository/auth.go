package repository

import (
	"context"
	"time"
)

type AuthRepository interface {
	StoreRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (string, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

