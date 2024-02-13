package repository

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/uptrace/bun"
)

type CartRepository interface {
	FindByAccountID(db bun.IDB, ctx context.Context, accountID string) (entity.Cart, bool, error)
	Insert(db bun.IDB, ctx context.Context, cart entity.Cart) error
	Update(db bun.IDB, ctx context.Context, cart entity.Cart) error
}
