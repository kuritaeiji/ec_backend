package registory

import (
	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/enduser/domain/event"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/domain/service"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/uptrace/bun"
	"go.uber.org/dig"
)

// DIコンテナを作成する
func NewContainer() (*dig.Container, error) {
	container := dig.New()

	err := container.Provide(config.NewDB, dig.As(new(bun.IDB)))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = AddControllerTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUsecaseTo(container)
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

	err = AddRepositoryTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUtilsTo(container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// テスト用Iコンテナを作成する
func NewTestContainer() (*dig.Container, error) {
	container := dig.New()

	err := container.Provide(config.NewDB, dig.As(new(bun.IDB)))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = AddControllerTo(container)
	if err != nil {
		return nil, err
	}

	err = AddUsecaseTo(container)
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

	err = AddRepositoryTo(container)
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

// ドメインイベントパブリッシャーをDIコンテナに追加する
// 各イベントに対するサブスクライバーを定義する
func AddDomainEventPublisherTo(container *dig.Container) error {
	err := container.Provide(func() share.DomainEventPublisher {
		publisher := share.NewDomainEventPublisher()
		event.SubscribeAccountDomainEvent(publisher)
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
