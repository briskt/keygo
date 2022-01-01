package http

import (
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Transaction(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			couldBeDBorYourErr := db.Transaction(func(tx *gorm.DB) error {

				// add the transaction to the context
				c.Set("tx", tx)

				// call the next handler; if it errors stop and return the error
				if yourError := next(c); yourError != nil {
					return yourError
				}

				// check the response status code. if the code is NOT 200..399
				// then it is considered "NOT SUCCESSFUL" and an error will be returned
				res := c.Response()
				if res.Status < 200 || res.Status >= 400 {
					return errNonSuccess
				}

				// return nil will commit the whole transaction
				return nil
			})

			// couldBeDBorYourErr could be one of possible values:
			// * nil - everything went well, if so, return
			// * yourError - an error returned from your application, middleware, etc...
			// * a database error - this is returned if there were problems committing the transaction
			// * a errNonSuccess - this is returned if the response status code is not between 200..399
			if couldBeDBorYourErr != nil && errors.Unwrap(couldBeDBorYourErr) != errNonSuccess {
				return couldBeDBorYourErr
			}
			return nil
		}
	}
}

var errNonSuccess = errors.New("non-success error code")
