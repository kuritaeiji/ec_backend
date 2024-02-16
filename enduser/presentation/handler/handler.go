package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

// ハンドラーのセットアップ
func SetupHandlers(e *echo.Echo, loginG *echo.Group, container *dig.Container) error {
	err := setupHealthcheckHandler(e, container)
	if err != nil {
		return err
	}

	err = setupAccountHandler(e, container)
	if err != nil {
		return err
	}

	err = setupSessionAccountHandler(e, loginG, container)
	if err != nil {
		return err
	}

	return nil
}
