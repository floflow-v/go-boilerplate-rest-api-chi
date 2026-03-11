package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrationFS embed.FS

func Run(db *sql.DB) error {

	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	return goose.Up(db, ".")
}
