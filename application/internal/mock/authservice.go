package mock

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo"
)

type AuthService struct {
	auths []keygo.Auth
}

func NewAuthService() keygo.AuthService {
	return &AuthService{}
}

func (m *AuthService) FindAuthByID(context echo.Context, uuid uuid.UUID) (keygo.Auth, error) {
	panic("implement AuthService FindAuthByID")
}

func (m *AuthService) FindAuths(context echo.Context, filter keygo.AuthFilter) ([]keygo.Auth, int, error) {
	panic("implement AuthService FindAuths")
}

func (m *AuthService) CreateAuth(context echo.Context, auth keygo.Auth) (keygo.Auth, error) {
	if m.auths == nil {
		m.auths = make([]keygo.Auth, 0)
	}
	m.auths = append(m.auths, auth)
	return auth, nil
}

func (m *AuthService) DeleteAuth(context echo.Context, uuid uuid.UUID) error {
	panic("implement AuthService DeleteAuth")
}
