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
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Deleted    gorm.DeletedAt
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
	t.touch()
	t.PlainText = randomString()
	t.Hash = hashToken(t.PlainText)

	if err := t.Validate(); err != nil {
		return err
	}

	err := Tx(ctx).Omit("User").Create(t).Error
	return err
}

func (t *Token) touch() {
	now := time.Now()
	t.ExpiresAt = now.Add(tokenLifetime)
	t.LastUsedAt = &now
}

func randomString() string {
	rb := make([]byte, tokenBytes)

	_, err := rand.Read(rb)
	if err != nil {
		panic("rand.Read failed in randomString, " + err.Error())
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

// findToken is a helper function to return a token object by its ID
// Returns ERR_NOTFOUND if record doesn't exist
func findTokenByID(ctx echo.Context, id string) (Token, error) {
	var token Token
	err := Tx(ctx).First(&token, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return Token{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Token not found"}
	}
	return token, err
}

// updateToken updates expires_at and last_used_at on existing token object
// Returns new state of the token object
func updateToken(ctx echo.Context, token Token) (Token, error) {
	if err := token.Validate(); err != nil {
		return token, err
	}

	token.touch()

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

// TokenService is a service for managing API auth tokens
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

// CreateToken creates a new token object.
func (t TokenService) CreateToken(ctx echo.Context, tokenCreate app.TokenCreate) (app.Token, error) {
	token := Token{
		UserID: tokenCreate.UserID,
		AuthID: tokenCreate.AuthID,
	}

	err := token.create(ctx)
	if err != nil {
		return app.Token{}, err
	}

	return convertToken(token), nil
}

func (t TokenService) DeleteToken(ctx echo.Context, id string) error {
	return deleteToken(ctx, id)
}

func (t TokenService) UpdateToken(ctx echo.Context, id string) error {
	token, err := findTokenByID(ctx, id)
	if err != nil {
		return err
	}
	token.touch()
	_, err = updateToken(ctx, token)
	return err
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
