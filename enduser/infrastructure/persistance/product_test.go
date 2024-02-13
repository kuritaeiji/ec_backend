package persistance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/kuritaeiji/ec_backend/config"
	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/kuritaeiji/ec_backend/enduser/domain/enum"
	"github.com/kuritaeiji/ec_backend/enduser/domain/repository"
	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/kuritaeiji/ec_backend/util"
	"github.com/kuritaeiji/ec_backend/util/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
)

type productRepositoryTestSuite struct {
	suite.Suite
	productRepository repository.ProductRepository
	db                bun.IDB
	timeUtilsMock     *mocks.TimeUtils
	timeUtils         util.TimeUtils
}

func TestProductRepository(t *testing.T) {
	err := config.SetupEnv()
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("環境変数設定時にエラーが発生しました。\n%+v", err))
	}
	timeUtils := util.NewTimeUtils()
	timeUtilsMock := mocks.NewTimeUtils(t, timeUtils)
	suite.Run(t, &productRepositoryTestSuite{
		productRepository: persistance.NewProductRepository(timeUtilsMock),
		db:                config.NewDB(),
		timeUtilsMock:     timeUtilsMock,
		timeUtils:         util.NewTimeUtils(),
	})
}

func (suite *productRepositoryTestSuite) tearDown() {
	tables := []any{
		new(persistance.Product),
		new(persistance.ProductStatus),
		new(persistance.ProductImage),
		new(persistance.ProductPrice),
		new(persistance.ProductSalePrice),
		new(persistance.ReviewScore),
		new(persistance.Category),
	}
	for _, table := range tables {
		_, err := suite.db.NewTruncateTable().Model(table).Exec(context.Background())
		if err != nil {
			assert.FailNow(suite.T(), fmt.Sprintf("テーブルデータ（%v）削除時に失敗", table))
		}
	}
}

func (suite *productRepositoryTestSuite) insertProducts(products []persistance.Product) {
	for _, p := range products {
		_, err := suite.db.NewInsert().Model(&p).Exec(context.Background())
		if err != nil {
			assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
		}

		for _, ps := range p.ProductStatuses {
			_, err := suite.db.NewInsert().Model(&ps).Exec(context.Background())
			if err != nil {
				assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
			}
		}

		for _, pi := range p.ProductImages {
			_, err := suite.db.NewInsert().Model(&pi).Exec(context.Background())
			if err != nil {
				assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
			}
		}

		for _, pp := range p.ProductPrices {
			_, err := suite.db.NewInsert().Model(&pp).Exec(context.Background())
			if err != nil {
				assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
			}
		}

		for _, psp := range p.ProductSalePrices {
			_, err := suite.db.NewInsert().Model(&psp).Exec(context.Background())
			if err != nil {
				assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
			}
		}

		for _, rs := range p.ReviewScores {
			_, err := suite.db.NewInsert().Model(&rs).Exec(context.Background())
			if err != nil {
				assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
			}
		}

		_, err = suite.db.NewInsert().Model(&p.Category).Exec(context.Background())
		if err != nil {
			assert.FailNow(suite.T(), fmt.Sprintf("商品作成時にエラー発生\n+%+v", errors.WithStack(err)))
		}
	}
}

