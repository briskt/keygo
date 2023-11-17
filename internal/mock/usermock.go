package mock

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type UserService struct {
	users map[string]app.User

	FindUsersFn func(ctx echo.Context, filter app.UserFilter) ([]app.User, int, error)
}

func NewUserService() UserService {
	return UserService{
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
	if m.FindUsersFn != nil {
		return m.FindUsersFn(context, filter)
	}
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

func (m *UserService) CreateUser(context echo.Context, input app.UserCreateInput) (app.User, error) {
	if err := input.Validate(); err != nil {
		return app.User{}, err
	}
	now := time.Now()
	user := app.User{
		ID:        newID(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		AvatarURL: input.AvatarURL,
		Role:      input.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *UserService) UpdateUser(context echo.Context, id string, input app.UserUpdateInput) (app.User, error) {
	if err := input.Validate(); err != nil {
		return app.User{}, err
	}
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