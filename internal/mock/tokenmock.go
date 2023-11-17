package mock

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type TokenService struct {
	tokens map[string]app.Token // key is timestamp

	FindTokenFn   func(ctx echo.Context, raw string) (app.Token, error)
	UpdateTokenFn func(ctx echo.Context, id string, input app.TokenUpdateInput) error
}

func NewTokenService() TokenService {
	return TokenService{}
}

func (m *TokenService) DeleteAllTokens() {
	m.tokens = map[string]app.Token{}
}

// Init preloads the mock "database" with tokens
func (m *TokenService) Init(fakeTokens []app.Token) {
	m.tokens = make(map[string]app.Token, len(fakeTokens))
	for i := range fakeTokens {
		m.tokens[fakeTokens[i].PlainText] = fakeTokens[i]
	}
}

func (m *TokenService) FindToken(ctx echo.Context, raw string) (app.Token, error) {
	if m.FindTokenFn != nil {
		return m.FindTokenFn(ctx, raw)
	}
	if t, ok := m.tokens[raw]; ok {
		return t, nil
	}
	return app.Token{}, fmt.Errorf("token %s not found", raw)
}

func (m *TokenService) CreateToken(ctx echo.Context, input app.TokenCreateInput) (app.Token, error) {
	if err := input.Validate(); err != nil {
		return app.Token{}, err
	}
	if m.tokens == nil {
		m.tokens = make(map[string]app.Token)
	}
	mockRandomToken := strconv.Itoa(int(time.Now().Unix()))
	newToken := app.Token{
		AuthID:    input.AuthID,
		UserID:    input.UserID,
		User:      app.User{ID: input.UserID},
		PlainText: mockRandomToken,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.tokens[mockRandomToken] = newToken
	return newToken, nil
}

func (m *TokenService) DeleteToken(ctx echo.Context, tokenID string) error {
	panic("implement TokenService DeleteToken")
}

func (m *TokenService) UpdateToken(ctx echo.Context, id string, input app.TokenUpdateInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	if m.UpdateTokenFn != nil {
		return m.UpdateTokenFn(ctx, id, input)
	}

	t, ok := m.tokens[id]
	if !ok {
		return &app.Error{Code: app.ERR_NOTFOUND, Message: "Token not found"}
	}

	if input.ExpiresAt != nil {
		t.ExpiresAt = *input.ExpiresAt
	}
	if input.LastUsedAt != nil {
		t.LastUsedAt = input.LastUsedAt
	}

	now := time.Now()
	t.UpdatedAt = now
	m.tokens[id] = t
	return nil
}
