package migrations

import (
	"log"
	"os"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(os.DirFS(".")); err != nil {
		log.Fatal(err)
	}
}