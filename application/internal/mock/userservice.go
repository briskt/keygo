package mock

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo"
)

type UserService struct{}

func NewUserService() keygo.UserService {
	return &UserService{}
}

func (m *UserService) FindUserByID(context echo.Context, uuid uuid.UUID) (keygo.User, error) {
	panic("implement UserService FindUserByID")
}

func (m *UserService) FindUsers(context echo.Context, filter keygo.UserFilter) ([]keygo.User, int, error) {
	panic("implement UserService FindUsers")
}

func (m *UserService) CreateUser(context echo.Context, user keygo.User) (keygo.User, error) {
	panic("implement UserService CreateUser")
}

func (m *UserService) UpdateUser(context echo.Context, uuid uuid.UUID, update keygo.UserUpdate) (keygo.User, error) {
	panic("implement UserService UpdateUser")
}

func (m *UserService) DeleteUser(context echo.Context, uuid uuid.UUID) error {
	panic("implement UserService DeleteUser")
}
