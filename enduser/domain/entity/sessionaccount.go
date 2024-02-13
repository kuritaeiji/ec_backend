package entity

import (
	"context"
	"net/http"
	"time"

	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/uptrace/bun"
)

type (
	// セッションアカウント集約
	SessionAccount struct {
		AccountID string
		SessionID string

		Events []share.DomainEvent
	}

	// セッションアカウント作成イベント
	SessionAccountCreatedEvent struct {
		AccountID                  string
		SessionCartSessionID       string
		ExistsSessionCartSessionID bool
		DB                         bun.IDB
		Ctx                        context.Context
	}
)

const (
	sessionAccountCreatedEventName = "SessionAccountCreatedEvent"
	SessionAccountExpiration       = 14 * 24 * time.Hour // セッションアカウントの有効期限は2週間
	SessionAccountCookieName       = "AccountSessionID"
)

func (event SessionAccountCreatedEvent) Name() share.DomainEventName {
	return share.DomainEventName(sessionAccountCreatedEventName)
}

// セッションアカウントを作成する
func CreateSessionAccount(account Account, sessionCartSessionID string, existsSessionCartSessionID bool, db bun.IDB, ctx context.Context) (http.Cookie, SessionAccount) {
	// Cookieを作成する
	sessionID := util.IDutils.GenerateID()
	cookie := util.CookieUtils.CreateCookie(SessionAccountCookieName, sessionID, time.Now().Add(SessionAccountExpiration))

	// セッションアカウントを作成し返却する
	return cookie, SessionAccount{
		AccountID: account.ID,
		SessionID: sessionID,
		Events: []share.DomainEvent{
			SessionAccountCreatedEvent{
				AccountID:                  account.ID,
				SessionCartSessionID:       sessionCartSessionID,
				ExistsSessionCartSessionID: existsSessionCartSessionID,
				DB:                         db,
				Ctx:                        ctx,
			},
		},
	}
}

func (sessionAccount *SessionAccount) ClearEvents() []share.DomainEvent {
	events := sessionAccount.Events
	sessionAccount.Events = []share.DomainEvent{}
	return events
}
