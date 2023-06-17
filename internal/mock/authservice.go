package mock

import (
	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type AuthService struct {
	auths []app.Auth
}

func NewAuthService() app.AuthService {
	return &AuthService{}
}

func (m *AuthService) FindAuthByID(context echo.Context, id string) (app.Auth, error) {
	panic("implement AuthService FindAuthByID")
}

func (m *AuthService) FindAuths(context echo.Context, filter app.AuthFilter) ([]app.Auth, int, error) {
	panic("implement AuthService FindAuths")
}

func (m *AuthService) CreateAuth(context echo.Context, auth app.Auth) (app.Auth, error) {
	if m.auths == nil {
		m.auths = make([]app.Auth, 0)
	}
	m.auths = append(m.auths, auth)
	return auth, nil
}

func (m *AuthService) DeleteAuth(context echo.Context, id string) error {
	panic("implement AuthService DeleteAuth")
}
