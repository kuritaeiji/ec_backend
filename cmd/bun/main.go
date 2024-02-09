package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kuritaeiji/ec_backend/cmd/bun/migrations"
	"github.com/kuritaeiji/ec_backend/config"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	err := config.SetupEnv()
	if err != nil {
		log.Fatal(err)
	}

	db := config.NewDB()

	app := cli.App{
		Name: "migration",
		Usage: "bun migration",
		Commands: newMigrationCommands(migrate.NewMigrator(db, migrations.Migrations)),
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func exit(err error) error {
	if err != nil {
		return cli.Exit(err, 1)
	} else {
		return nil
	}
}

func newMigrationCommands(migrator *migrate.Migrator) []*cli.Command {
	return []*cli.Command{
		{
			Name: "init",
			Usage: "create migration table",
			Action: func(ctx *cli.Context) error {
				if err := migrator.Init(ctx.Context); err != nil {
					return cli.Exit(err, 1)
				}

				fmt.Println("マイグレーションテーブルを作成しました")
				return nil
			},
		},
		{
			Name: "migrate",
			Usage: "migrate database",
			Action: func(ctx *cli.Context) error {
				if err := migrator.Lock(ctx.Context); err != nil {
					return exit(err)
				}
				defer func() {
					err := migrator.Unlock(ctx.Context)
					if err != nil {
						fmt.Println(err)
					}
				}()

				group, err := migrator.Migrate(ctx.Context)
				if err != nil {
					return exit(err)
				}
				if group.IsZero() {
					fmt.Println("migrateするテーブルが存在しません")
					return nil
				}

				fmt.Printf("マイグレーションに成功しました %s\n", group)
				return nil
			},
		},
		{
			Name: "rollback",
			Usage: "rollback the last migration group",
			Action: func(ctx *cli.Context) error {
				if err := migrator.Lock(ctx.Context); err != nil {
					return exit(err)
				}
				defer func() {
					err := migrator.Unlock(ctx.Context)
					if err != nil {
						fmt.Println(err)
					}
				}()

				group, err := migrator.Rollback(ctx.Context)
				if err != nil {
					return exit(err)
				}

				if group.IsZero() {
					fmt.Println("ロールバックするグループが存在しません")
					return nil
				}

				fmt.Printf("ロールバックしました %s\n", group)
				return nil
			},
		},
		{
			Name: "create_migration_file",
			Usage: "create migration file",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "sql",
					Usage: "create sql migration file",
					Value: false,
				},
			},
			Action: func(ctx *cli.Context) error {

				fname := strings.Join(ctx.Args().Slice(), "_")
				if ctx.Bool("sql") {
					mfs, err := migrator.CreateSQLMigrations (ctx.Context, fname)
					if err != nil {
						return exit(err)
					}
					for _, mf := range mfs {
						fmt.Printf("作成したマイグレーションファイル: %s (%s)\n", mf.Name, mf.Path)
					}
					return nil
				} else {
					mf, err := migrator.CreateGoMigration(ctx.Context, fname)
					if err != nil {
						return exit(err)
					}

					fmt.Printf("作成したマイグレーションファイル: %s (%s)\n", mf.Name, mf.Path)
					return nil
				}
			},
		},
		{
			Name: "status",
			Usage: "print migration status",
			Action: func(ctx *cli.Context) error {

				ms, err := migrator.MigrationsWithStatus(ctx.Context)
				if err != nil {
					return exit(err)
				}

				fmt.Printf("マイグレーション: %s\n", ms)
				fmt.Printf("適用していないマイグレーション: %s\n", ms.Unapplied())
				fmt.Printf("最後のマイグレーショングループ: %s\n", ms.LastGroup())
				return nil
			},
		},
	}
}