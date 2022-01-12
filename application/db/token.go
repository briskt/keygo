package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
)

const (
	tokenLifetime = time.Hour
	tokenBytes    = 32
)

type Token struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`

	Auth   Auth
	AuthID uuid.UUID

	Hash  string
	Token string `gorm:"-"`

	LastLoginAt time.Time
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}

// Validate returns an error if any fields are invalid on the Token object.
func (t *Token) Validate() error {
	if t.AuthID == uuid.Nil {
		return keygo.Errorf(keygo.ERR_INVALID, "AuthID required.")
	} else if t.Hash == "" {
		return keygo.Errorf(keygo.ERR_INVALID, "Hash required.")
	}
	return nil
}

// create a new token object in the database. On success, the ID is set to the new database
// ID & timestamp fields are set to the current time
func (t *Token) create(ctx echo.Context, clientID string) error {
	t.ExpiresAt = time.Now().Add(tokenLifetime)
	t.LastLoginAt = time.Now()
	t.Token = getRandomToken()
	t.Hash = hashToken(clientID + t.Token)

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

// findTokenByID is a helper function to return a token object by ID
// Returns ERR_NOTFOUND if record doesn't exist
func findTokenByID(ctx echo.Context, id uuid.UUID) (Token, error) {
	var token Token
	result := Tx(ctx).Find(&token, id)
	if result.Error == sql.ErrNoRows {
		return Token{}, &keygo.Error{Code: keygo.ERR_NOTFOUND, Message: "Token not found"}
	}
	return token, result.Error
}

// deleteToken permanently removes a token object by ID
func deleteToken(ctx echo.Context, id uuid.UUID) error {
	// Verify object exists & that the user is the owner of the token
	//if token, err := findTokenByID(tx, id); err != nil {
	//	return err
	//} else if token.UserID != keygo.UserIDFromContext(ctx) {
	//	return keygo.Errorf(keygo.ERR_UNAUTHORIZED, "You are not allowed to delete this token")
	//}

	token := keygo.Token{ID: id}
	result := Tx(ctx).Delete(&token)
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
var _ keygo.TokenService = (*TokenService)(nil)

// TokenService represents a service for managing API auth tokens
type TokenService struct{}

// NewTokenService returns a new instance of TokenService
func NewTokenService() *TokenService {
	return &TokenService{}
}

func (t TokenService) FindTokenByID(ctx echo.Context, id uuid.UUID) (keygo.Token, error) {
	token, err := findTokenByID(ctx, id)
	if err != nil {
		return keygo.Token{}, err
	} else if err = token.loadAuth(ctx); err != nil {
		return keygo.Token{}, err
	}

	return convertToken(token), nil
}

func (t TokenService) CreateToken(ctx echo.Context, authID uuid.UUID, clientID string) (keygo.Token, error) {
	token := Token{
		AuthID: authID,
	}
	if err := token.create(ctx, clientID); err != nil {
		return keygo.Token{}, err
	}

	if err := token.loadAuth(ctx); err != nil {
		return keygo.Token{}, err
	}
	return convertToken(token), nil
}

func (t TokenService) DeleteToken(ctx echo.Context, id uuid.UUID) error {
	return deleteToken(ctx, id)
}

func convertToken(token Token) keygo.Token {
	return keygo.Token{
		ID:          token.ID,
		Auth:        convertAuth(token.Auth),
		AuthID:      token.AuthID,
		Token:       token.Token,
		LastLoginAt: token.LastLoginAt,
		ExpiresAt:   token.ExpiresAt,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}
}
