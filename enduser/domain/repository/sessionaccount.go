package repository

import (
	"context"
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/share"
)

type SessionAccountRepository interface {
	Insert(ctx context.Context, sessionAccount entity.SessionAccount, expiration time.Duration, eventPublisher share.DomainEventPublisher) error
	FindBySessionID(ctx context.Context, sessionID string) (entity.SessionAccount, bool, error)
}
