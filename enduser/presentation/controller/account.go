package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/labstack/echo/v4"
)

type AccountController struct {
	accountUsecase usecase.AccountUsecase
}

func NewAccountController(accountUsecase usecase.AccountUsecase) AccountController {
	return AccountController{
		accountUsecase: accountUsecase,
	}
}

// メールアドレスによる新規アカウント登録時のフォーム
type AccountCreationForm struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// メールアドレスによって新規アカウントを登録する
func (ac AccountController) CreateAccountByEmail(c echo.Context) error {
	form := new(AccountCreationForm)
	err := c.Bind(form)
	if err != nil {
		return err
	}

	err = ac.accountUsecase.CreateAccountByEmail(c.Request().Context(), form.Email, form.Password, form.PasswordConfirmation)
	if err != nil {
		if oe, ok := err.(share.OriginalError); ok {
			return c.JSON(http.StatusOK, share.OriginalErrorToResult(oe))
		}

		return err
	}

	return c.JSON(http.StatusOK, share.SuccessResult())
}

// 新規アカウント登録時のメールアドレスを認証する
func (ac AccountController) AuthenticateEmail(c echo.Context) error {
	// クエリパラメータからJWTを取得する
	tokenString := c.QueryParam("token")

	// セッションカートCookieを取り出す
	cookie, ok, err := util.CookieUtils.GetCookie(c, entity.SessionCartCookieName)
	if err != nil {
		return err
	}

	// ユースケース層にメールアドレス認証の処理を委譲する
	accountSessionCookie, err := ac.accountUsecase.AuthenticateEmail(c.Request().Context(), tokenString, cookie.Value, ok)
	if err != nil {
		if originalErr, ok := err.(share.OriginalError); ok {
			return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?message=%s", os.Getenv("FRONT_URL"), originalErr.Messages[0]))
		}

		// TODO エラー画面にリダイレクトさせる
		return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s/error", os.Getenv("FRONT_URL")))
	}

	// セッションアカウントのセッションIDをCookieとしてセットする
	c.SetCookie(&accountSessionCookie)
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?message=%s", os.Getenv("FRONT_URL"), "メールアドレスを認証し、ログインしました"))
}
