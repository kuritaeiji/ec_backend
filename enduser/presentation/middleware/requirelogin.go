package middleware

import (
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/labstack/echo/v4"
)

type RequireLoginMiddleware struct{}

func NewRequireLoginMiddleware() RequireLoginMiddleware {
	return RequireLoginMiddleware{}
}

// ログインしていない場合にエラーコードとエラーメッセージを返却する
func (m RequireLoginMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, ok := SessionAccountFromContext(c.Request().Context())
		if !ok {
			return share.CreateOriginalError(share.ErrorCodeNoLogin, []string{"ログインしてください"})
		}

		return next(c)
	}
}
