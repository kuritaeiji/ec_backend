package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

// ミドルウェアを適用し、*echo.Echoとログインが必須の*echo.Groupを返却する
func SetupMiddleware(e *echo.Echo, container *dig.Container) (*echo.Echo, *echo.Group) {
	var loginG *echo.Group

	container.Invoke(func (sessionMiddleware SessionMiddleware, requireLoginMiddleware RequireLoginMiddleware)  {
		e.Use(middleware.Recover())
		e.Use(sessionMiddleware.Middleware)
		loginG = e.Group("", requireLoginMiddleware.Middleware)
	})

	return e, loginG
}