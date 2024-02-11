package migrations

import (
	"context"
	"fmt"

	"github.com/kuritaeiji/ec_backend/enduser/infrastructure/persistance"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		_, err := db.NewCreateTable().Model(new(persistance.Product)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.ProductStatus)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.ProductImage)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.ProductPrice)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.ProductSalePrice)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.ReviewScore)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.Category)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().Model(new(persistance.Product)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.ProductStatus)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.ProductImage)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.ProductPrice)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.ProductSalePrice)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.ReviewScore)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.Category)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
