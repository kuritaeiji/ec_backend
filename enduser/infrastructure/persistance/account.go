package persistance

import (
	"context"
	"database/sql"
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/uptrace/bun"
)

// アカウントテーブル
type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID                string `bun:",pk"`
	Email             string `bun:",notnull,unique"`
	PasswordDigest    *string
	AuthType          int
	ExternalAccountID *string
	IsActive          bool `bun:",notnull"`
	StripeCustomerId  *string
	ReviewNickname    string    `bun:",notnull"`
	DeleteDateTime    time.Time `bun:",soft_delete,nullzero"`
}

type accountRepository struct {
}

func NewAccountRepository() repository.AccountRepository {
	return accountRepository{}
}

func (ar accountRepository) FindByEmail(db bun.IDB, ctx context.Context, email string) (entity.Account, bool, error) {
	account := entity.Account{}
	err := db.NewSelect().Model(&account).Where("email = ?", email).Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return account, false, nil
	}

	return account, true, err
}

func (ar accountRepository) Insert(db bun.IDB, ctx context.Context, account entity.Account, domainEventPublisher share.DomainEventPublisher) error {
	_, err := db.NewInsert().Model(&account).Exec(ctx)
	if err != nil {
		return err
	}

	// アカウント集約内のイベントを発行する
	err = domainEventPublisher.Publish(account.Events)
	if err != nil {
		return err
	}

	return nil
}
