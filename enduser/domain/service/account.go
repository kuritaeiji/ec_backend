package service

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/enduser/domain/event"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/domain/validator"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type AccountDomainService struct {
	accountRepository repository.AccountRepository
	validationUtils   util.ValidationUtils
	logger            echo.Logger
}

func NewAccountService(accountRepository repository.AccountRepository, validationUtils util.ValidationUtils, logger echo.Logger) AccountDomainService {
	return AccountDomainService{
		accountRepository: accountRepository,
		validationUtils:   validationUtils,
		logger:            logger,
	}
}

const initialReviewNickname = "匿名"

func (as AccountDomainService) CreateAccountByEmail(email *string, password *string, passwordConfirmation *string, db bun.IDB, ctx context.Context) (entity.Account, error) {
	// バリデーション実施
	validationAccount := validator.ValidationAccountForCreation{
		Email:    email,
		Password: password,
		PasswordConfirmation: passwordConfirmation,
		AuthType: enum.AuthTypeEmail,
	}

	err := as.validationUtils.Struct(validationAccount)
	if err != nil {
		return entity.Account{}, as.validationUtils.CreateValidationMessages(err)
	}

	// メールアドレスが一意であることを確認
	account, isUnique, err := as.emailIsUnique(*email, db, ctx)
	if err != nil {
		return entity.Account{}, err
	}
	// DBに同一メールアドレスが登録されており、認証済みの場合、ログインを促すエラーメッセージを返却する
	if !isUnique && account.IsActive {
		return entity.Account{}, share.OriginalError{Code: share.ErrorCodeValidation, Messages: []string{"既に使用されているメールアドレスです。ログインしてください"}}
	}
	// DBに同一メールアドレスが登録されており、未認証の場合、認証を促すエラーメッセージを返却する
	if !isUnique && !account.IsActive {
		return entity.Account{}, share.OriginalError{Code: share.ErrorCodeValidation, Messages: []string{"既に使用されているメールアドレスです。認証メールを確認してください"}}
	}

	// パスワードダイジェストを作成する
	passwordDigest, err := util.BcryptUtils.GeneratePasswordDigest(*password)
	if err != nil {
		as.logger.Error("Bcryptによるパスワードのハッシュ化失敗\n", err)
		return entity.Account{}, err
	}

	// メールアドレスによるアカウント登録イベントを作成する
	event := event.AccountCreatedByEmailEvent{Email: *email}

	return entity.Account{
		ID:                util.IDutils.GenerateID(),
		Email:             *email,
		PasswordDigest:    &passwordDigest,
		AuthType:          enum.AuthTypeEmail,
		ExternalAccountID: nil,
		IsActive:          false,
		StripeCustomerId:  nil,
		ReviewNickname:    initialReviewNickname,
		Events: []share.DomainEvent{event},
	}, nil
}

// 同一メールアドレスのアカウントが存在してない場合はtrueを返却し、そうでない場合はfalseを返却する
func (as AccountDomainService) emailIsUnique(email string, db bun.IDB, ctx context.Context) (entity.Account, bool, error) {
	account, ok, err := as.accountRepository.FindByEmail(db, ctx, email)
	if err != nil {
		return account, false, err
	}

	return account, !ok, nil
}
