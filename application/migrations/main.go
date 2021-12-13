package main

import (
	"database/sql"
	"embed"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func main() {
	var db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
