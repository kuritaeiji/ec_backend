package handler

import (
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func setupAccountHandler(e *echo.Echo, container *dig.Container) error {
	err := container.Invoke(func(ac controller.AccountController) {
		e.POST("/account", ac.CreateAccountByEmail)
	})
	return err
}
