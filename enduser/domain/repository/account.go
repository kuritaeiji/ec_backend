package repository

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/uptrace/bun"
)

type AccountRepository interface {
	FindByEmail(db bun.IDB, ctx context.Context, email string) (entity.Account, bool, error)
	Insert(db bun.IDB, ctx context.Context, account *entity.Account, domainEventPublisher share.DomainEventPublisher) error
	Update(db bun.IDB, ctx context.Context, account *entity.Account, domainEventPublisher share.DomainEventPublisher) error
}
