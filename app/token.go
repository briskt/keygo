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
	ID string `json:"id"`

	User   User   `json:"user"`
	UserID string `json:"userID"`

	AuthID    string `json:"authID"`
	PlainText string `json:"plainText"`

	LastUsedAt *time.Time `json:"lastUsedAt"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type TokenCreate struct {
	UserID string
	AuthID string
}
