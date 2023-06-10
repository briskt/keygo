package server

import (
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/briskt/keygo"
)

func TxMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	errNotOK := errors.New("http error, rolling back transaction")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := db.Transaction(func(tx *gorm.DB) error {
				c.Set(keygo.ContextKeyTx, tx)

				if err := next(c); err != nil {
					return err
				}

				// If the status is not a "success", roll back transaction by returning an error
				res := c.Response()

				// let 200s and 300s through
				if res.Status < 200 || res.Status >= 400 {
					return errNotOK
				}

				return nil
			})
			if err != nil {
				if errors.Unwrap(err) == errNotOK {
					return nil
				}
				return err
			}
			return nil
		}
	}
}
