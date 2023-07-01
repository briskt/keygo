package mock

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type UserService struct {
	users map[string]app.User
}

func NewUserService() app.UserService {
	return &UserService{
		users: map[string]app.User{},
	}
}

func (m *UserService) FindUserByID(context echo.Context, id string) (app.User, error) {
	panic("implement UserService FindUserByID")
}

func (m *UserService) FindUsers(context echo.Context, filter app.UserFilter) ([]app.User, int, error) {
	users := make([]app.User, len(m.users))
	i := 0
	for _, u := range m.users {
		users[i] = u
		i++
	}
	return users, 0, nil
}

func (m *UserService) CreateUser(context echo.Context, user app.User) (app.User, error) {
	user.ID = newID()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	m.users[user.ID] = user
	return user, nil
}

func (m *UserService) UpdateUser(context echo.Context, id string, update app.UserUpdate) (app.User, error) {
	panic("implement UserService UpdateUser")
}

func (m *UserService) DeleteUser(context echo.Context, id string) error {
	panic("implement UserService DeleteUser")
}

func (m *UserService) TouchLastLoginAt(context echo.Context, s string) error {
	now := time.Now()
	u := m.users[s]
	u.LastLoginAt = &now
	m.users[s] = u
	return nil
}
