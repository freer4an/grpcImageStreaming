package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func MakeMigrations(dbUrl string) error {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return fmt.Errorf("sql open %w", err)
	}
	defer db.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up %w", err)
	}
	return nil
}
