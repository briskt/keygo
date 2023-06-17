package app

import (
	"github.com/labstack/echo/v4"
)

const (
	// ContextKeyUser stores the current logged-in user
	ContextKeyUser = "user"

	// ContextKeyTx stores the database transaction
	ContextKeyTx = "tx"

	// ContextKeyToken stores the Token passed by the client
	ContextKeyToken = "token"
)

// NewContextWithUser returns a new context with the given user.
func NewContextWithUser(ctx echo.Context, user User) echo.Context {
	ctx.Set(ContextKeyUser, user)
	return ctx
}

// CurrentUser returns the current logged-in user.
func CurrentUser(ctx echo.Context) User {
	user, _ := ctx.Get(ContextKeyUser).(User)
	return user
}
