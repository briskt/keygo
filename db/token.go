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

const tokenBytes = 32

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

func (t *Token) BeforeCreate(_ *gorm.DB) error {
	t.ID = newID()
	return nil
}

// create a new token object in the database. On success, the ID is set to the new database
// ID & timestamp fields are set to the current time
func (t *Token) create(ctx echo.Context) error {
	if t.PlainText == "" {
		t.PlainText = randomString()
	}
	t.Hash = hashToken(t.PlainText)

	err := Tx(ctx).Omit("User").Create(t).Error
	return err
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
func updateToken(ctx echo.Context, token Token, input app.TokenUpdateInput) (Token, error) {
	if input.ExpiresAt != nil {
		token.ExpiresAt = *input.ExpiresAt
	}
	if input.LastUsedAt != nil {
		token.LastUsedAt = input.LastUsedAt
	}

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

func FindToken(ctx echo.Context, raw string) (Token, error) {
	token, err := findToken(ctx, raw)
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

// CreateToken creates a new token object.
func CreateToken(ctx echo.Context, input app.TokenCreateInput) (Token, error) {
	if err := input.Validate(); err != nil {
		return Token{}, err
	}

	token := Token{
		UserID:    input.UserID,
		AuthID:    input.AuthID,
		ExpiresAt: input.ExpiresAt,
	}

	err := token.create(ctx)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func DeleteToken(ctx echo.Context, id string) error {
	return deleteToken(ctx, id)
}

func UpdateToken(ctx echo.Context, id string, input app.TokenUpdateInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	token, err := findTokenByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = updateToken(ctx, token, input)
	return err
}

func ConvertToken(ctx echo.Context, token Token) (app.Token, error) {
	if err := token.loadUser(ctx); err != nil {
		return app.Token{}, err
	}

	user, err := ConvertUser(ctx, token.User)
	if err != nil {
		return app.Token{}, err
	}

	return app.Token{
		ID:         token.ID,
		User:       user,
		UserID:     token.UserID,
		AuthID:     token.AuthID,
		PlainText:  token.PlainText,
		LastUsedAt: token.LastUsedAt,
		ExpiresAt:  token.ExpiresAt,
		CreatedAt:  token.CreatedAt,
		UpdatedAt:  token.UpdatedAt,
	}, nil
}
