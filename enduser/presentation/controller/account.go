package controller

import (
	"net/http"

	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/share"
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

type AccountCreationForm struct {
	Email *string `json:"email"`
	Password *string `json:"password"`
	PasswordConfirmation *string `json:"passwordConfirmation"`
}

//メールアドレスによって新規アカウントを登録する
func (ac AccountController) CreateAccountByEmail(c echo.Context) error {
	form := new(AccountCreationForm)
	err := c.Bind(form)
	if err != nil {
		return err
	}

	err = ac.accountUsecase.CreateAccountByEmail(c.Request().Context(), form.Email, form.Password, form.PasswordConfirmation)
	if err != nil {
		if oe, ok := err.(share.OriginalError); ok {
			return c.JSON(http.StatusOK, oe)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}