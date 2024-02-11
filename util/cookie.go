package util

import (
	"net/http"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
)

type cookieUtils struct{}

var CookieUtils = cookieUtils{}

func (cu cookieUtils) CreateCookie(name string, value string, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Domain:   os.Getenv("COOKIE_DOMAIN"),
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
}

func (cu cookieUtils) GetCookie(c echo.Context, name string) (http.Cookie, bool, error) {
	// Cookieを取り出す
	cookie, err := c.Cookie(name)

	if err != nil {
		// Cookieが存在しない場合
		if errors.Is(err, http.ErrNoCookie) {
			return http.Cookie{}, false, nil
		}
		// Cookieが存在しない以外のエラーの場合
		return http.Cookie{}, false, errors.WithStack(err)
	}

	// Cookieが存在する場合
	return *cookie, true, nil
}
