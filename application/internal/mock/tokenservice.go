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
	tokens map[string]keygo.Token // key is clientID+timestamp
}

func NewTokenService() keygo.TokenService {
	return &TokenService{}
}

func (m *TokenService) FindToken(ctx echo.Context, raw string) (keygo.Token, error) {
	if t, ok := m.tokens[raw]; ok {
		return t, nil
	}
	return keygo.Token{}, fmt.Errorf("token %s not found", raw)
}

func (m *TokenService) CreateToken(ctx echo.Context, authID uuid.UUID, clientID string) (keygo.Token, error) {
	if m.tokens == nil {
		m.tokens = make(map[string]keygo.Token)
	}
	mockRandomToken := strconv.Itoa(int(time.Now().Unix()))
	newToken := keygo.Token{
		AuthID:    authID,
		PlainText: mockRandomToken,
		ExpiresAt: time.Now().Add(time.Minute),
	}
	m.tokens[clientID+mockRandomToken] = newToken
	return newToken, nil
}

func (m *TokenService) DeleteToken(ctx echo.Context, tokenID uuid.UUID) error {
	panic("implement TokenService DeleteToken")
}

func (m *TokenService) Init(t keygo.Token, clientID string) {
	m.tokens = map[string]keygo.Token{clientID + t.PlainText: t}
}
