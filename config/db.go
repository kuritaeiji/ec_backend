package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/go-sql-driver/mysql"
)

func SetupDB() (*bun.DB, func(), error) {
	sql, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	db := bun.NewDB(sql, mysqldialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, func() {
		sql.Close()
		db.Close()
	}, nil
}
