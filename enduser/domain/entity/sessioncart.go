package entity

import "time"

type (
	//セッションカート集約
	SessionCart struct {
		SessionID           string
		SessionCartProducts []SessionCartProduct
	}

	//セッションカート商品
	SessionCartProduct struct {
		ProductID string
		Count     int
	}
)

const (
	SessionCartExpiration = 30 * 24 * time.Hour // セッションカートの有効期限は30日
	SessionCartCookieName = "SessionCartSessionID"
)

// セッションカート内の商品の商品ID配列を返却する
func (sessionCart SessionCart) ProductIDs() []string {
	ids := make([]string, len(sessionCart.SessionCartProducts))
	for _, sessionCartProduct := range sessionCart.SessionCartProducts {
		ids = append(ids, sessionCartProduct.ProductID)
	}

	return ids
}
