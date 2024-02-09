package usecase

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/domain/service"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/uptrace/bun"
)

type AccountUsecase struct {
	accountDomainService service.AccountDomainService
	accountRepository    repository.AccountRepository
	domainEventPublisher share.DomainEventPublisher
	db                   bun.IDB
}

func NewAccountUsecase(
	accountDomainService service.AccountDomainService,
	accountRepository repository.AccountRepository,
	domainEventPublisher share.DomainEventPublisher,
	db bun.IDB,
) AccountUsecase {
	return AccountUsecase{
		accountDomainService: accountDomainService,
		accountRepository:    accountRepository,
		domainEventPublisher: domainEventPublisher,
		db:                   db,
	}
}

// メールアドレスによって新規アカウントを登録する
func (au AccountUsecase) CreateAccountByEmail(ctx context.Context, email string, password string, passwordConfirmation string) error {
	err := au.db.RunInTx(ctx, nil, func(ctxt context.Context, tx bun.Tx) error {
		account, err := au.accountDomainService.CreateAccountByEmail(email, password, passwordConfirmation, tx, ctxt)
		if err != nil {
			return err
		}

		err = au.accountRepository.Insert(tx, ctxt, account, au.domainEventPublisher)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
