package subscriber

import (
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/adapter"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/bridge"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
)

type (
	SendAuthenticationEmailSubscriber struct {
		emailAdapter adapter.EmailAdapter
	}
)

func NewAccountCreatedByEmailSubscriber(emailAdapter adapter.EmailAdapter) SendAuthenticationEmailSubscriber {
	return SendAuthenticationEmailSubscriber{
		emailAdapter: emailAdapter,
	}
}

// メールアドレスによるアカウント登録イベントを購読する
func (as SendAuthenticationEmailSubscriber) TargetEvents() []share.DomainEvent {
	return []share.DomainEvent{entity.AccountCreatedByEmailEvent{}}
}

// メールアドレスによるアカウント登録イベントが発行されたときに、認証メールを送信する
func (as SendAuthenticationEmailSubscriber) Subscribe(event share.DomainEvent) error {
	accountEvent := event.(entity.AccountCreatedByEmailEvent)
	email := accountEvent.Email

	jwtString, err := util.JwtUtils.CreateJwt(email, 24*time.Hour)
	if err != nil {
		return errors.WithStack(err)
	}

	text := fmt.Sprintf(`<a href="%s/account/email/auth?token=%s">メールアドレスを認証する</a><br/>有効期限は24時間`, os.Getenv("BACKEND_URL"), jwtString)

	err = as.emailAdapter.SendEmail(bridge.From, email, "認証メール", text)
	if err != nil {
		return err
	}

	return nil
}
