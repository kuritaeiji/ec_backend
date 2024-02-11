package entity

import (
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
)

type (
	// 商品集約
	Product struct {
		ID             string
		CategoryID     string
		CategoryName   string
		Name           string
		ReviewScore    int
		Price          int
		SalePrice      int
		Description    string
		Status         enum.ProductStatus
		StockCount     int
		Version        int
		CreateDateTime time.Time

		ProductImages []ProductImage
	}

	// 商品画像
	ProductImage struct {
		ID        string
		ProductID string
		Order     int
		Path      string
		Image     string // TODO 型を正しくする必要あり
	}
)

// 商品が販売中の場合trueを、そうでない場合falseを返却する
func (product Product) isOnSale() bool {
	return product.Status == enum.OnSale
}
