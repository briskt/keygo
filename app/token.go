package app

import (
	"time"
)

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
