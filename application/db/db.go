package db

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/briskt/keygo"
)

var DB *gorm.DB

func init() {
	DB = OpenDB()
}

func OpenDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("required environment variable DATABASE_URL is not set")
	}

	config := gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)}
	conn, err := gorm.Open(postgres.Open(dsn), &config)
	if err != nil {
		panic("failed to open database '" + dsn + "': " + err.Error())
	}
	return conn
}

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	}
	if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	}
	if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}

func Tx(ctx echo.Context) *gorm.DB {
	tmp := ctx.Get(keygo.ContextKeyTx)
	tx, ok := tmp.(*gorm.DB)
	if !ok {
		panic("no transaction found in context")
	}
	return tx
}

// NullTime represents a helper wrapper for time.Time. It automatically converts
// time fields to/from RFC 3339 format. Also supports NULL for zero time.
type NullTime time.Time

// Scan reads a time value from the database.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	}
	if s, ok := value.(string); ok {
		*(*time.Time)(n), _ = time.Parse(time.RFC3339, s)
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

// Value formats a time value for the database.
func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format(time.RFC3339), nil
}
