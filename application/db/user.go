package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName string
	LastName  string
	Email     string
	AvatarURL string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

// Ensure service implements interface.
var _ keygo.UserService = (*UserService)(nil)

// UserService represents a service for managing users.
type UserService struct{}

// NewUserService returns a new instance of UserService.
func NewUserService() *UserService {
	return &UserService{}
}

// FindUserByID retrieves a user by ID along with their associated auth objects.
func (s *UserService) FindUserByID(ctx echo.Context, id uuid.UUID) (keygo.User, error) {
	user, err := findUserByID(ctx, id)
	if err != nil {
		return keygo.User{}, err
	}
	return convertUser(user), nil
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(ctx echo.Context, filter keygo.UserFilter) ([]keygo.User, int, error) {
	users, n, err := findUsers(ctx, filter)
	if err != nil {
		return []keygo.User{}, 0, err
	}
	keygoUsers := make([]keygo.User, len(users))
	for i := range users {
		keygoUsers[i] = convertUser(users[i])
	}
	return keygoUsers, n, nil
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx echo.Context, user keygo.User) (keygo.User, error) {
	if err := user.Validate(); err != nil {
		return keygo.User{}, err
	}
	newUser, err := createUser(ctx, User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	return convertUser(newUser), err
}

// UpdateUser updates a user object.
func (s *UserService) UpdateUser(ctx echo.Context, id uuid.UUID, upd keygo.UserUpdate) (keygo.User, error) {
	user, err := updateUser(ctx, id, upd)
	if err != nil {
		return keygo.User{}, err
	}
	return convertUser(user), nil
}

// DeleteUser permanently deletes a user and all child objects
func (s *UserService) DeleteUser(ctx echo.Context, id uuid.UUID) error {
	if err := deleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

// findUserByID is a helper function to fetch a user by ID.
func findUserByID(ctx echo.Context, id uuid.UUID) (User, error) {
	var user User
	result := Tx(ctx).First(&user, id)
	return user, result.Error
}

// findUserByEmail is a helper function to fetch a user by email.
func findUserByEmail(ctx echo.Context, email string) (User, error) {
	var user User
	err := Tx(ctx).Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, &keygo.Error{Code: keygo.ERR_NOTFOUND, Message: "User not found"}
	}
	return user, err
}

// findUsers returns a list of users. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsers(ctx echo.Context, filter keygo.UserFilter) ([]User, int, error) {
	var users []User
	result := Tx(ctx).Find(&users)
	return users, len(users), result.Error
}

// createUser creates a new user. Sets the new database ID to user.ID and sets
// the timestamps to the current time.
func createUser(ctx echo.Context, user User) (User, error) {
	result := Tx(ctx).Create(&user)
	return user, result.Error
}

// updateUser updates fields on a user object.
func updateUser(ctx echo.Context, id uuid.UUID, upd keygo.UserUpdate) (User, error) {
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
func deleteUser(ctx echo.Context, id uuid.UUID) error {
	result := Tx(ctx).Where("id = ?", id).Delete(&User{})
	return result.Error
}

func convertUser(u User) keygo.User {
	return keygo.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func convertKeygoUser(u keygo.User) User {
	return User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
