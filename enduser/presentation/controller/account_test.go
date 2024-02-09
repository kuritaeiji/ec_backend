package controller_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/kuritaeiji/ec_backend/enduser/registory"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/test"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"go.uber.org/dig"
)

type accountControllerTestSuite struct {
	suite.Suite
	container *dig.Container
}

func TestAccountController(t *testing.T) {
	err := config.SetupEnv()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("環境変数設定時にエラーが発生しました。\n%+v", err))
	}
	container, err := registory.NewTestContainer()
	if err != nil {
		assert.Fail(t, fmt.Sprintf("コンテナ作成時にエラーが発生しました。\n%+v", err))
	}
	suite.Run(t, &accountControllerTestSuite{
		container: container,
	})
}

func (suite *accountControllerTestSuite) tearDown() {
	cErr := suite.container.Invoke(func(db bun.IDB) {
		_, err := db.NewTruncateTable().Model(new(persistance.Account)).Exec(context.Background())
		if err != nil {
			suite.Fail(fmt.Sprintf("テーブルデータ削除時にエラー発生f\n+%v", err))
		}
	})
	if cErr != nil {
		suite.FailNow(cErr.Error())
	}
}

// アカウント作成（バリデーション）
func (suite *accountControllerTestSuite) TestCreateAccountValidation() {
	// given（前提条件）
	tests := []struct {
		name   string
		form   controller.AccountCreationForm
		status int
		result *share.Result
	}{
		{"メールアドレスとパスワードが空文字列の場合バリデーションエラー", controller.AccountCreationForm{"", "", ""}, http.StatusOK, &share.Result{Code: share.ResultCodeValidation, Messages: []string{fmt.Sprintf(util.RequriedMsg, "メールアドレス"), fmt.Sprintf(util.RequriedMsg, "パスワード")}}},
		{"メールアドレスの形式が正しくなくパスワードに不正な文字が含まれる場合バリデーションエラー", controller.AccountCreationForm{"testemail", "パスワードパスワード", "パスワードパスワード"}, http.StatusOK, &share.Result{Code: share.ResultCodeValidation, Messages: []string{fmt.Sprintf(util.EmailMsg, "メールアドレス"), fmt.Sprintf(util.AvailableSymbolPaswordMsg, "パスワード", `!"#$%&'()`)}}},
		{"パスワードとパスワード確認用が一致しない場合バリデーションエラー", controller.AccountCreationForm{"test@example.com", "password", "wrongpassword"}, http.StatusOK, &share.Result{Code: share.ResultCodeValidation, Messages: []string{fmt.Sprintf(util.PasswordConfirmationMsg, "パスワード")}}},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/account", test.FormToReader(tt.form))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			var err error

			// when（操作）
			cErr := suite.container.Invoke(func(con controller.AccountController) {
				err = con.CreateAccountByEmail(c)
			})

			// then（期待する結果）
			assert.Nil(t, cErr)
			assert.Nil(t, err)
			assert.Equal(t, tt.status, rec.Code)
			res := new(share.Result)
			test.ReaderToResponse(rec.Body, res)
			suite.Equal(tt.result, res)
		})
	}
}

// アカウント作成（既にDBに同じメールアドレスのアカウントが存在する場合）
func (suite *accountControllerTestSuite) TestCreateAccountEmailIsUnique() {
	// given（前提条件）
	email := "test@test.com"
	tests := []struct {
		name    string
		form    controller.AccountCreationForm
		setupFn func(t *testing.T, db bun.IDB)
		status  int
		result  *share.Result
	}{
		{
			name: "同一メールアドレスかつ未認証のアカウントが既に存在する場合、認証を促すエラーメッセージを返却する",
			form: controller.AccountCreationForm{email, "password", "password"},
			setupFn: func(t *testing.T, db bun.IDB) {
				_, err := db.NewInsert().Model(&persistance.Account{
					ID:                "test",
					Email:             email,
					PasswordDigest:    test.ToPointer("test"),
					AuthType:          int(enum.AuthTypeEmail),
					ExternalAccountID: nil,
					IsActive:          false,
					StripeCustomerId:  nil,
					ReviewNickname:    "test",
				}).Exec(context.Background())
				if err != nil {
					assert.FailNow(t, err.Error())
				}
			},
			status: http.StatusOK,
			result: &share.Result{Code: share.ResultCodeValidation, Messages: []string{"既に使用されているメールアドレスです。認証メールを確認してください"}},
		},
		{
			"同一メールアドレスかつ認証済みのアカウントが既に存在する場合、ログインを促すエラーメッセージを返却する",
			controller.AccountCreationForm{email, "password", "password"},
			func(t *testing.T, db bun.IDB) {
				_, err := db.NewInsert().Model(&persistance.Account{
					ID:                "test",
					Email:             email,
					PasswordDigest:    test.ToPointer("test"),
					AuthType:          int(enum.AuthTypeEmail),
					ExternalAccountID: nil,
					IsActive:          true,
					StripeCustomerId:  nil,
					ReviewNickname:    "test",
				}).Exec(context.Background())
				if err != nil {
					assert.FailNow(t, err.Error())
				}
			},
			http.StatusOK,
			&share.Result{Code: share.ResultCodeValidation, Messages: []string{"既に使用されているメールアドレスです。ログインしてください"}},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// when（操作）
			defer suite.tearDown()

			cErr := suite.container.Invoke(func(con controller.AccountController, db bun.IDB) {
				tt.setupFn(t, db)

				req := httptest.NewRequest(http.MethodPost, "/account", test.FormToReader(tt.form))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				e := echo.New()
				c := e.NewContext(req, rec)
				var err error

				err = con.CreateAccountByEmail(c)

				// then（期待する結果）
				assert.Nil(t, err)
				assert.Equal(t, tt.status, rec.Code)
				res := new(share.Result)
				test.ReaderToResponse(rec.Body, res)
				suite.Equal(tt.result, res)
			})
			assert.Nil(t, cErr)
		})
	}
}

// アカウント作成に成功する
func (suite *accountControllerTestSuite) TestCreateAccountSuccess() {
	// given（前提条件）
	defer suite.tearDown()

	// when（操作）
	cErr := suite.container.Invoke(func(con controller.AccountController, db bun.IDB) {
		email := "test@test.com"
		password := "password"
		req := httptest.NewRequest(http.MethodPost, "/account", test.FormToReader(controller.AccountCreationForm{
			Email:                email,
			Password:             password,
			PasswordConfirmation: password,
		}))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)
		err := con.CreateAccountByEmail(c)

		// then（期待する結果）
		suite.Nil(err)
		suite.Equal(http.StatusOK, rec.Code)
		res := new(share.Result)
		test.ReaderToResponse(rec.Body, res)
		suite.Equal(share.SuccessResult(), *res)

		count, err := db.NewSelect().Model(new(persistance.Account)).Count(context.Background())
		if err != nil {
			suite.FailNow(err.Error())
		}
		suite.Equal(1, count)

		account := new(persistance.Account)
		err = db.NewSelect().Model(account).Where("email = ?", email).Scan(context.Background())
		if err != nil {
			suite.FailNow(err.Error())
		}
		suite.Equal(email, account.Email)
		suite.Equal(int(enum.AuthTypeEmail), account.AuthType)
		suite.Nil(account.ExternalAccountID)
		suite.Equal(false, account.IsActive)
		suite.Nil(account.StripeCustomerId)
		suite.Equal("匿名", account.ReviewNickname)
		suite.True(util.BcryptUtils.MatchPassword(*account.PasswordDigest, password))
	})
	suite.Nil(cErr)
}
