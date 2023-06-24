package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/briskt/keygo/app"
)

type Auth struct {
	ID string `gorm:"primaryKey"`

	// User can have one or more methods of authentication
	// However, only one per provider is allowed per user
	UserID string
	User   User

	// The authentication provider & the provider's user ID
	Provider   string
	ProviderID string

	// Timestamps of creation & last update
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (a *Auth) BeforeCreate(tx *gorm.DB) error {
	a.ID = newID()
	return nil
}

// Validate returns an error if any fields are invalid on the Auth object.
func (a *Auth) Validate() error {
	if a.UserID == "" {
		return app.Errorf(app.ERR_INVALID, "user required")
	}
	if a.Provider == "" {
		return app.Errorf(app.ERR_INVALID, "provider required")
	}
	if a.ProviderID == "" {
		return app.Errorf(app.ERR_INVALID, "provider ID required")
	}
	return nil
}

// Ensure service implements interface
var _ app.AuthService = (*AuthService)(nil)

// AuthService represents a service for managing OAuth authentication
type AuthService struct{}

// NewAuthService returns a new instance of AuthService
func NewAuthService() *AuthService {
	return &AuthService{}
}

// FindAuthByID retrieves an authentication object by ID along with the associated user.
// Returns ERR_NOTFOUND if ID does not exist.
func (s *AuthService) FindAuthByID(ctx echo.Context, id string) (app.Auth, error) {
	auth, err := findAuthByID(ctx, id)
	if err != nil {
		return app.Auth{}, err
	}
	if err = auth.loadUser(ctx); err != nil {
		return app.Auth{}, err
	}

	return convertAuth(auth), nil
}

// FindAuths retrieves authentication objects based on a filter.
//
// Also returns the total number of objects that match the filter. This may
// differ from the returned object count if the Limit field is set.
func (s *AuthService) FindAuths(ctx echo.Context, filter app.AuthFilter) ([]app.Auth, int, error) {
	// Fetch the individual authentication objects from the database.
	auths, n, err := findAuths(ctx, filter)
	if err != nil {
		return []app.Auth{}, n, err
	}

	keygoAuths := make([]app.Auth, len(auths))
	for i, auth := range auths {
		if err = auth.loadUser(ctx); err != nil {
			return []app.Auth{}, n, err
		}
		keygoAuths[i] = convertAuth(auth)
	}
	return keygoAuths, n, nil
}

// CreateAuth Creates a new auth object. If a User is attached to the provided auth, then the created
// auth object is linked to the existing user, otherwise a new user object is created and linked.
//
// On success, the auth.ID is set to the new authentication ID
func (s *AuthService) CreateAuth(ctx echo.Context, keygoAuth app.Auth) (app.Auth, error) {
	auth := convertKeygoAuth(keygoAuth)

	// Check to see if the auth already exists for the given source.
	if other, err := findAuthByProviderID(ctx, auth.Provider, auth.ProviderID); err == nil {
		// If an auth already exists for the source user, update with the new tokens.
		if other, err = updateAuth(ctx, other.ID); err != nil {
			return app.Auth{}, fmt.Errorf("cannot create auth: id=%s err=%w", other.ID, err)
		} else if err = other.loadUser(ctx); err != nil {
			return app.Auth{}, err
		}

		return convertAuth(other), nil
	} else if app.ErrorCode(err) != app.ERR_NOTFOUND {
		return app.Auth{}, fmt.Errorf("canot find auth by source user: %w", err)
	}

	// Check if auth has a new user object passed in. It is considered "new" if
	// the caller doesn't know the database ID for the user.
	if auth.UserID == "" {
		// Look up the user by email address. If no user can be found then
		// create a new user with the auth.User object passed in.
		if user, err := findUserByEmail(ctx, auth.User.Email); err == nil { // user exists
			auth.User = user
		} else if app.ErrorCode(err) == app.ERR_NOTFOUND { // user does not exist
			if auth.User, err = createUser(ctx, auth.User); err != nil {
				return app.Auth{}, fmt.Errorf("could not create user for auth: %w", err)
			}
		} else {
			return app.Auth{}, fmt.Errorf("cannot find user by email: %w", err)
		}

		// Assign the created/found user ID back to the auth object.
		auth.UserID = auth.User.ID
	}

	// Create new auth object & attach associated user.
	err := auth.create(ctx)
	if err != nil {
		return app.Auth{}, err
	}

	if err = auth.loadUser(ctx); err != nil {
		return app.Auth{}, err
	}

	return convertAuth(auth), nil
}

// DeleteAuth permanently deletes an authentication object from the system by ID
// The parent user object is not removed
func (s *AuthService) DeleteAuth(ctx echo.Context, id string) error {
	return deleteAuth(ctx, id)
}

// findAuthByID is a helper function to return an auth object by ID
// Returns ERR_NOTFOUND if auth doesn't exist
func findAuthByID(ctx echo.Context, id string) (Auth, error) {
	var auth Auth
	result := Tx(ctx).Preload("User").Find(&auth, "id = ?", id)
	if result.Error == sql.ErrNoRows {
		return Auth{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Auth not found"}
	}
	return auth, result.Error
}

// findAuthByProviderID is a helper function to return an auth object by source ID.
// Returns ERR_NOTFOUND if auth doesn't exist.
func findAuthByProviderID(ctx echo.Context, provider, providerID string) (Auth, error) {
	var auth Auth
	err := Tx(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&auth).Error
	if err == gorm.ErrRecordNotFound {
		return Auth{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Auth not found"}
	}
	return auth, err
}

// findAuths returns a list of auth objects that match a filter. Also returns
// a total count of matches which may differ from results if filter.Limit is set.
func findAuths(ctx echo.Context, filter app.AuthFilter) (_ []Auth, n int, err error) {
	// TODO: implement query filter
	var auths []Auth
	result := Tx(ctx).Find(&auths)
	return auths, len(auths), result.Error
}

// create a new auth object in the database. On success, the ID is set to the new database
// ID & timestamp fields are set to the current time
func (a *Auth) create(ctx echo.Context) error {
	if err := a.Validate(); err != nil {
		return err
	}

	result := Tx(ctx).Omit("User").Create(a)
	return result.Error
}

// updateAuth updates tokens & expiry on exist auth object
// Returns new state of the auth object
func updateAuth(ctx echo.Context, id string) (Auth, error) {
	// Fetch current object state.
	auth, err := findAuthByID(ctx, id)
	if err != nil {
		return Auth{}, err
	}

	if err = auth.Validate(); err != nil {
		return auth, err
	}

	result := Tx(ctx).Omit("User").Save(&auth)
	return auth, result.Error
}

// deleteAuth permanently removes an auth object by ID
func deleteAuth(ctx echo.Context, id string) error {
	// Verify object exists & that the user is the owner of the auth
	//if auth, err := findAuthByID(tx, id); err != nil {
	//	return err
	//} else if auth.UserID != keygo.UserIDFromContext(ctx) {
	//	return keygo.Errorf(keygo.ERR_UNAUTHORIZED, "You are not allowed to delete this auth")
	//}

	result := Tx(ctx).Where("id = ?", id).Delete(&Auth{})
	return result.Error
}

// loadUser is a helper function to fetch & attach the associated user
// to the auth object.
func (a *Auth) loadUser(ctx echo.Context) (err error) {
	if a.User, err = findUserByID(ctx, a.UserID); err != nil {
		return fmt.Errorf("attach auth user: %w", err)
	}
	return nil
}

func convertAuth(a Auth) app.Auth {
	return app.Auth{
		ID:         a.ID,
		UserID:     a.UserID,
		User:       convertUser(a.User),
		Provider:   a.Provider,
		ProviderID: a.ProviderID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func convertKeygoAuth(a app.Auth) Auth {
	return Auth{
		ID:         a.ID,
		UserID:     a.UserID,
		User:       convertKeygoUser(a.User),
		Provider:   a.Provider,
		ProviderID: a.ProviderID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
