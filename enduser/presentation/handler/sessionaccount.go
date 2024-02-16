package handler

import (
	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func setupSessionAccountHandler(e *echo.Echo, loginG *echo.Group, container *dig.Container) error {
	err := container.Invoke(func (sessionAccountController controller.SessionAccountController)  {
		e.GET("/login", sessionAccountController.LoginByEmailAndPassword)
		loginG.DELETE("/logout", sessionAccountController.Logout)
	})
	return errors.WithStack(err)
}
