package db_test

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/schparky/keygo/migrations"
)

// Ensure the test database can open & close.
func TestDB(t *testing.T) {
	_ = MustOpenDB(t)
}

// MustOpenDB returns a new, open DB. Fatal on error.
func MustOpenDB(tb testing.TB) *gorm.DB {
	tb.Helper()

	tx, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic("error opening database, " + err.Error())
	}

	migrations.MigrateDown()
	migrations.MigrateUp()
	return tx
}
