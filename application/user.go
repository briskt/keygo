package keygo

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatarURL"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate returns an error if the user contains invalid fields.
// This only performs basic validation.
func (u *User) Validate() error {
	if u.FirstName == "" {
		return Errorf(ERR_INVALID, "FirstName required")
	}
	if u.Email == "" {
		return Errorf(ERR_INVALID, "Email required")
	}
	return nil
}

// UserService represents a service for managing users
type UserService interface {
	// FindUserByID retrieves a user by ID
	FindUserByID(echo.Context, uuid.UUID) (User, error)

	// FindUsers retrieves a list of users by filter
	FindUsers(echo.Context, UserFilter) ([]User, int, error)

	// CreateUser creates a new user
	CreateUser(echo.Context, User) (User, error)

	// UpdateUser updates a user object
	UpdateUser(echo.Context, uuid.UUID, UserUpdate) (User, error)

	// DeleteUser permanently deletes a user and all child objects
	DeleteUser(echo.Context, uuid.UUID) error
}

// UserFilter represents a filter passed to FindUsers()
type UserFilter struct {
	// Filtering fields.
	ID     *uuid.UUID `json:"id"`
	Email  *string    `json:"email"`
	APIKey *string    `json:"apiKey"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser()
type UserUpdate struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Email     *string `json:"email"`
}
