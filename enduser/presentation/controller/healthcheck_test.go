package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckController(t *testing.T) {
	// given（前提条件）
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	con := NewHealthcheckController()

	// when（操作）
	err := con.Healthcheck(c)

	// then（期待する結果）
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
