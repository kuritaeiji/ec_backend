package usecase

import (
	"context"
	"net/http"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/middleware"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/uptrace/bun"
)

type (
	SessionAccountUsecase struct {
		sessionAccountRepository repository.SessionAccountRepository
		accountRepository        repository.AccountRepository
		domainEventPublisher     share.DomainEventPublisher
		db                       bun.IDB
	}
)

var errEmailOrPasswordIsInvalid = share.CreateOriginalError(share.ErrorCodeOther, []string{"メールアドレスまたはパスワードが間違っています"})

func NewSessionAccountUsecase(
	sessionAccountRepository repository.SessionAccountRepository,
	accountRepository repository.AccountRepository,
	domainEventPublisher share.DomainEventPublisher,
	db bun.IDB,
) SessionAccountUsecase {
	return SessionAccountUsecase{
		sessionAccountRepository: sessionAccountRepository,
		accountRepository:        accountRepository,
		domainEventPublisher:     domainEventPublisher,
		db:                       db,
	}
}

// メールアドレス・パスワードでログインする
// セッションアカウントクッキーを返却する
func (sau SessionAccountUsecase) LoginByEmailAndPassword(ctx context.Context, email string, password string) (http.Cookie, error) {
	var sessionAccountCookie http.Cookie
	err := sau.db.RunInTx(ctx, nil, func(ctxt context.Context, tx bun.Tx) error {
		// メールアドレスをもとにアカウント集約を取得する。アカウント集約を取得できない場合はエラーメッセージを返却する。アカウントの認証タイプがメールアドレスではない場合はエラーメッセージを返却する
		account, ok, err := sau.accountRepository.FindByEmail(tx, ctxt, email)
		if err != nil {
			return err
		}

		// アカウント集約が存在しない場合
		if !ok {
			return errEmailOrPasswordIsInvalid
		}

		// 認証タイプがメールアドレスではない場合
		if account.AuthType != enum.AuthTypeEmail {
			return errEmailOrPasswordIsInvalid
		}

		// パスワードとパスワードダイジェストを比較する。パスワードが一致しない場合はエラーメッセージを返却する
		if !util.BcryptUtils.MatchPassword(*account.PasswordDigest, password) {
			return errEmailOrPasswordIsInvalid
		}

		// セッションアカウントを作成する
		sessionCart, existsSessionCart := middleware.SessionCartFromContext(ctx)
		var sessionAccount entity.SessionAccount
		sessionAccountCookie, sessionAccount = entity.CreateSessionAccount(account, sessionCart, existsSessionCart, tx, ctxt)
		return sau.sessionAccountRepository.Insert(ctxt, &sessionAccount, entity.SessionAccountExpiration, sau.domainEventPublisher)
	})

	return sessionAccountCookie, err
}

// ログアウトする
func (sau SessionAccountUsecase) Logout(ctx context.Context) error {
	sessionAccount, _ := middleware.SessionAccountFromContext(ctx)
	return sau.sessionAccountRepository.Delete(ctx, sessionAccount)
}
