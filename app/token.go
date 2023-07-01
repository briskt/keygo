package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

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

// TokenService represents a service for managing tokens
type TokenService interface {
	// FindToken looks up a token object by raw, unhashed token, and returns the Token object
	// with associated User
	// Returns ERR_NOTFOUND if token does not exist
	FindToken(ctx echo.Context, raw string) (Token, error)

	// CreateToken creates a new token object
	//
	// On success, the token.ID is set to the new token ID
	CreateToken(echo.Context, Token) (Token, error)

	// DeleteToken permanently deletes a token object from the system by ID.
	// The parent user object is not removed.
	DeleteToken(ctx echo.Context, tokenID string) error
}
