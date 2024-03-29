package persistance

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
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

func NewAccountRepository() accountRepository {
	return accountRepository{}
}

func (ar accountRepository) FindByEmail(db bun.IDB, ctx context.Context, email string) (entity.Account, bool, error) {
	account := Account{}
	err := db.NewSelect().Model(&account).Where("email = ?", email).Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return entity.Account{}, false, nil
	}

	return ar.toEntity(account), true, errors.WithStack(err)
}

func (ar accountRepository) Insert(db bun.IDB, ctx context.Context, account *entity.Account, domainEventPublisher share.DomainEventPublisher) error {
	mAccount := ar.toModel(*account)
	_, err := db.NewInsert().Model(&mAccount).Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if domainEventPublisher == nil {
		return nil
	}

	// アカウント集約内のイベントを発行する
	events := account.ClearEvents()
	err = domainEventPublisher.Publish(events)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ar accountRepository) Update(db bun.IDB, ctx context.Context, account *entity.Account, domainEventPublisher share.DomainEventPublisher) error {
	mAccout := ar.toModel(*account)
	_, err := db.NewUpdate().Model(&mAccout).WherePK().Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if domainEventPublisher == nil {
		return nil
	}

	events := account.ClearEvents()
	err = domainEventPublisher.Publish(events)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ar accountRepository) toEntity(account Account) entity.Account {
	return entity.Account{
		ID:                account.ID,
		Email:             account.Email,
		PasswordDigest:    account.PasswordDigest,
		AuthType:          enum.AuthType(account.AuthType),
		ExternalAccountID: account.ExternalAccountID,
		IsActive:          account.IsActive,
		StripeCustomerID:  account.StripeCustomerId,
		ReviewNickname:    account.ReviewNickname,
	}
}

func (ar accountRepository) toModel(account entity.Account) Account {
	return Account{
		ID:                account.ID,
		Email:             account.Email,
		PasswordDigest:    account.PasswordDigest,
		AuthType:          int(account.AuthType),
		ExternalAccountID: account.ExternalAccountID,
		IsActive:          account.IsActive,
		StripeCustomerId:  account.StripeCustomerID,
		ReviewNickname:    account.ReviewNickname,
	}
}
