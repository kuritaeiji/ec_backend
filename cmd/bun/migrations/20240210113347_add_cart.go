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
		_, err := db.NewCreateTable().Model(new(persistance.Cart)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewCreateTable().Model(new(persistance.CartProduct)).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		_, err := db.NewDropTable().Model(new(persistance.Cart)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		_, err = db.NewDropTable().Model(new(persistance.CartProduct)).IfExists().Exec(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
