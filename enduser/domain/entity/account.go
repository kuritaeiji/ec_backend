package entity

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/uptrace/bun"
)

// アカウント集約
type (
	Account struct {
		ID                string        `json:"id"`
		Email             string        `json:"email"`
		PasswordDigest    *string       `json:"passwordDigest"`
		AuthType          enum.AuthType `json:"authType"`
		ExternalAccountID *string       `json:"externalAccountID"`
		IsActive          bool          `json:"isActive"`
		StripeCustomerID  *string       `json:"stripeCustmerID"`
		ReviewNickname    string        `json:"reviewNickname"`

		Events []share.DomainEvent
	}

	// メールアドレスによるアカウント登録イベント
	AccountCreatedByEmailEvent struct {
		Email string
	}
	// アカウント有効化イベント
	AccountActivatedEvent struct {
		Account *Account
		DB      bun.IDB
		Ctx     context.Context
	}
)

const (
	accountCreatedByEmailEventName share.DomainEventName = "AccountCreatedByEmailEvent"
	accountActivatedEventName      share.DomainEventName = "AccountActivatedEvent"
)

func (ae AccountCreatedByEmailEvent) Name() share.DomainEventName {
	return accountCreatedByEmailEventName
}

func (ae AccountActivatedEvent) Name() share.DomainEventName {
	return accountActivatedEventName
}

func (account *Account) SetStripeCustomerID(stripeCustomerID string) {
	account.StripeCustomerID = &stripeCustomerID
}

// ドメインイベント配列を削除し、ドメインイベント配列を返却する
func (account *Account) ClearEvents() []share.DomainEvent {
	events := account.Events
	account.Events = []share.DomainEvent{}
	return events
}
