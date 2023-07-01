package mock

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

// tokenLife is the amount of time a new token can be used
// TODO: since token renewal and expiration is duplicated here and in the db package, it maybe indicates this logic should be moved elsewhere.
const tokenLife = time.Minute

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

func (m *TokenService) CreateToken(ctx echo.Context, token app.Token) (app.Token, error) {
	if m.tokens == nil {
		m.tokens = make(map[string]app.Token)
	}
	mockRandomToken := strconv.Itoa(int(time.Now().Unix()))
	newToken := app.Token{
		AuthID:    token.AuthID,
		PlainText: mockRandomToken,
		ExpiresAt: time.Now().Add(tokenLife),
	}
	m.tokens[mockRandomToken] = newToken
	return newToken, nil
}

func (m *TokenService) DeleteToken(ctx echo.Context, tokenID string) error {
	panic("implement TokenService DeleteToken")
}

func (m *TokenService) UpdateToken(ctx echo.Context, id string) error {
	t := m.tokens[id]
	now := time.Now()
	t.UpdatedAt = now
	t.ExpiresAt = now.Add(tokenLife)
	t.LastUsedAt = &now
	m.tokens[id] = t
	return nil
}
