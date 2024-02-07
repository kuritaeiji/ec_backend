package handler

import (
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func setupHealthcheckHandler(e *echo.Echo, container *dig.Container) error {
	err := container.Invoke(func(hc controller.HealthcheckController) {
		e.GET("/healthcheck", hc.Healthcheck)
	})
	return err
}
