package app

import (
	"time"
)

const (
	UserRoleBasic = "Basic"
	UserRoleAdmin = "Admin"
)

// UserFilter is a filter passed to FindUsers()
type UserFilter struct {
	// Filtering fields.
	Email    *string
	TenantID *string

	// Restrict to subset of results.
	Offset int
	Limit  int
}

// User is the full model that identifies an app User
type User struct {
	ID          string
	FirstName   string
	LastName    string
	Email       string
	AvatarURL   string
	Role        string
	TenantID    string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserCreateInput is a set of fields to define a new user for CreateUser()
type UserCreateInput struct {
	FirstName string
	LastName  string
	Email     string
	AvatarURL string
	Role      string
	TenantID  string
}

// Validate returns an error if the struct contains invalid information
func (uc *UserCreateInput) Validate() error {
	if uc.Email == "" {
		return Errorf(ERR_INVALID, "Email is required")
	}
	return nil
}

// UserUpdateInput is a set of fields to be updated via UpdateUser()
type UserUpdateInput struct {
	FirstName *string
	LastName  *string
	Email     *string
}

// Validate returns an error if the struct contains invalid information.
func (uu *UserUpdateInput) Validate() error {
	if uu.Email != nil && *uu.Email == "" {
		return Errorf(ERR_INVALID, "Email is required")
	}
	return nil
}
