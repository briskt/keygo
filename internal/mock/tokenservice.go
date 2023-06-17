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
}

func NewTokenService() app.TokenService {
	return &TokenService{}
}

// Init preloads the mock "database" with tokens
func (m *TokenService) Init(fakeTokens []app.Token) {
	m.tokens = make(map[string]app.Token, len(fakeTokens))
	for i := range fakeTokens {
		m.tokens[fakeTokens[i].PlainText] = fakeTokens[i]
	}
}

func (m *TokenService) FindToken(ctx echo.Context, raw string) (app.Token, error) {
	if t, ok := m.tokens[raw]; ok {
		return t, nil
	}
	return app.Token{}, fmt.Errorf("token %s not found", raw)
}

func (m *TokenService) CreateToken(ctx echo.Context, authID string) (app.Token, error) {
	if m.tokens == nil {
		m.tokens = make(map[string]app.Token)
	}
	mockRandomToken := strconv.Itoa(int(time.Now().Unix()))
	newToken := app.Token{
		AuthID:    authID,
		PlainText: mockRandomToken,
		ExpiresAt: time.Now().Add(time.Minute),
	}
	m.tokens[mockRandomToken] = newToken
	return newToken, nil
}

func (m *TokenService) DeleteToken(ctx echo.Context, tokenID string) error {
	panic("implement TokenService DeleteToken")
}
