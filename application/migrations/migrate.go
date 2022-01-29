package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Fresh(db *sql.DB) {
	fmt.Println("migrating down/up")
	Down(db)
	Up(db)
}

func Up(db *sql.DB) {
	fmt.Println("migrating up")

	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(db, "."); err != nil {
		panic("goose up failed: " + err.Error())
	}
}

func Down(db *sql.DB) {
	fmt.Println("migrating down")

	goose.SetBaseFS(embedMigrations)
	if err := goose.DownTo(db, ".", 0); err != nil {
		panic("goose down failed: " + err.Error())
	}
}
