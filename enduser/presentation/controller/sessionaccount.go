package controller

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/labstack/echo/v4"
)

type (
	SessionAccountController struct {
		sessionAccountUsecase usecase.SessionAccountUsecase
	}

	LoginByEmailAndPasswordForm struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

func NewSessionAccountControler(sessionAccountUsecase usecase.SessionAccountUsecase) SessionAccountController {
	return SessionAccountController{
		sessionAccountUsecase: sessionAccountUsecase,
	}
}

// メールアドレス・パスワードでログインする
func (sac SessionAccountController) LoginByEmailAndPassword(c echo.Context) error {
	// リクエストボディーを取得する
	var form LoginByEmailAndPasswordForm
	err := c.Bind(&form)
	if err != nil {
		return errors.WithStack(err)
	}

	// ログインする
	sessionAccountCookie, err := sac.sessionAccountUsecase.LoginByEmailAndPassword(c.Request().Context(), form.Email, form.Password)
	if err != nil {
		if originalErr, ok := err.(share.OriginalError); ok {
			return c.JSON(http.StatusOK, originalErr)
		}

		return err
	}

	// セッションアカウントのクッキーを作成する
	c.SetCookie(&sessionAccountCookie)
	return c.JSON(http.StatusOK, share.SuccessResult())
}

// ログアウトする
func (sac SessionAccountController) Logout(c echo.Context) error {
	err := sac.sessionAccountUsecase.Logout(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, share.SuccessResult())
}
