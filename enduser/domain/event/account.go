package event

import (
	"fmt"

	"github.com/kuritaeiji/ec_backend/share"
)

func SubscribeAccountDomainEvent(publisher share.DomainEventPublisher) {
	publisher.Subscribe(AccountCreatedByEmailEvent{}, accountCreatedByEmailSubscriber{})
}

type (
	AccountCreatedByEmailEvent struct {
		Email string
	}
	accountCreatedByEmailSubscriber struct{}
)

const (
	accountCreatedByEmailEventName share.DomainEventName = "AccountCreatedByEmailEvent"
)

func (ae AccountCreatedByEmailEvent) Name() share.DomainEventName {
	return accountCreatedByEmailEventName
}

func (as accountCreatedByEmailSubscriber) Subscribe(event share.DomainEvent) error {
	accountEvent := event.(AccountCreatedByEmailEvent)
	// TODO 認証メールを送信する
	fmt.Println(accountEvent)

	return nil
}
