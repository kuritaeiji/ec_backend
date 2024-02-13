package entity_test

import (
	"testing"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/stretchr/testify/assert"
)

// 商品ステータス          在庫とカート内の商品数               既にカートに商品が存在する
// 販売中・停止中・中止     在庫>=商品数・1<=在庫<商品数 在庫=0  カートに同一商品が存在する・存在しない
func TestMoveSessionCartToCart(t *testing.T) {
	// given（前提条件）
	cartID := "1"
	productID := "2"
	productID2 := "3"

	type params struct {
		cart        entity.Cart
		sessionCart entity.SessionCart
		products    []entity.Product
	}

	tests := []struct {
		Name         string
		Params       params
		ExpectedCart entity.Cart
	}{
		{
			Name: "セッションカートに商品が存在しない場合、何もしない",
			Params: params{
				cart:        entity.Cart{ID: cartID},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{}},
				products:    []entity.Product{},
			},
			ExpectedCart: entity.Cart{ID: cartID},
		},
		{
			Name: "セッションカート内の商品IDと一致する商品集約が存在しない場合、何もしない",
			Params: params{
				cart:        entity.Cart{ID: cartID},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{{ProductID: productID, Count: 1}}},
				products:    []entity.Product{{ID: "2", Status: enum.OnSale, StockCount: 10}},
			},
			ExpectedCart: entity.Cart{ID: cartID},
		},
		{
			Name: "商品ステータスが販売中かつ在庫>=セッションカート商品個数かつ既にカートに同一商品が存在する場合、カート内の商品の個数をセッションカートの商品個数だけ増やす",
			Params: params{
				cart:        entity.Cart{ID: cartID, CartProducts: []entity.CartProduct{{ID: "1", CartID: cartID, ProductID: productID, Count: 1}}},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{{ProductID: productID, Count: 1}, {ProductID: productID2, Count: 1}}},
				products:    []entity.Product{{ID: productID, Status: enum.OnSale, StockCount: 1}, {ID: productID2, Status: enum.OnSale, StockCount: 2}},
			},
			ExpectedCart: entity.Cart{ID: cartID, CartProducts: []entity.CartProduct{{CartID: cartID, ProductID: productID, Count: 2}, {ID: "1", CartID: cartID, ProductID: productID2, Count: 1}}},
		},
		{
			Name: "商品ステータスが販売中かつ1<=在庫数<セッションカートの商品個数かつカートに同一商品が存在しない場合、カート内に在庫数分の商品を追加する",
			Params: params{
				cart:        entity.Cart{ID: cartID},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{{ProductID: productID, Count: 2}, {ProductID: productID2, Count: 3}}},
				products:    []entity.Product{{ID: productID, Status: enum.OnSale, StockCount: 1}, {ID: productID2, Status: enum.OnSale, StockCount: 2}},
			},
			ExpectedCart: entity.Cart{ID: cartID, CartProducts: []entity.CartProduct{{CartID: cartID, ProductID: productID, Count: 1}, {CartID: cartID, ProductID: productID2, Count: 2}}},
		},
		{
			Name: "商品ステータスが販売中かつ在庫が0かつカートに同一商品が存在しない場合何もしない",
			Params: params{
				cart:        entity.Cart{ID: cartID},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{{ProductID: productID, Count: 1}}},
				products:    []entity.Product{{ID: productID, Status: enum.OnSale, StockCount: 0}},
			},
			ExpectedCart: entity.Cart{ID: cartID, CartProducts: []entity.CartProduct{}},
		},
		{
			Name: "商品ステータスが販売中以外の場合、何もしない",
			Params: params{
				cart:        entity.Cart{ID: cartID},
				sessionCart: entity.SessionCart{SessionCartProducts: []entity.SessionCartProduct{{ProductID: "1", Count: 1}, {ProductID: productID2, Count: 1}}},
				products:    []entity.Product{{ID: "1", Status: enum.SalesSuspend, StockCount: 2}, {ID: productID2, Status: enum.SalesEnded, StockCount: 2}},
			},
			ExpectedCart: entity.Cart{ID: cartID, CartProducts: []entity.CartProduct{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// when（操作）
			tt.Params.cart.MoveSessionCartProductsToCart(tt.Params.sessionCart, tt.Params.products)

			// then（期待する結果）
			assert.Equal(t, len(tt.ExpectedCart.CartProducts), len(tt.ExpectedCart.CartProducts))
			for i, expectedCartProduct := range tt.ExpectedCart.CartProducts {
				for j, acturalCartProduct := range tt.Params.cart.CartProducts {
					if i != j {
						continue
					}

					assert.Equal(t, expectedCartProduct.CartID, acturalCartProduct.CartID, "カートID")
					assert.Equal(t, expectedCartProduct.ProductID, acturalCartProduct.ProductID, "商品ID")
					assert.Equal(t, expectedCartProduct.Count, acturalCartProduct.Count, "個数")
				}
			}
		})
	}
}
