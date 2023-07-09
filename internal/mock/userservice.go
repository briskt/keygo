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

func (m *UserService) DeleteAllUsers() {
	m.users = map[string]app.User{}
}

func (m *UserService) FindUserByID(context echo.Context, id string) (app.User, error) {
	panic("implement UserService FindUserByID")
}

func (m *UserService) FindUsers(context echo.Context, filter app.UserFilter) ([]app.User, int, error) {
	var users []app.User
	for _, u := range m.users {
		if filter.Email != nil && *filter.Email != u.Email {
			continue
		}
		// TODO: implement (or remove) other filter fields

		users = append(users, u)
	}
	return users, len(users), nil
}

func (m *UserService) CreateUser(context echo.Context, userCreate app.UserCreate) (app.User, error) {
	now := time.Now()
	user := app.User{
		ID:        newID(),
		FirstName: userCreate.FirstName,
		LastName:  userCreate.LastName,
		Email:     userCreate.Email,
		AvatarURL: userCreate.AvatarURL,
		Role:      userCreate.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}
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
