package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthcheckController struct{}

func NewHealthcheckController() HealthcheckController {
	return HealthcheckController{}
}

func (hc HealthcheckController) Healthcheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
