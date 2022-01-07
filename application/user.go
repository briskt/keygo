package keygo

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	AvatarURL string    `db:"avatar_url"`
	Role      string    `db:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

// UserService represents a service for managing users
type UserService interface {
	// FindUserByID retrieves a user by ID
	FindUserByID(echo.Context, uuid.UUID) (User, error)

	// FindUsers retrieves a list of users by filter
	FindUsers(echo.Context, UserFilter) ([]User, int, error)

	// CreateUser creates a new user
	CreateUser(echo.Context, User) error

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
