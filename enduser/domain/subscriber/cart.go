package subscriber

import (
	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/share"
)

type (
	// カートを新規作成するサブスクライバー
	// アカウント有効化イベント発行時に実行される
	CreateCartSubscriber struct {
		cartRepository repository.CartRepository
	}

	// セッションカートからDBカートに商品を移動させるサブスクライバー
	// セッションアカウント作成時（ログイン時）に実行される
	MoveSessionCartProductToCartSubscriber struct {
		cartRepository        repository.CartRepository
		productRepository     repository.ProductRepository
		sessionCartRepository repository.SessionCartRepository
	}
)

func NewAccountActivatedEventSubscriber(cartRepository repository.CartRepository) CreateCartSubscriber {
	return CreateCartSubscriber{
		cartRepository: cartRepository,
	}
}

// アカウント有効化イベントを購読する
func (subscriber CreateCartSubscriber) TargetEvents() []share.DomainEvent {
	return []share.DomainEvent{entity.AccountActivatedEvent{}}
}

// カートを新規作成する
func (subscriber CreateCartSubscriber) Subscribe(event share.DomainEvent) error {
	accountActivatedEvent := event.(entity.AccountActivatedEvent)
	cart := entity.CreateCart(accountActivatedEvent.Account.ID)
	err := subscriber.cartRepository.Insert(accountActivatedEvent.DB, accountActivatedEvent.Ctx, cart)
	return err
}

func NewMoveSessionCartProductToCartSubscriber(
	cartRepository repository.CartRepository,
	productRepository repository.ProductRepository,
	sessionCartRepository repository.SessionCartRepository,
) MoveSessionCartProductToCartSubscriber {
	return MoveSessionCartProductToCartSubscriber{
		cartRepository: cartRepository,
	}
}

// セッションアカウント作成イベント（ログインイベント）を購読する
func (subscriber MoveSessionCartProductToCartSubscriber) TargetEvents() []share.DomainEvent {
	return []share.DomainEvent{entity.SessionAccountCreatedEvent{}}
}

// セッションカート内の商品をDBカートに移動させる
func (subscriber MoveSessionCartProductToCartSubscriber) Subscribe(event share.DomainEvent) error {
	sessionAccountCreatedEvent := event.(entity.SessionAccountCreatedEvent)

	// セッションカートのセッションIDがCookieとして存在しない場合、returnする
	if !sessionAccountCreatedEvent.ExistsSessionCartSessionID {
		return nil
	}

	// セッションカートを取得する
	sessionCart, ok, err := subscriber.sessionCartRepository.FindBySessionID(sessionAccountCreatedEvent.Ctx, sessionAccountCreatedEvent.SessionCartSessionID)
	if err != nil {
		return err
	}
	if !ok {
		// セッションカートが存在しない場合はreturnする
		return nil
	}

	// セッションカート内に商品が存在しない場合はreturnする
	if len(sessionCart.SessionCartProducts) > 0 {
		return nil
	}

	// アカウントに紐づくDBカートを取得する
	cart, ok, err := subscriber.cartRepository.FindByAccountID(sessionAccountCreatedEvent.DB, sessionAccountCreatedEvent.Ctx, sessionAccountCreatedEvent.AccountID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.WithStack(errors.New("カートが見つかりません"))
	}

	// セッションカート内の商品の商品集約リストを取得する
	products, err := subscriber.productRepository.FindByIDs(sessionAccountCreatedEvent.DB, sessionAccountCreatedEvent.Ctx, sessionCart.ProductIDs(), false)
	if err != nil {
		return err
	}

	// セッションカートからDBカートに商品を移動させる
	cart.MoveSessionCartProductsToCart(sessionCart, products)
	err = subscriber.cartRepository.Update(sessionAccountCreatedEvent.DB, sessionAccountCreatedEvent.Ctx, cart)
	if err != nil {
		return err
	}

	// セッションカートを削除する
	return subscriber.sessionCartRepository.Delete(sessionAccountCreatedEvent.Ctx, sessionCart)
}
