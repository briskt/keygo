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
	CreateToken(ctx echo.Context, input TokenCreateInput) (Token, error)

	// DeleteToken permanently deletes a token object from the system by ID.
	// The parent user object is not removed.
	DeleteToken(ctx echo.Context, id string) error

	// UpdateToken extends a token's ExpiresAt
	UpdateToken(ctx echo.Context, id string, input TokenUpdateInput) error
}

type Token struct {
	// TODO: remove private fields not appropriate for the API. (May require architecture changes.)
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

type TokenCreateInput struct {
	UserID    string
	AuthID    string
	ExpiresAt time.Time
}

// Validate returns an error if the struct contains invalid information
func (tc *TokenCreateInput) Validate() error {
	if tc.UserID == "" {
		return Errorf(ERR_INVALID, "UserID is required")
	}
	if tc.AuthID == "" {
		return Errorf(ERR_INVALID, "AuthID is required")
	}
	return nil
}

type TokenUpdateInput struct {
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
}

// Validate returns an error if the struct contains invalid information
func (tu *TokenUpdateInput) Validate() error {
	if tu.LastUsedAt != nil && tu.LastUsedAt.After(time.Now()) {
		return Errorf(ERR_INVALID, "LastUsedAt is in the future")
	}
	return nil
}
