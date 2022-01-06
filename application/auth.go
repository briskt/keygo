package keygo

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`

	// User can have one or more methods of authentication
	// However, only one per provider is allowed per user
	UserID uuid.UUID `json:"userID"`
	User   User      `json:"user"`

	// Avatar URL provided by the auth provider
	AvatarURL string `json:"avatarURL"`

	// The authentication provider & the provider's user ID
	Provider   string `json:"provider"`
	ProviderID string `json:"providerID"`

	// OAuth fields returned from the authentication provider
	// Not all providers use refresh tokens.
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	Expiry       time.Time `json:"-"`

	// Timestamps of creation & last update
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (a *Auth) BeforeCreate(tx *gorm.DB) error {
	a.ID = NewUUID()
	return nil
}

// Validate returns an error if any fields are invalid on the Auth object.
// This can be called by the database implementation to do some basic checks.
func (a *Auth) Validate() error {
	if a.UserID == uuid.Nil {
		return Errorf(ERR_INVALID, "User required.")
	} else if a.Provider == "" {
		return Errorf(ERR_INVALID, "Provider required.")
	} else if a.ProviderID == "" {
		return Errorf(ERR_INVALID, "Provider ID required.")
	} else if a.AccessToken == "" {
		return Errorf(ERR_INVALID, "Access token required.")
	}
	return nil
}

// GetAvatarURL returns a URL to the avatar image hosted by the authentication provider
// Returns an empty string if the authentication provider is invalid
func (a *Auth) GetAvatarURL(size int) string {
	switch a.Provider {
	case AuthProviderGoogle:
		return a.AvatarURL
	default:
		return ""
	}
}

// AuthService represents a service for managing auths
type AuthService interface {
	// FindAuthByID looks up an authentication object by ID along with the associated user
	// Returns ERR_NOTFOUND if ID does not exist
	FindAuthByID(echo.Context, uuid.UUID) (Auth, error)

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
	DeleteAuth(echo.Context, uuid.UUID) error
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
