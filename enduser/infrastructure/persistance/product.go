package persistance

import (
	"context"
	"time"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/uptrace/bun"
)

type (
	// 商品テーブル
	Product struct {
		bun.BaseModel `bun:"table:products"`

		ID                   string    `bun:",pk"`
		CategoryID           string    `bun:",notnull"`
		Name                 string    `bun:",notnull"`
		Description          string    `bun:",type:text,notnull"`
		StockCount           int       `bun:",notnull"`
		Version              int       `bun:",notnull"`
		CreateDateTime       time.Time `bun:",notnull"`
		CreateStaffAccountID string    `bun:",notnull"`
		UpdateDateTime       *time.Time
		UpdateStaffAccountID *string

		ProductStatuses   []ProductStatus    `bun:"rel:has-many,join:id=product_id"`
		ProductImages     []ProductImage     `bun:"rel:has-many,join:id=product_id"`
		ProductPrices     []ProductPrice     `bun:"rel:has-many,join:id=product_id"`
		ProductSalePrices []ProductSalePrice `bun:"rel:has-many,join:id=product_id"`
		ReviewScores      []ReviewScore      `bun:"rel:has-many,join:id=product_id"`
		Category          Category           `bun:"rel:belongs-to,join:category_id=id"`
	}

	// 商品ステータステーブル
	ProductStatus struct {
		bun.BaseModel `bun:"table:product_statuses"`

		ID                 string    `bun:",pk"`
		ProductID          string    `bun:",notnull"`
		Status             int       `bun:",notnull"`
		EffectiveStartDate time.Time `bun:",notnull,type:date"`
		EffectiveEndDate   string    `bun:",notnull,type:date"`
	}

	// 商品画像テーブル
	ProductImage struct {
		bun.BaseModel `bun:"table:product_images"`

		ID        string `bun:",pk"`
		ProductID string `bun:",notnull"`
		Order     int    `bun:",notnull"`
		Path      string `bun:",notnull"`
	}

	// 商品価格テーブル
	ProductPrice struct {
		bun.BaseModel `bun:"table:product_prices"`

		ID                 string    `bun:",pk"`
		ProductID          string    `bun:",notnull"`
		TaxInclusivePrice  int       `bun:",notnull"`
		EffectiveStartDate time.Time `bun:",notnull,type:date"`
		EffectiveEndDate   time.Time `bun:",notnull,type:date"`
	}

	// 商品セール価格テーブル
	ProductSalePrice struct {
		bun.BaseModel `bun:"table:product_sale_prices"`

		ID                 string    `bun:",pk"`
		ProductID          string    `bun:",notnull"`
		TaxInclusivePrice  int       `bun:",notnull"`
		EffectiveStartDate time.Time `bun:",notnull,type:date"`
		EffectiveEndDate   time.Time `bun:",notnull,type:date"`
	}

	// レビュー点数テーブル
	ReviewScore struct {
		bun.BaseModel `bun:"table:review_scores"`

		ID        string    `bun:",pk"`
		ProductID string    `bun:",notnull"`
		Score     int       `bun:",notnull"`
		Date      time.Time `bun:",notnull,type:date"`
	}

	// カテゴリーテーブル
	Category struct {
		bun.BaseModel `bun:"table:categories"`

		ID   string `bun:",pk"`
		Name string `bun:",notnull"`
	}

	productRepository struct{}
)

func NewProductRepository() productRepository {
	return productRepository{}
}

// 商品ID配列に一致する商品配列を返却する。引数withImageがtrueの場合はS3から画像を取得し、そうでない場合は取得しない。
func (pr productRepository) FindByIDs(db bun.IDB, ctx context.Context, ids []string, withImage bool) ([]entity.Product, error) {
	today := time.Now()

	var products []Product
	err := db.NewSelect().Model(&products).
		// 商品IDがidsに含まれる商品
		Where("products.id in (?)", bun.In(ids)).
		// システム日付が商品ステータスの適用開始日以上かつ適用終了日以下
		Where("? between product_statuses.EffectiveStartDate and product_statuses.EffectiveEndDate", today).
		// システム日付が商品価格の適用開始日以上かつ適用終了日以下
		Where("? between product_prices.EffectiveStartDate and product_prices.EffectiveEndDate", today).
		// システム日付が商品セール価格の適用開始日以上かつ適用終了日以下
		Where("? between product_sale_prices.EffectiveStartDate and product_sale_prices.EffectiveEndDate", today).
		// レビュー点数の日付がシステム日付
		Where("review_scores.date = ?", today).
		Scan(ctx)

	if err != nil {
		return []entity.Product{}, err
	}

	eProducts := make([]entity.Product, len(products))
	for _, product := range products {
		eProducts = append(eProducts, pr.toEntity(product, withImage))
	}
	return eProducts, nil
}

func (pr productRepository) toEntity(product Product, withImage bool) entity.Product {
	images := make([]entity.ProductImage, len(product.ProductImages))
	for _, image := range product.ProductImages {
		images = append(images, entity.ProductImage{
			ID:        image.ID,
			ProductID: image.ProductID,
			Order:     image.Order,
			Path:      image.Path,
			Image:     "", // TODO 正しい値にする必要あり
		})
	}

	return entity.Product{
		ID:             product.ID,
		CategoryID:     product.CategoryID,
		CategoryName:   product.Category.Name,
		Name:           product.Name,
		ReviewScore:    product.ReviewScores[0].Score,
		Price:          product.ProductPrices[0].TaxInclusivePrice,
		SalePrice:      product.ProductSalePrices[0].TaxInclusivePrice,
		Description:    product.Description,
		Status:         enum.ProductStatus(product.ProductStatuses[0].Status),
		StockCount:     product.StockCount,
		Version:        product.Version,
		CreateDateTime: product.CreateDateTime,
		ProductImages:  images,
	}
}
