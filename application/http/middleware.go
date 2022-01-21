package http

import (
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const ContextKeyTx = "tx"

func Transaction(db *gorm.DB) echo.MiddlewareFunc {
	errNotOK := errors.New("http error caught in transaction middleware")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := db.Transaction(func(tx *gorm.DB) error {
				c.Set(ContextKeyTx, tx)

				if err := next(c); err != nil {
					return err
				}

				// If the status is not a "success", roll back transaction by returning an error
				res := c.Response()
				// let 200-series through
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