func (suite *productRepositoryTestSuite) TestFindByIds() {
	// given（前提条件）
	type expected struct {
		Products []entity.Product
		Err      error
	}

	type params struct {
		IDs       []string
		WithImage bool
	}

	// 2000年1/10を現在日付にする
	date10 := suite.timeUtils.DateJP(2000, 1, 10)
	suite.timeUtilsMock.On("NowJP").Return(date10)

	// 2000年1/1・1/8・1/9・1/11・1/31を作成する
	date1 := suite.timeUtils.DateJP(2000, 1, 1)
	date8 := suite.timeUtils.DateJP(2000, 1, 8)
	date9 := suite.timeUtils.DateJP(2000, 1, 9)
	date11 := suite.timeUtils.DateJP(2000, 1, 11)
	date31 := suite.timeUtils.DateJP(2000, 1, 31)

	tests := []struct {
		Name     string
		Setup    func(t *testing.T)
		Params   params
		Expected expected
	}{
		{
			Name: "商品IDに一致する商品が存在しない場合、空配列が返却される",
			Setup: func(t *testing.T) {
				product := persistance.Product{
					ID:                   "100",
					CategoryID:           "100",
					Name:                 "100",
					Description:          "100",
					StockCount:           10,
					Version:              1,
					CreateDateTime:       suite.timeUtils.TimeToUTC(date1),
					CreateStaffAccountID: "100",
					UpdateDateTime:       nil,
					UpdateStaffAccountID: nil,

					ProductStatuses:   []persistance.ProductStatus{},
					ProductImages:     []persistance.ProductImage{},
					ProductPrices:     []persistance.ProductPrice{},
					ProductSalePrices: []persistance.ProductSalePrice{},
					ReviewScores:      []persistance.ReviewScore{},
					Category:          persistance.Category{},
				}
				suite.insertProducts([]persistance.Product{product})
			},
			Params:   params{IDs: []string{"1"}, WithImage: false},
			Expected: expected{Products: []entity.Product{}, Err: nil},
		},
		{
			Name: "商品IDに一致する商品が存在する場合、商品配列が返却される",
			Setup: func(t *testing.T) {
				productID := "1"
				categoryID := "1"
				product := persistance.Product{
					ID:                   productID,
					CategoryID:           categoryID,
					Name:                 "商品名",
					Description:          "商品説明",
					StockCount:           10,
					Version:              1,
					CreateDateTime:       suite.timeUtils.TimeToUTC(date1),
					CreateStaffAccountID: "100",
					UpdateDateTime:       nil,
					UpdateStaffAccountID: nil,

					ProductStatuses: []persistance.ProductStatus{
						{ID: "1", ProductID: productID, Status: int(enum.OnSale), EffectiveStartDate: date1, EffectiveEndDate: date9},
						{ID: "2", ProductID: productID, Status: int(enum.SalesSuspend), EffectiveStartDate: date10, EffectiveEndDate: date31},
					},
					ProductImages: []persistance.ProductImage{
						{ID: "1", ProductID: productID, Order: 1, Path: "/path1"},
						{ID: "2", ProductID: productID, Order: 2, Path: "/path2"},
					},
					ProductPrices: []persistance.ProductPrice{
						{ID: "1", ProductID: productID, TaxInclusivePrice: 100, EffectiveStartDate: date1, EffectiveEndDate: date10},
						{ID: "2", ProductID: productID, TaxInclusivePrice: 120, EffectiveStartDate: date11, EffectiveEndDate: date31},
					},
					ProductSalePrices: []persistance.ProductSalePrice{
						{ID: "1", ProductID: productID, TaxInclusivePrice: 90, EffectiveStartDate: date1, EffectiveEndDate: date8},
						{ID: "2", ProductID: productID, TaxInclusivePrice: 110, EffectiveStartDate: date9, EffectiveEndDate: date31},
					},
					ReviewScores: []persistance.ReviewScore{
						{ID: "1", ProductID: productID, Score: 1, Date: date9},
						{ID: "2", ProductID: productID, Score: 2, Date: date10},
						{ID: "3", ProductID: productID, Score: 3, Date: date11},
					},
					Category: persistance.Category{ID: categoryID, Name: "カテゴリー1"},
				}
				suite.insertProducts([]persistance.Product{product})
			},
			Params: params{IDs: []string{"1"}, WithImage: false},
			Expected: expected{Products: []entity.Product{{
				ID:             "1",
				CategoryID:     "1",
				CategoryName:   "カテゴリー1",
				Name:           "商品名",
				ReviewScore:    2,
				Price:          100,
				SalePrice:      110,
				Description:    "商品説明",
				Status:         enum.SalesSuspend,
				StockCount:     10,
				Version:        1,
				CreateDateTime: date1,

				ProductImages: []entity.ProductImage{
					{ID: "1", ProductID: "1", Order: 1, Path: "/path1", Image: ""},
					{ID: "2", ProductID: "1", Order: 2, Path: "/path2", Image: ""},
				},
			}}, Err: nil},
		},
		{
			Name: "レビュー点数が存在しない場合、商品のレビュー点数は0点になる",
			Setup: func(t *testing.T) {
				productID := "1"
				categoryID := "1"
				product := persistance.Product{
					ID:                   productID,
					CategoryID:           categoryID,
					Name:                 "商品名",
					Description:          "商品説明",
					StockCount:           10,
					Version:              1,
					CreateDateTime:       suite.timeUtilsMock.TimeToUTC(date1),
					CreateStaffAccountID: "100",
					UpdateDateTime:       nil,
					UpdateStaffAccountID: nil,

					ProductStatuses: []persistance.ProductStatus{
						{ID: "1", ProductID: productID, Status: int(enum.OnSale), EffectiveStartDate: date1, EffectiveEndDate: date9},
						{ID: "2", ProductID: productID, Status: int(enum.SalesSuspend), EffectiveStartDate: date10, EffectiveEndDate: date31},
					},
					ProductImages: []persistance.ProductImage{
						{ID: "1", ProductID: productID, Order: 1, Path: "/path1"},
						{ID: "2", ProductID: productID, Order: 2, Path: "/path2"},
					},
					ProductPrices: []persistance.ProductPrice{
						{ID: "1", ProductID: productID, TaxInclusivePrice: 100, EffectiveStartDate: date1, EffectiveEndDate: date10},
						{ID: "2", ProductID: productID, TaxInclusivePrice: 120, EffectiveStartDate: date11, EffectiveEndDate: date31},
					},
					ProductSalePrices: []persistance.ProductSalePrice{
						{ID: "1", ProductID: productID, TaxInclusivePrice: 90, EffectiveStartDate: date1, EffectiveEndDate: date8},
						{ID: "2", ProductID: productID, TaxInclusivePrice: 110, EffectiveStartDate: date9, EffectiveEndDate: date31},
					},
					ReviewScores: []persistance.ReviewScore{
						{ID: "1", ProductID: productID, Score: 1, Date: date9},
						{ID: "3", ProductID: productID, Score: 3, Date: date11},
					},
					Category: persistance.Category{ID: categoryID, Name: "カテゴリー1"},
				}
				suite.insertProducts([]persistance.Product{product})
			},
			Params: params{IDs: []string{"1"}, WithImage: false},
			Expected: expected{Products: []entity.Product{
				{
					ID:             "1",
					CategoryID:     "1",
					CategoryName:   "カテゴリー1",
					Name:           "商品名",
					ReviewScore:    0,
					Price:          100,
					SalePrice:      110,
					Description:    "商品説明",
					Status:         enum.SalesSuspend,
					StockCount:     10,
					Version:        1,
					CreateDateTime: date1,

					ProductImages: []entity.ProductImage{
						{ID: "1", ProductID: "1", Order: 1, Path: "/path1", Image: ""},
						{ID: "2", ProductID: "1", Order: 2, Path: "/path2", Image: ""},
					},
				},
			}, Err: nil},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.Name, func(t *testing.T) {
			defer suite.tearDown()

			tt.Setup(t)

			// when（操作）
			products, err := suite.productRepository.FindByIDs(suite.db, context.Background(), tt.Params.IDs, tt.Params.WithImage)

			// then（期待する結果）
			assert.Equal(t, tt.Expected.Products, products)
			assert.Equal(t, tt.Expected.Err, err)
		})
	}
}
