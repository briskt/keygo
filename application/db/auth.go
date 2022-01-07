package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
)

// Ensure service implements interface.
var _ keygo.AuthService = (*AuthService)(nil)

// AuthService represents a service for managing OAuth authentication.
type AuthService struct{}

// NewAuthService returns a new instance of AuthService attached to DB.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// FindAuthByID retrieves an authentication object by ID along with the associated user.
// Returns ERR_NOTFOUND if ID does not exist.
func (s *AuthService) FindAuthByID(ctx echo.Context, id uuid.UUID) (keygo.Auth, error) {
	auth, err := findAuthByID(ctx, id)
	if err != nil {
		return auth, err
	} else if err = attachAuthAssociations(ctx, auth); err != nil {
		return auth, err
	}

	return auth, nil
}

// FindAuths retrieves authentication objects based on a filter.
//
// Also returns the total number of objects that match the filter. This may
// differ from the returned object count if the Limit field is set.
func (s *AuthService) FindAuths(ctx echo.Context, filter keygo.AuthFilter) ([]keygo.Auth, int, error) {
	// Fetch the individual authentication objects from the database.
	auths, n, err := findAuths(ctx, filter)
	if err != nil {
		return auths, n, err
	}

	// Iterate over returned objects and attach user objects.
	// This works well for SQLite because it is in-process but remote database
	// servers will incur a high per-query latency so queries should be batched.
	for _, auth := range auths {
		if err = attachAuthAssociations(ctx, auth); err != nil {
			return auths, n, err
		}
	}
	return auths, n, nil
}

// CreateAuth Creates a new authentication object If a User is attached to auth,
// then the auth object is linked to an existing user, otherwise a new user
// object is created.
//
// On success, the auth.ID is set to the new authentication ID
func (s *AuthService) CreateAuth(ctx echo.Context, auth keygo.Auth) (keygo.Auth, error) {
	// Check to see if the auth already exists for the given source.
	if other, err := findAuthByProviderID(ctx, auth.Provider, auth.ProviderID); err == nil {
		// If an auth already exists for the source user, update with the new tokens.
		if other, err = updateAuth(ctx, other.ID, auth.AccessToken, auth.RefreshToken, auth.Expiry); err != nil {
			return keygo.Auth{}, fmt.Errorf("cannot create auth: id=%d err=%w", other.ID, err)
		} else if err = attachAuthAssociations(ctx, other); err != nil {
			return keygo.Auth{}, err
		}

		return other, nil
	} else if keygo.ErrorCode(err) != keygo.ERR_NOTFOUND {
		return keygo.Auth{}, fmt.Errorf("canot find auth by source user: %w", err)
	}

	// Check if auth has a new user object passed in. It is considered "new" if
	// the caller doesn't know the database ID for the user.
	if auth.UserID == uuid.Nil {
		// Look up the user by email address. If no user can be found then
		// create a new user with the auth.User object passed in.
		if user, err := findUserByEmail(ctx, auth.User.Email); err == nil { // user exists
			auth.User = user
		} else if keygo.ErrorCode(err) == keygo.ERR_NOTFOUND { // user does not exist
			if auth.User, err = createUser(ctx, auth.User); err != nil {
				return keygo.Auth{}, fmt.Errorf("cannot create user: %w", err)
			}
		} else {
			return keygo.Auth{}, fmt.Errorf("cannot find user by email: %w", err)
		}

		// Assign the created/found user ID back to the auth object.
		auth.UserID = auth.User.ID
	}

	// Create new auth object & attach associated user.
	newAuth, err := createAuth(ctx, auth)
	if err != nil {
		return keygo.Auth{}, err
	}

	return newAuth, attachAuthAssociations(ctx, auth)
}

// DeleteAuth permanently deletes an authentication object from the system by ID
// The parent user object is not removed
func (s *AuthService) DeleteAuth(ctx echo.Context, id uuid.UUID) error {
	return deleteAuth(ctx, id)
}

// findAuthByID is a helper function to return an auth object by ID
// Returns ERR_NOTFOUND if auth doesn't exist
func findAuthByID(ctx echo.Context, id uuid.UUID) (keygo.Auth, error) {
	var auth keygo.Auth
	result := Tx(ctx).Find(&auth, id)
	if result.Error == sql.ErrNoRows {
		return keygo.Auth{}, &keygo.Error{Code: keygo.ERR_NOTFOUND, Message: "Auth not found"}
	}
	return auth, result.Error
}

// findAuthByProviderID is a helper function to return an auth object by source ID.
// Returns ERR_NOTFOUND if auth doesn't exist.
func findAuthByProviderID(ctx echo.Context, provider, providerID string) (keygo.Auth, error) {
	var auth keygo.Auth
	err := Tx(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&auth).Error
	if err == gorm.ErrRecordNotFound {
		return keygo.Auth{}, &keygo.Error{Code: keygo.ERR_NOTFOUND, Message: "Auth not found"}
	}
	return auth, err
}

// findAuths returns a list of auth objects that match a filter. Also returns
// a total count of matches which may differ from results if filter.Limit is set.
func findAuths(ctx echo.Context, filter keygo.AuthFilter) (_ []keygo.Auth, n int, err error) {
	// TODO: implement query filter
	var auths []keygo.Auth
	result := Tx(ctx).Find(&auths)
	return auths, len(auths), result.Error
}

// createAuth creates a new auth object in the database. On success, the
// ID is set to the new database ID & timestamp fields are set to the current time.
func createAuth(ctx echo.Context, auth keygo.Auth) (keygo.Auth, error) {
	if err := auth.Validate(); err != nil {
		return keygo.Auth{}, err
	}

	result := Tx(ctx).Omit("User").Create(&auth)
	return auth, result.Error
}

// updateAuth updates tokens & expiry on exist auth object
// Returns new state of the auth object
func updateAuth(ctx echo.Context, id uuid.UUID, accessToken, refreshToken string, expiry time.Time) (keygo.Auth, error) {
	// Fetch current object state.
	auth, err := findAuthByID(ctx, id)
	if err != nil {
		return keygo.Auth{}, err
	}

	// Update fields
	auth.AccessToken = accessToken
	auth.RefreshToken = refreshToken
	auth.Expiry = expiry

	if err = auth.Validate(); err != nil {
		return auth, err
	}

	result := Tx(ctx).Save(&auth)
	return auth, result.Error
}

// deleteAuth permanently removes an auth object by ID
func deleteAuth(ctx echo.Context, id uuid.UUID) error {
	// Verify object exists & that the user is the owner of the auth
	//if auth, err := findAuthByID(tx, id); err != nil {
	//	return err
	//} else if auth.UserID != keygo.UserIDFromContext(ctx) {
	//	return keygo.Errorf(keygo.ERR_UNAUTHORIZED, "You are not allowed to delete this auth")
	//}

	auth := keygo.User{ID: id}
	result := Tx(ctx).Delete(&auth)
	return result.Error
}

// attachAuthAssociations is a helper function to fetch & attach the associated user
// to the auth object.
func attachAuthAssociations(ctx echo.Context, auth keygo.Auth) (err error) {
	if auth.User, err = findUserByID(ctx, auth.UserID); err != nil {
		return fmt.Errorf("attach auth user: %w", err)
	}
	return nil
}
