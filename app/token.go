package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

// TokenService is a service for managing tokens
type TokenService interface {
	// FindToken looks up a token object by raw, unhashed token, and returns the Token object
	// with associated User
	// Returns ERR_NOTFOUND if token does not exist
	FindToken(ctx echo.Context, raw string) (Token, error)

	// ListTokensForUser returns all tokens for the given user
	ListTokensForUser(ctx echo.Context, userID string) ([]Token, error)

	// CreateToken creates a new token object
	//
	// On success, the token.ID is set to the new token ID
	CreateToken(ctx echo.Context, tokenCreate TokenCreate) (Token, error)

	// DeleteToken permanently deletes a token object from the system by ID.
	// The parent user object is not removed.
	DeleteToken(ctx echo.Context, id string) error

	// UpdateToken extends a token's ExpiresAt
	UpdateToken(ctx echo.Context, id string) error
}

type Token struct {
	ID string

	User   User
	UserID string

	AuthID    string
	PlainText string

	LastUsedAt *time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TokenCreate struct {
	UserEmail string
	AuthID    string
}
