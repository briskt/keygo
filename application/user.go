package keygo

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UserService represents a service for managing users
type UserService interface {
	// FindUserByID retrieves a user by ID
	FindUserByID(id uuid.UUID) (User, error)

	// FindUsers retrieves a list of users by filter
	FindUsers(filter UserFilter) ([]User, int, error)

	// CreateUser creates a new user
	CreateUser(user User) error

	// UpdateUser updates a user object
	UpdateUser(id uuid.UUID, upd UserUpdate) (User, error)

	// DeleteUser permanently deletes a user and all child objects
	DeleteUser(id uuid.UUID) error
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
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
}
