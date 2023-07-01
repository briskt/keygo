package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	AvatarURL   string
	Role        string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	FindUserByID(ctx echo.Context, id string) (User, error)

	// FindUsers retrieves a list of users by filter
	FindUsers(ctx echo.Context, userFilter UserFilter) ([]User, int, error)

	// CreateUser creates a new user
	CreateUser(ctx echo.Context, user User) (User, error)

	// UpdateUser updates a user object
	UpdateUser(ctx echo.Context, id string, userUpdate UserUpdate) (User, error)

	// DeleteUser permanently deletes a user and all child objects
	DeleteUser(ctx echo.Context, id string) error

	// TouchLastLoginAt sets the LastLoginAt field to the current time
	TouchLastLoginAt(ctx echo.Context, id string) error
}

// UserFilter represents a filter passed to FindUsers()
type UserFilter struct {
	// Filtering fields.
	ID     *string `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

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
