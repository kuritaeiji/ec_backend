package registory

import (
	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/application/usecase"
	"github.com/kuritaeiji/ec_backend/enduser/domain/service"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/enduser/presentation/controller"
	"github.com/kuritaeiji/ec_backend/share"
	"github.com/kuritaeiji/ec_backend/util"
	"go.uber.org/dig"
)

// DIコンテナを作成する
func NewContainer() (*dig.Container, error) {
	container := dig.New()

	err := addControllerTo(container)
	if err != nil {
		return nil, err
	}

	err = addUsecaseTo(container)
	if err != nil {
		return nil, err
	}

	err = addDomainServiceTo(container)
	if err != nil {
		return nil, err
	}

	err = addRepositoryTo(container)
	if err != nil {
		return nil, err
	}

	err = addUtilsTo(container)
	if err != nil {
		return nil, err
	}

	err = addShareTo(container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// コントローラーをDIコンテナに追加する
func addControllerTo(container *dig.Container) error {
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

func addUsecaseTo(container *dig.Container) error {
	err := container.Provide(usecase.NewAccountUsecase)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ドメインサービスをDIコンテナに追加する
func addDomainServiceTo(container *dig.Container) error {
	err := container.Provide(service.NewAccountService)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}

// リポジトリをDIコンテナに追加する
func addRepositoryTo(container *dig.Container) error {
	err := container.Provide(persistance.NewAccountRepository)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ユーティリティーをDIコンテナに追加する
func addUtilsTo(container *dig.Container) error {
	err := container.Provide(util.NewLogger)
	if err != nil {
		return errors.WithStack(err)
	}

	err = container.Provide(util.NewValidationUtils)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// 共有クラスをDIコンテナに追加する
func addShareTo(container *dig.Container) error {
	err := container.Provide(share.NewDomainEventPublisher)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
