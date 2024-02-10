package event

import (
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/adapter"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/bridge"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
)

func SubscribeAccountDomainEvent(publisher share.DomainEventPublisher, subscriber AccountCreatedByEmailSubscriber) {
	publisher.Subscribe(AccountCreatedByEmailEvent{}, subscriber)
}

type (
	AccountCreatedByEmailEvent struct {
		Email string
	}
	AccountCreatedByEmailSubscriber struct {
		emailAdapter adapter.EmailAdapter
	}
)

const (
	accountCreatedByEmailEventName share.DomainEventName = "AccountCreatedByEmailEvent"
)

func NewAccountCreatedByEmailSubscriber(emailAdapter adapter.EmailAdapter) AccountCreatedByEmailSubscriber {
	return AccountCreatedByEmailSubscriber{
		emailAdapter: emailAdapter,
	}
}

func (ae AccountCreatedByEmailEvent) Name() share.DomainEventName {
	return accountCreatedByEmailEventName
}

// メールアドレスによるアカウント登録イベントが発行されたときに、認証メールを送信する
func (as AccountCreatedByEmailSubscriber) Subscribe(event share.DomainEvent) error {
	accountEvent := event.(AccountCreatedByEmailEvent)
	email := accountEvent.Email

	jwtString, err := util.JwtUtils.CreateJwt(email, 24*time.Hour)
	if err != nil {
		return errors.WithStack(err)
	}

	text := fmt.Sprintf(`<a href="%s?token=%s">メールアドレスを認証する</a><br/>有効期限は24時間`, os.Getenv("BACKEND_URL"), jwtString)

	err = as.emailAdapter.SendEmail(bridge.From, email, "認証メール", text)
	if err != nil {
		return err
	}

	return nil
}
