package keygo

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	// Stores the current logged-in user in the context.
	userContextKey = "user"
)

// NewContextWithUser returns a new context with the given user.
func NewContextWithUser(ctx echo.Context, user User) echo.Context {
	ctx.Set(userContextKey, user)
	return ctx
}

// UserFromContext returns the current logged-in user.
func UserFromContext(ctx echo.Context) User {
	user, _ := ctx.Get(userContextKey).(User)
	return user
}

// UserIDFromContext is a helper function that returns the ID of the current
// logged-in user. Returns zero if no user is logged in.
func UserIDFromContext(ctx echo.Context) uuid.UUID {
	user := UserFromContext(ctx)
	return user.ID
}
