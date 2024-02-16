package repository

import (
	"context"
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
)

type SessionCartRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (entity.SessionCart, bool, error)
	Delete(ctx context.Context, sessionCart entity.SessionCart) error
	UpdateExpiration(ctx context.Context, sessionCart entity.SessionCart, expiration time.Duration) error
}
