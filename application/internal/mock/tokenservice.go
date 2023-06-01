package mock

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
)

type TokenService struct {
	tokens map[string]keygo.Token // key is timestamp
}

func NewTokenService() keygo.TokenService {
	return &TokenService{}
}

// Init preloads the mock "database" with tokens
func (m *TokenService) Init(fakeTokens []keygo.Token) {
	m.tokens = make(map[string]keygo.Token, len(fakeTokens))
	for i := range fakeTokens {
		m.tokens[fakeTokens[i].PlainText] = fakeTokens[i]
	}
}

func (m *TokenService) FindToken(ctx echo.Context, raw string) (keygo.Token, error) {
	if t, ok := m.tokens[raw]; ok {
		return t, nil
	}
	return keygo.Token{}, fmt.Errorf("token %s not found", raw)
}

func (m *TokenService) CreateToken(ctx echo.Context, authID uuid.UUID) (keygo.Token, error) {
	if m.tokens == nil {
		m.tokens = make(map[string]keygo.Token)
	}
	mockRandomToken := strconv.Itoa(int(time.Now().Unix()))
	newToken := keygo.Token{
		AuthID:    authID,
		PlainText: mockRandomToken,
		ExpiresAt: time.Now().Add(time.Minute),
	}
	m.tokens[mockRandomToken] = newToken
	return newToken, nil
}

func (m *TokenService) DeleteToken(ctx echo.Context, tokenID uuid.UUID) error {
	panic("implement TokenService DeleteToken")
}
