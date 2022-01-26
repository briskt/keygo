package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func MigrateUp() {
	fmt.Println("migrating up")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("failed to open DB, " + err.Error())
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic("failed to close DB, " + err.Error())
		}
	}()

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "."); err != nil {
		panic("goose up failed, " + err.Error())
	}
}

func MigrateDown() {
	fmt.Println("migrating down")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("failed to open DB, " + err.Error())
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic("failed to close DB, " + err.Error())
		}
	}()

	goose.SetBaseFS(embedMigrations)

	if err := goose.DownTo(db, ".", 0); err != nil {
		panic("goose down failed, " + err.Error())
	}
}
