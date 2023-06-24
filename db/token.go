package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
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

	Auth   Auth
	AuthID string

	Hash      string
	PlainText string `gorm:"-"`

	LastLoginAt time.Time
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	t.ID = newID()
	return nil
}

// Validate returns an error if any fields are invalid on the Token object.
func (t *Token) Validate() error {
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
	t.LastLoginAt = time.Now()
	t.PlainText = getRandomToken()
	t.Hash = hashToken(t.PlainText)

	if err := t.Validate(); err != nil {
		return err
	}

	err := Tx(ctx).Create(t).Error
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
	if err == sql.ErrNoRows {
		return Token{}, &app.Error{Code: app.ERR_NOTFOUND, Message: "Token not found"}
	}
	return token, err
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

// loadAuth is a helper function to fetch & attach the associated Auth
// to the token object.
func (t *Token) loadAuth(ctx echo.Context) (err error) {
	if t.Auth, err = findAuthByID(ctx, t.AuthID); err != nil {
		return fmt.Errorf("attach token auth: %w", err)
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
	if err = token.loadAuth(ctx); err != nil {
		return app.Token{}, err
	}

	return convertToken(token), nil
}

func (t TokenService) CreateToken(ctx echo.Context, authID string) (app.Token, error) {
	token := Token{
		AuthID: authID,
	}
	if err := token.create(ctx); err != nil {
		return app.Token{}, err
	}

	if err := token.loadAuth(ctx); err != nil {
		return app.Token{}, err
	}
	return convertToken(token), nil
}

func (t TokenService) DeleteToken(ctx echo.Context, id string) error {
	return deleteToken(ctx, id)
}

func convertToken(token Token) app.Token {
	return app.Token{
		ID:          token.ID,
		Auth:        convertAuth(token.Auth),
		AuthID:      token.AuthID,
		PlainText:   token.PlainText,
		LastLoginAt: token.LastLoginAt,
		ExpiresAt:   token.ExpiresAt,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}
}
