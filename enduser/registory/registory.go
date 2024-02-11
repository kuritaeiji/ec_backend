package registory

import (
	"log"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/enduser/domain/adapter"
	"github.com/kuritaeiji/ec_backend/enduser/domain/adapter/mocks"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/domain/service"
	"github.com/kuritaeiji/ec_backend/enduser/domain/subscriber"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/bridge"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun"
	"go.uber.org/dig"
)

// DIコンテナを作成する
func NewContainer() (*dig.Container, error) {
	container := dig.New()

	err := AddDBTo(container)
	if err != nil {
		return nil, err
	}

	err = AddRepositoryTo(container)
	if err != nil {
		return nil, err
	}

	err = AddAdapterTo(container)
	if err != nil {
		return nil, err
	}

	err = AddDomainServiceTo(container)
	if err != nil {
		return nil, err
	}

	err = AddDomainEventPublisherTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUsecaseTo(container)
	if err != nil {
		return nil, err
	}

	err = AddControllerTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUtilsTo(container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// テスト用DIコンテナを作成する
func NewTestContainer(t *testing.T) (*dig.Container, error) {
	container := dig.New()

	err := AddTestingTo(container, t)
	if err != nil {
		return nil, err
	}

	err = AddDBTo(container)
	if err != nil {
		return nil, err
	}

	err = AddRepositoryTo(container)
	if err != nil {
		return nil, err
	}

	err = AddMockAdapterTo(container)
	if err != nil {
		return nil, err
	}

	err = AddDomainServiceTo(container)
	if err != nil {
		return nil, err
	}

	err = AddDomainEventPublisherTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUsecaseTo(container)
	if err != nil {
		return nil, err
	}

	err = AddControllerTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUtilsTo(container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// DBをDIコンテナに追加する
func AddDBTo(container *dig.Container) error {
	err := container.Provide(config.NewDB, dig.As(new(bun.IDB)))
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(config.NewRedisClient)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// コントローラーをDIコンテナに追加する
func AddControllerTo(container *dig.Container) error {
	err := container.Provide(controller.NewHealthcheckController)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(controller.NewAccountController)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ユースケースをDIコンテナに追加する
func AddUsecaseTo(container *dig.Container) error {
	err := container.Provide(usecase.NewAccountUsecase)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ドメインサービスをDIコンテナに追加する
func AddDomainServiceTo(container *dig.Container) error {
	err := container.Provide(service.NewAccountService)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ドメインイベントサブスクライバーをDIコンテナに追加する
// ドメインイベントパブリッシャーをDIコンテナに追加する
// 各イベントに対するサブスクライバーを定義する
func AddDomainEventPublisherTo(container *dig.Container) error {
	err := container.Provide(subscriber.NewAccountCreatedByEmailSubscriber)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(subscriber.NewAccountActivatedEventSubscriber)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(subscriber.NewMoveSessionCartProductToCartSubscriber)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(func() share.DomainEventPublisher {
		publisher := share.NewDomainEventPublisher()
		err := container.Invoke(func(
			sendAuthenticationEmailSubscriber subscriber.SendAuthenticationEmailSubscriber,
			createCartSubscriber subscriber.CreateCartSubscriber,
			moveSessionCartProductToCartSubscriber subscriber.MoveSessionCartProductToCartSubscriber,
		) {
			// どのイベントをサブスクライブするかを設定する
			publisher.Subscribe(sendAuthenticationEmailSubscriber.TargetEvents(), sendAuthenticationEmailSubscriber)
			publisher.Subscribe(createCartSubscriber.TargetEvents(), createCartSubscriber)
			publisher.Subscribe(moveSessionCartProductToCartSubscriber.TargetEvents(), moveSessionCartProductToCartSubscriber)
		})
		if err != nil {
			log.Fatal(errors.WithStack(err))
		}
		return publisher
	})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// リポジトリをDIコンテナに追加する
func AddRepositoryTo(container *dig.Container) error {
	err := container.Provide(persistance.NewAccountRepository, dig.As(new(repository.AccountRepository)))
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(persistance.NewCartRepository, dig.As(new(repository.CartRepository)))
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(persistance.NewProductRepository, dig.As(new(repository.ProductRepository)))
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(persistance.NewSessionAccountRepository, dig.As(new(repository.SessionAccountRepository)))
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(persistance.NewSessionCartRepository, dig.As(new(repository.SessionCartRepository)))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// アダプターをDIコンテナに追加する
func AddAdapterTo(container *dig.Container) error {
	err := container.Provide(bridge.NewEmailAdapter, dig.As(new(adapter.EmailAdapter)))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// モック化されたアダプターをDIコンテナに追加する
func AddMockAdapterTo(container *dig.Container) error {
	err := container.Provide(mocks.NewEmailAdapter, dig.As((new(adapter.EmailAdapter))))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ユーティリティーをDIコンテナに追加する
func AddUtilsTo(container *dig.Container) error {
	err := container.Provide(util.NewLogger)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(util.NewValidationUtils, dig.As(new(util.ValidationUtils)))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// モックを作成する際に*testing.Tが必要なのでDIコンテナに追加する
func AddTestingTo(container *dig.Container, t *testing.T) error {
	err := container.Provide(func() *testing.T {
		return t
	}, dig.As(new(interface {
		mock.TestingT
		Cleanup(func())
	})))

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
