package keygo

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Token struct {
	ID uuid.UUID `json:"id"`

	Auth   Auth      `json:"auth"`
	AuthID uuid.UUID `json:"authID"`

	PlainText string `json:"plainText"`

	LastLoginAt time.Time `json:"lastLoginAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TokenService represents a service for managing tokens
type TokenService interface {
	// FindToken looks up a token object by raw, unhashed token, and returns the Token object
	// with associated Auth and Auth.User
	// Returns ERR_NOTFOUND if token does not exist
	FindToken(ctx echo.Context, raw string) (Token, error)

	// CreateToken creates a new token object
	//
	// On success, the token.ID is set to the new token ID
	CreateToken(ctx echo.Context, authID uuid.UUID, clientID string) (Token, error)

	// DeleteToken permanently deletes a token object from the system by ID.
	// The parent user object is not removed.
	DeleteToken(ctx echo.Context, tokenID uuid.UUID) error
}
