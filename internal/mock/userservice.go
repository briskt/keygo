package mock

import (
	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type UserService struct{}

func NewUserService() app.UserService {
	return &UserService{}
}

func (m *UserService) FindUserByID(context echo.Context, id string) (app.User, error) {
	panic("implement UserService FindUserByID")
}

func (m *UserService) FindUsers(context echo.Context, filter app.UserFilter) ([]app.User, int, error) {
	panic("implement UserService FindUsers")
}

func (m *UserService) CreateUser(context echo.Context, user app.User) (app.User, error) {
	panic("implement UserService CreateUser")
}

func (m *UserService) UpdateUser(context echo.Context, id string, update app.UserUpdate) (app.User, error) {
	panic("implement UserService UpdateUser")
}

func (m *UserService) DeleteUser(context echo.Context, id string) error {
	panic("implement UserService DeleteUser")
}
