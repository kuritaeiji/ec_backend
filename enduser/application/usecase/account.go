package usecase

import (
	"context"
	"net/http"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/domain/service"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/uptrace/bun"
)

type AccountUsecase struct {
	accountDomainService     service.AccountDomainService
	accountRepository        repository.AccountRepository
	sessionAccountRepository repository.SessionAccountRepository
	domainEventPublisher     share.DomainEventPublisher
	db                       bun.IDB
}

func NewAccountUsecase(
	accountDomainService service.AccountDomainService,
	accountRepository repository.AccountRepository,
	sessionAccountRepository repository.SessionAccountRepository,
	domainEventPublisher share.DomainEventPublisher,
	db bun.IDB,
) AccountUsecase {
	return AccountUsecase{
		accountDomainService: accountDomainService,
		accountRepository:    accountRepository,
		sessionAccountRepository: sessionAccountRepository,
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

// 新規アカウント登録時のメールアドレスを認証する
func (au AccountUsecase) AuthenticateEmail(ctx context.Context, tokenString string, sessionCartSessionID string, existsSessionCartSessionID bool) (http.Cookie, error) {
	var accountSessionCookie http.Cookie
	err := au.db.RunInTx(ctx, nil, func(ctxt context.Context, tx bun.Tx) error {
		// アカウントのメールアドレス認証とアカウントの有効化を行う
		account, err := au.accountDomainService.AuthenticateEmail(tx, ctxt, tokenString)
		if err != nil {
			return err
		}

		err = au.accountRepository.Update(tx, ctxt, account, au.domainEventPublisher)
		if err != nil {
			return err
		}

		// セッションアカウントを作成する
		var sessionAccount entity.SessionAccount
		accountSessionCookie, sessionAccount = entity.CreateSessionAccount(account, sessionCartSessionID, existsSessionCartSessionID, tx, ctxt)
		err = au.sessionAccountRepository.Insert(ctxt, sessionAccount, entity.SessionAccountExpiration, au.domainEventPublisher)
		return err
	})

	return accountSessionCookie, err
}
