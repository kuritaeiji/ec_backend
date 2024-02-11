package entity

import (
	"github.com/kuritaeiji/ec_backend/util"
)

type (
	//カート集約
	Cart struct {
		ID        string
		AccountID string
		Version   int

		CartProducts []CartProduct
	}

	//カート商品
	CartProduct struct {
		ID        string
		CartID    string
		ProductID string
		Count     int
	}
)

// カート集約を作成する
func CreateCart(accountID string) Cart {
	return Cart{
		ID:           util.IDutils.GenerateID(),
		AccountID:    accountID,
		Version:      1,
		CartProducts: []CartProduct{},
	}
}

// セッションカートカート内の商品をカート集約に移動する
// セッションカート内の商品が販売中かつ在庫が存在することをチェックするためにセッションカートの商品集約リストを引数に取る
func (cart *Cart) MoveSessionCartProductsToCart(sessionCart SessionCart, products []Product) {
	for _, sessionCartProduct := range sessionCart.SessionCartProducts {
		product, ok := findProduct(products, sessionCartProduct.ProductID)
		//商品が存在しない場合はcontinueする
		if !ok {
			continue
		}

		//カート集約内にセッションカートの商品と一致する商品が存在する場合取得する
		cartProduct, ok := cart.findCartProduct(sessionCartProduct.ProductID)

		// 商品が販売中の場合のみセッションカートに商品を追加する
		if product.isOnSale() {
			// 在庫が追加したい個数以上の場合は商品を追加したい個数分セッションカートに追加する
			if product.StockCount >= sessionCartProduct.Count {
				// 既に同じ商品がカート内に存在する場合はカート内の商品の個数にセッションカートの商品の個数分だけ追加する
				if ok {
					cartProduct.Count += sessionCartProduct.Count
				} else {
					// 同じ商品がカート内に存在しない場合はカート内の商品をセッションカートの商品の個数分だけ追加する
					cart.CartProducts = append(cart.CartProducts, CartProduct{
						ID:        util.IDutils.GenerateID(),
						CartID:    cart.ID,
						ProductID: sessionCartProduct.ProductID,
						Count:     sessionCartProduct.Count,
					})
				}
				continue
			}

			//在庫が追加したい個数より少ないが1個以上の在庫を持つ場合は在庫分だけセッションカートに追加する
			if product.StockCount >= 1 {
				// 既に同じ商品がカート内に存在する場合はカート内の商品の個数に在庫数だけ追加する
				if ok {
					cartProduct.Count += product.StockCount
				} else {
					// 同じ商品がカート内に存在しない場合はカート内の商品を在庫数分だけ追加する
					cart.CartProducts = append(cart.CartProducts, CartProduct{
						ID:        util.IDutils.GenerateID(),
						CartID:    cart.ID,
						ProductID: sessionCartProduct.ProductID,
						Count:     product.StockCount,
					})
				}
				continue
			}

			//在庫が存在しない場合は何もしない
		}
	}
}

// 引数productIDに一致するカート内の商品を返却する
func (cart Cart) findCartProduct(productID string) (CartProduct, bool) {
	for _, cartProduct := range cart.CartProducts {
		if cartProduct.ID == productID {
			return cartProduct, true
		}
	}

	return CartProduct{}, false
}

// 引数productsから引数productIDに一致する商品を返却する
func findProduct(products []Product, productID string) (Product, bool) {
	for _, product := range products {
		if product.ID == productID {
			return product, true
		}
	}

	return Product{}, false
}
