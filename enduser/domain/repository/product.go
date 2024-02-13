package repository

import (
	"context"

	"github.com/kuritaeiji/ec_backend/enduser/domain/entity"
	"github.com/uptrace/bun"
)

type ProductRepository interface {
	// 商品ID配列に一致する商品配列を返却する。引数withImageがtrueの場合はS3から画像を取得し、そうでない場合は取得しない。
	FindByIDs(db bun.IDB, ctx context.Context, ids []string, withImage bool) ([]entity.Product, error)
}
