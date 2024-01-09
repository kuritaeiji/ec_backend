package handler

import (
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/labstack/echo/v4"
)

func SetupHealthcheckHandler(e *echo.Echo) {
	healthcheckController :=  controller.NewHealthcheckController()

	e.GET("/healthcheck", healthcheckController.Healthcheck)
}
