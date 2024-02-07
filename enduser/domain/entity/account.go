package entity

import (
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/share"
)

// アカウント集約
type Account struct {
	ID                string        `json:"id"`
	Email             string        `json:"email"`
	PasswordDigest    *string       `json:"passwordDigest"`
	AuthType          enum.AuthType `json:"authType"`
	ExternalAccountID *string       `json:"externalAccountID"`
	IsActive          bool          `json:"isActive"`
	StripeCustomerId  *string       `json:"stripeCustmerID"`
	ReviewNickname    string        `json:"reviewNickname"`

	Events []share.DomainEvent
}
