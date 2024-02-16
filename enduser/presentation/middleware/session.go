package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/labstack/echo/v4"
)

type SessionMiddleware struct {
	sessionAccountRepository repository.SessionAccountRepository
	sessionCartRepository repository.SessionCartRepository
	timeUtils util.TimeUtils
	logger echo.Logger
}

func NewSessionMiddleware(
	sessionAccountRepository repository.SessionAccountRepository,
	sessionCartRepository repository.SessionCartRepository,
	timeUtils util.TimeUtils,
	logger echo.Logger,
) SessionMiddleware {
	return SessionMiddleware{
		sessionAccountRepository: sessionAccountRepository,
		sessionCartRepository: sessionCartRepository,
		timeUtils: timeUtils,
		logger: logger,
	}
}

// セッションアカウントとセッションカートを取得する
// セッションアカウントの存在の有無とセッションカートの存在の有無も取得する
// セッションアカウントの有効期限が1週間より小さい場合有効期限を2週間に伸ばす
// セッションカートの有効期限が1週間より小さい場合有効期限を30日に伸ばす
func (m SessionMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// セッションアカウント
		sessionAccount, sessionAccountCookie, existsSessionAccount, err := m.getSessionAccount(c, ctx)
		if err != nil {
			return err
		}

		// セッションアカウントが存在し、セッションアカウントの有効期限が1週間未満の場合、有効期限を2週間に伸ばす
		if existsSessionAccount && sessionAccountCookie.Expires.Before(m.timeUtils.NowJP().Add(7 * 24 * time.Hour)) {
			sessionAccountCookie.Expires = m.timeUtils.NowJP().Add(entity.SessionAccountExpiration)
			err := m.sessionAccountRepository.UpdateExpiration(ctx, sessionAccount, entity.SessionAccountExpiration)
			if err == nil {
				c.SetCookie(&sessionAccountCookie)
			} else {
				// セッションアカウントの有効期限を伸ばす際にエラーが発生してもエラーを返却しない
				// 次にAPIを呼び出す際に有効期限を伸ばせばよいから
				m.logger.Errorf("%+v", err)
			}
		}

		if existsSessionAccount {
			// セッションアカウントをContextに登録する
			ctx = m.contextWithSessionAccount(ctx, sessionAccount)
		}

		// セッションカート
		sessionCart, sessionCartCookie, existsSessionCart, err := m.getSessionCart(c, ctx)
		if err != nil {
			return err
		}

		if existsSessionCart {
			// セッションカートをContextに登録する
			ctx = m.contextWithSessionCart(ctx, sessionCart)
		}

		// セッションカートが存在し、セッションカートの有効期限が2週間未満の場合、有効期限を30日に伸ばす
		if existsSessionCart && sessionCartCookie.Expires.Before(m.timeUtils.NowJP().Add(14 * 24 * time.Hour)) {
			sessionCartCookie.Expires = m.timeUtils.NowJP().Add(entity.SessionCartExpiration)
			err := m.sessionCartRepository.UpdateExpiration(ctx, sessionCart, entity.SessionCartExpiration)
			if err == nil {
				c.SetCookie(&sessionCartCookie)
			} else {
				// セッションカートの有効期限を伸ばす際にエラーが発生してもエラーを返却しない
				// 次にAPIを呼び出す際に有効期限を伸ばせばよいから
				m.logger.Errorf("%+v", err)
			}
		}

		// echo.Contextのhttp.Requestに新しいcontext.Contextをセットする
		c.SetRequest(c.Request().WithContext(ctx))

		// 次のミドルウェア実行
		return next(c)
	}
}

// セッションアカウント・セッションアカウントクッキーを取得する
func (m SessionMiddleware) getSessionAccount(c echo.Context, ctx context.Context) (entity.SessionAccount, http.Cookie, bool, error) {
	// セッションアカウントクッキーを取得する
	sessionAccountCookie, ok, err := util.CookieUtils.GetCookie(c, entity.SessionAccountCookieName)
	if err != nil {
		return entity.SessionAccount{}, http.Cookie{}, false, err
	}

	if !ok {
		return entity.SessionAccount{}, http.Cookie{}, false, nil
	}

	// セッションアカウントを取得する
	sessionAccount, ok, err := m.sessionAccountRepository.FindBySessionID(ctx, sessionAccountCookie.Value)
	if err != nil {
		return entity.SessionAccount{}, http.Cookie{}, false, err
	}

	if !ok {
		return entity.SessionAccount{}, http.Cookie{}, false, nil
	}

	return sessionAccount, sessionAccountCookie, true, nil
}

// セッションカート・セッションカートクッキーを取得する
func (m SessionMiddleware) getSessionCart(c echo.Context, ctx context.Context) (entity.SessionCart, http.Cookie, bool, error) {
	sessionCartCookie, ok, err := util.CookieUtils.GetCookie(c, entity.SessionCartCookieName)
	if err != nil {
		return entity.SessionCart{}, http.Cookie{}, false, err
	}

	if !ok {
		return entity.SessionCart{}, http.Cookie{}, false, nil
	}

	sessionCart, ok, err := m.sessionCartRepository.FindBySessionID(ctx, sessionCartCookie.Value)
	if err != nil {
		return entity.SessionCart{}, http.Cookie{}, false, err
	}

	if !ok {
		return entity.SessionCart{}, http.Cookie{}, false, nil
	}

	return sessionCart, sessionCartCookie, true, nil
}

const (
	sessionAccountCtxKey = "SessionAccountCtx"
	sessionCartCtxKey = "SessionCartCtx"
)

// セッションアカウントをContextに登録する
func (m SessionMiddleware) contextWithSessionAccount(ctx context.Context, sessionAccount entity.SessionAccount) context.Context {
	return context.WithValue(ctx, sessionAccountCtxKey, sessionAccount)
}

// セッションカートをContextに登録する
func (m SessionMiddleware) contextWithSessionCart(ctx context.Context, sessionCart entity.SessionCart) context.Context {
	return context.WithValue(ctx, sessionCartCtxKey, sessionCart)
}

// セッションアカウントをContextから取り出す
func SessionAccountFromContext(ctx context.Context) (entity.SessionAccount, bool) {
	sessionAccount, ok := ctx.Value(sessionAccountCtxKey).(entity.SessionAccount)
	return sessionAccount, ok
}

// セッションカートをContextから取り出す
func SessionCartFromContext(ctx context.Context) (entity.SessionCart, bool) {
	sessionCart, ok := ctx.Value(sessionCartCtxKey).(entity.SessionCart)
	return sessionCart, ok
}
