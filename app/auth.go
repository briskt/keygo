package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Authentication providers
const (
	AuthProviderGoogle = "google"
)

// Auth represents a set of OAuth credentials. These are linked to a User so a
// single user could authenticate through multiple providers.
//
// The authentication system links users by email address, however, some users
// on some providers don't provide their public email so we may not be able to link them
// by email address.
type Auth struct {
	ID string `json:"id"`

	// User can have one or more methods of authentication
	// However, only one per provider is allowed per user
	UserID string `json:"userID"`
	User   User   `json:"user"`

	// The authentication provider
	Provider string `json:"provider"`

	// The user's ID for the provider
	ProviderID string `json:"providerID"`

	// Timestamps of creation & last update
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AuthService represents a service for managing auths
type AuthService interface {
	// FindAuthByID looks up an authentication object by ID along with the associated user
	// Returns ERR_NOTFOUND if ID does not exist
	FindAuthByID(echo.Context, string) (Auth, error)

	// FindAuths retrieves authentication objects based on a filter. Also returns the
	// total number of objects that match the filter. This may differ from the
	// returned object count if the Limit field is set.
	FindAuths(echo.Context, AuthFilter) ([]Auth, int, error)

	// CreateAuth creates a new authentication object If a User is attached to auth, then
	// the auth object is linked to an existing user, otherwise a new user
	// object is created.
	//
	// On success, the auth.ID is set to the new authentication ID
	CreateAuth(echo.Context, Auth) (Auth, error)

	// DeleteAuth permanently deletes an authentication object from the system by ID.
	// The parent user object is not removed.
	DeleteAuth(echo.Context, string) error
}

// AuthFilter represents a filter accepted by FindAuths()
type AuthFilter struct {
	// Filtering fields
	ID         *int    `json:"id"`
	UserID     *int    `json:"userID"`
	Provider   *string `json:"provider"`
	ProviderID *string `json:"providerID"`

	// Pagination parameters
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// swagger:model
type AuthStatus struct {
	// IsAuthenticated is true when the supplied session cookie is valid and references a valid user
	IsAuthenticated bool `json:"IsAuthenticated"`

	// Expiry is the date and time when the session is scheduled to expire. It is invalid if `IsAuthenticated` is false.
	//
	// swagger:strfmt date-time
	Expiry time.Time `json:"Expiry"`

	// UserID is the ID of the authenticated user. It is invalid if `IsAuthenticated` is false.
	UserID string `json:"UserID"`
}
