package subscriber

import (
	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/adapter"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/share"
)

// アカウント有効化時にStripeのカスタマーを作成し、アカウント集約のストライプ顧客IDを更新する
type CreateStripeCustomerSubscriber struct {
	stripeAdapter        adapter.StripeAdapter
	accountRepository    repository.AccountRepository
}

func NewCreateStripeCustomerSubscriber(
	stripeAdapter adapter.StripeAdapter,
	accountRepository repository.AccountRepository,
) CreateStripeCustomerSubscriber {
	return CreateStripeCustomerSubscriber{
		stripeAdapter:        stripeAdapter,
		accountRepository:    accountRepository,
	}
}

// アカウント有効化イベントを購読する
func (subscriber CreateStripeCustomerSubscriber) TargetEvents() []share.DomainEvent {
	return []share.DomainEvent{entity.AccountActivatedEvent{}}
}

// アカウント有効化イベント発行時に、StripeのCustomerを作成する
func (subscriber CreateStripeCustomerSubscriber) Subscribe(event share.DomainEvent) error {
	accountActivatedEvent := event.(entity.AccountActivatedEvent)

	// StripeのCustomerを作成する
	stripeCustomerID, err := subscriber.stripeAdapter.CreateCustomer()
	if err != nil {
		return errors.WithStack(err)
	}

	// アカウント集約にストライプ顧客IDを設定する
	account := accountActivatedEvent.Account
	account.SetStripeCustomerID(stripeCustomerID)
	return subscriber.accountRepository.Update(accountActivatedEvent.DB, accountActivatedEvent.Ctx, account, nil)
}
