package repository

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
)

type SessionCartRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (entity.SessionCart, error)
	Delete(ctx context.Context, sessionCart entity.SessionCart) error
}
