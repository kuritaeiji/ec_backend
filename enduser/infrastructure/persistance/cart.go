package persistance

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/uptrace/bun"
)

type (
	//カートテーブル
	Cart struct {
		bun.BaseModel `bun:"table:carts"`

		ID           string        `bun:",pk"`
		AccountID    string        `bun:",notnull,unique"`
		Version      int           `bun:",notnull"`
		CartProducts []CartProduct `bun:"rel:has-many,join:id=cart_id"`
	}

	//カート商品テーブル
	CartProduct struct {
		bun.BaseModel `bun:"table:cart_products"`

		ID        string `bun:",pk"`
		CartID    string `bun:",notnull"`
		ProductID string `bun:",notnull"`
		Count     int    `bun:",notnull"`
	}

	//カートリポジトリの実装
	cartRepository struct{}
)

func NewCartRepository() cartRepository {
	return cartRepository{}
}

func (cr cartRepository) FindByAccountID(db bun.IDB, ctx context.Context, accountID string) (entity.Cart, bool, error) {
	var cart Cart
	err := db.NewSelect().Model(&cart).Where("account_id = ?", accountID).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Cart{}, false, nil
		}

		return entity.Cart{}, false, errors.WithStack(err)
	}

	return cr.toEntity(cart), true, nil
}

// カート集約を登録する
func (cr cartRepository) Insert(db bun.IDB, ctx context.Context, cart entity.Cart) error {
	mCart := cr.toModel(cart)

	//カートを登録する
	_, err := db.NewInsert().Model(&mCart).Exec(ctx)

	return errors.WithStack(err)
}

// カート集約を更新する
func (cr cartRepository) Update(db bun.IDB, ctx context.Context, cart entity.Cart) error {
	mCart := cr.toModel(cart)

	//カート内の商品をすべて削除する
	_, err := db.NewDelete().Model(new(CartProduct)).Where("cart_id = ?", mCart.ID).Exec(ctx)
	if err != nil {
		return err
	}

	//すべてのカート商品を登録する
	_, err = db.NewInsert().Model(&mCart.CartProducts).Exec(ctx)
	if err != nil {
		return err
	}

	//カートを更新する（楽観ロックする）
	mCart.Version = mCart.Version + 1
	res, err := db.NewUpdate().Model(&mCart).WherePK().Where("version = ?", mCart.Version-1).Exec(ctx)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		return ErrOptimisticLocking
	}

	return nil
}

func (cr cartRepository) toModel(cart entity.Cart) Cart {
	cartProducts := make([]CartProduct, 0, len(cart.CartProducts))
	for _, p := range cart.CartProducts {
		cartProducts = append(cartProducts, CartProduct{
			ID:        p.ID,
			CartID:    p.CartID,
			ProductID: p.ProductID,
			Count:     p.Count,
		})
	}

	return Cart{
		ID:           cart.ID,
		AccountID:    cart.AccountID,
		Version:      cart.Version,
		CartProducts: cartProducts,
	}
}

func (cr cartRepository) toEntity(cart Cart) entity.Cart {
	cartProducts := make([]entity.CartProduct, 0, len(cart.CartProducts))
	for _, cartProduct := range cart.CartProducts {
		cartProducts = append(cartProducts, entity.CartProduct{
			ID:        cartProduct.ID,
			CartID:    cartProduct.CartID,
			ProductID: cartProduct.ProductID,
			Count:     cartProduct.Count,
		})
	}

	return entity.Cart{
		ID:           cart.ID,
		AccountID:    cart.AccountID,
		Version:      cart.Version,
		CartProducts: cartProducts,
	}
}
