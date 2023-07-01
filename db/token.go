package db

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/briskt/keygo/app"
)

const (
	tokenLifetime = time.Hour * 24
	tokenBytes    = 32
)

type Token struct {
	ID string `gorm:"primaryKey"`

	User   User
	UserID string

	AuthID    string // OAuth sub (subject)
	Hash      string
	PlainText string `gorm:"-"`

	LastUsedAt *time.Time
	ExpiresAt  time.Time // FIXME: change to pointer
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	t.ID = newID()
	return nil
}

// Validate returns an error if any fields are invalid on the Token object.
func (t *Token) Validate() error {
	if t.UserID == "" {
		return app.Errorf(app.ERR_INVALID, "UserID required.")
	}
	if t.AuthID == "" {
		return app.Errorf(app.ERR_INVALID, "AuthID required.")
	}
	if t.Hash == "" {
		return app.Errorf(app.ERR_INVALID, "Hash required.")
	}
	return nil
}

// create a new token object in the database. On success, the ID is set to the new database
// ID & timestamp fields are set to the current time
func (t *Token) create(ctx echo.Context) error {
	t.ExpiresAt = time.Now().Add(tokenLifetime)
	t.LastUsedAt = time.Now()
	t.PlainText = getRandomToken()
	t.Hash = hashToken(t.PlainText)

	if err := t.Validate(); err != nil {
		return err
	}

	err := Tx(ctx).Omit("User").Create(t).Error
	return err
}

func getRandomToken() string {
	rb := make([]byte, tokenBytes)

	_, err := rand.Read(rb)
	if err != nil {
		panic("rand.Read failed in getRandomToken, " + err.Error())
	}

	return base64.URLEncoding.EncodeToString(rb)
}

func hashToken(accessToken string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(accessToken)))
}

// findToken is a helper function to return a token object by unhashed token string
// Returns ERR_NOTFOUND if record doesn't exist
func findToken(ctx echo.Context, raw string) (Token, error) {
	var token Token
	err := Tx(ctx).Where("hash = ?", hashToken(raw)).First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return Token{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Token not found"}
	}
	return token, err
}

// findToken is a helper function to return a token object by AuthID
// Returns ERR_NOTFOUND if record doesn't exist
func findTokenByAuthID(ctx echo.Context, authID string) (Token, error) {
	var token Token
	err := Tx(ctx).Where("auth_id = ?", authID).First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return Token{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Token not found"}
	}
	return token, err
}

// updateToken updates expires_at and last_used_at on existing token object
// Returns new state of the token object
// FIXME: this is never called. It should be called in the authentication middleware.
func updateToken(ctx echo.Context, token Token) (Token, error) {
	if err := token.Validate(); err != nil {
		return token, err
	}

	token.ExpiresAt = time.Now().Add(tokenLifetime)
	token.LastUsedAt = time.Now()

	result := Tx(ctx).Omit("User").Save(&token)
	return token, result.Error
}

// deleteToken permanently removes a token object by ID
func deleteToken(ctx echo.Context, id string) error {
	// Verify object exists & that the user is the owner of the token
	//if token, err := findTokenByID(tx, id); err != nil {
	//	return err
	//} else if token.UserID != keygo.UserIDFromContext(ctx) {
	//	return keygo.Errorf(keygo.ERR_UNAUTHORIZED, "You are not allowed to delete this token")
	//}

	result := Tx(ctx).Where("id = ?", id).Delete(&Token{})
	return result.Error
}

// loadUser is a helper function to fetch & attach the associated User
// to the token object.
func (t *Token) loadUser(ctx echo.Context) (err error) {
	if t.User, err = findUserByID(ctx, t.UserID); err != nil {
		return fmt.Errorf("attach token user: %w", err)
	}
	return nil
}

// Ensure service implements interface.
var _ app.TokenService = (*TokenService)(nil)

// TokenService represents a service for managing API auth tokens
type TokenService struct{}

// NewTokenService returns a new instance of TokenService
func NewTokenService() *TokenService {
	return &TokenService{}
}

func (t TokenService) FindToken(ctx echo.Context, raw string) (app.Token, error) {
	token, err := findToken(ctx, raw)
	if err != nil {
		return app.Token{}, err
	}
	if err = token.loadUser(ctx); err != nil {
		return app.Token{}, err
	}

	return convertToken(token), nil
}

// CreateToken creates a new token object. If a User is attached to the provided token, then the created
// token object is linked to the existing user, otherwise a new user object is created and linked.
//
// On success, the token.ID is set to the new token ID
func (t TokenService) CreateToken(ctx echo.Context, appToken app.Token) (app.Token, error) {
	token := convertAppToken(appToken)

	// Check if token has a new user object passed in. It is considered "new" if
	// the caller doesn't know the database ID for the user.
	if token.UserID == "" {
		// Look up the user by email address. If no user can be found then
		// create a new user with the token.User object passed in.
		if user, err := findUserByEmail(ctx, token.User.Email); err == nil { // user exists
			token.User = user
		} else if app.ErrorCode(err) == app.ERR_NOTFOUND {
			// user does not exist with the given email address -- create a new user
			if token.User, err = createUser(ctx, token.User); err != nil {
				return app.Token{}, fmt.Errorf("could not create user for token: %w", err)
			}
		} else {
			return app.Token{}, fmt.Errorf("cannot find user by email: %w", err)
		}

		// Assign the created/found user ID back to the token object.
		token.UserID = token.User.ID
	}

	// Create new token object & attach associated user.
	err := token.create(ctx)
	if err != nil {
		return app.Token{}, err
	}

	if err = token.loadUser(ctx); err != nil {
		return app.Token{}, err
	}

	return convertToken(token), nil
}

func (t TokenService) DeleteToken(ctx echo.Context, id string) error {
	return deleteToken(ctx, id)
}

func convertToken(token Token) app.Token {
	return app.Token{
		ID:         token.ID,
		User:       convertUser(token.User),
		UserID:     token.UserID,
		AuthID:     token.AuthID,
		PlainText:  token.PlainText,
		LastUsedAt: token.LastUsedAt,
		ExpiresAt:  token.ExpiresAt,
		CreatedAt:  token.CreatedAt,
		UpdatedAt:  token.UpdatedAt,
	}
}

func convertAppToken(token app.Token) Token {
	return Token{
		ID:         token.ID,
		User:       convertKeygoUser(token.User),
		UserID:     token.UserID,
		AuthID:     token.AuthID,
		PlainText:  token.PlainText,
		LastUsedAt: token.LastUsedAt,
		ExpiresAt:  token.ExpiresAt,
		CreatedAt:  token.CreatedAt,
		UpdatedAt:  token.UpdatedAt,
	}
}
