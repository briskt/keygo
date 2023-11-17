package db

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/briskt/keygo/app"
)

type User struct {
	ID          string `gorm:"primaryKey;type:string"`
	FirstName   string
	LastName    string
	Email       string
	AvatarURL   string
	Role        string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deleted     gorm.DeletedAt
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.ID = newID()
	return nil
}

// Ensure service implements interface.
var _ app.UserService = (*UserService)(nil)

// UserService is a service for managing users.
type UserService struct{}

// NewUserService returns a new instance of UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// FindUserByID retrieves a user by ID along with their associated auth objects.
func (s *UserService) FindUserByID(ctx echo.Context, id string) (app.User, error) {
	user, err := findUserByID(ctx, id)
	if err != nil {
		return app.User{}, err
	}
	return convertUser(ctx, user)
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(ctx echo.Context, filter app.UserFilter) ([]app.User, int, error) {
	users, n, err := findUsers(ctx, filter)
	if err != nil {
		return []app.User{}, 0, err
	}
	keygoUsers := make([]app.User, len(users))
	for i := range users {
		u, err := convertUser(ctx, users[i])
		if err != nil {
			return nil, 0, err
		}
		keygoUsers[i] = u
	}
	return keygoUsers, n, nil
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx echo.Context, userCreate app.UserCreateInput) (app.User, error) {
	if err := userCreate.Validate(); err != nil {
		return app.User{}, err
	}
	newUser, err := createUser(ctx, User{
		FirstName: userCreate.FirstName,
		LastName:  userCreate.LastName,
		Email:     userCreate.Email,
		AvatarURL: userCreate.AvatarURL,
		Role:      userCreate.Role,
	})
	if err != nil {
		return app.User{}, err
	}
	return convertUser(ctx, newUser)
}

// UpdateUser updates a user object.
func (s *UserService) UpdateUser(ctx echo.Context, id string, input app.UserUpdateInput) (app.User, error) {
	if err := input.Validate(); err != nil {
		return app.User{}, err
	}
	user, err := updateUser(ctx, id, input)
	if err != nil {
		return app.User{}, err
	}
	return convertUser(ctx, user)
}

// DeleteUser permanently deletes a user and all child objects
func (s *UserService) DeleteUser(ctx echo.Context, id string) error {
	if err := deleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

// TouchLastLoginAt sets the LastLoginAt field to the current time
func (s *UserService) TouchLastLoginAt(ctx echo.Context, id string) error {
	result := Tx(ctx).Model(&User{}).Where("id = ?", id).Update("last_login_at", time.Now())
	return result.Error
}

// findUserByID is a helper function to fetch a user by ID.
func findUserByID(ctx echo.Context, id string) (User, error) {
	var user User
	result := Tx(ctx).First(&user, "id = ?", id)
	return user, result.Error
}

// findUsers returns a list of users. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsers(ctx echo.Context, filter app.UserFilter) ([]User, int, error) {
	// TODO: implement (or remove) other filter parameters
	var users []User
	q := Tx(ctx)
	if filter.Email != nil {
		q = q.Where("email = ?", filter.Email)
	}
	result := q.Find(&users)
	return users, len(users), result.Error
}

// createUser creates a new user. Sets the new database ID to user.ID and sets
// the timestamps to the current time.
func createUser(ctx echo.Context, user User) (User, error) {
	// TODO: remove this when ready
	user.Role = "Admin"
	result := Tx(ctx).Create(&user)
	return user, result.Error
}

// updateUser updates fields on a user object.
func updateUser(ctx echo.Context, id string, upd app.UserUpdateInput) (User, error) {
	user, err := findUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	if upd.Email != nil {
		user.Email = *upd.Email
	}
	if upd.FirstName != nil {
		user.FirstName = *upd.FirstName
	}
	if upd.LastName != nil {
		user.LastName = *upd.LastName
	}

	result := Tx(ctx).Save(&user)
	return user, result.Error
}

// deleteUser permanently removes a user by ID.
func deleteUser(ctx echo.Context, id string) error {
	result := Tx(ctx).Where("id = ?", id).Delete(&User{})
	return result.Error
}

func convertUser(_ echo.Context, u User) (app.User, error) {
	return app.User{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		Role:        u.Role,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}, nil
}
