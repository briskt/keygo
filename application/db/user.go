package db

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Role      string    `db:"role"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
	return findUserByID(ctx, id)
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(ctx echo.Context, filter keygo.UserFilter) ([]keygo.User, int, error) {
	return findUsers(ctx, filter)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx echo.Context, user keygo.User) error {
	_, err := createUser(ctx, user)
	return err
}

// UpdateUser updates a user object.
func (s *UserService) UpdateUser(ctx echo.Context, id uuid.UUID, upd keygo.UserUpdate) (keygo.User, error) {
	user, err := updateUser(ctx, id, upd)
	return user, err
}

// DeleteUser permanently deletes a user and all child objects
func (s *UserService) DeleteUser(ctx echo.Context, id uuid.UUID) error {
	if err := deleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}

// findUserByID is a helper function to fetch a user by ID.
func findUserByID(ctx echo.Context, id uuid.UUID) (keygo.User, error) {
	var user keygo.User
	result := Tx(ctx).First(&user, id)
	return user, result.Error
}

// findUserByEmail is a helper function to fetch a user by email.
func findUserByEmail(ctx echo.Context, email string) (keygo.User, error) {
	var user keygo.User
	err := Tx(ctx).Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return keygo.User{}, &keygo.Error{Code: keygo.ERR_NOTFOUND, Message: "User not found"}
	}
	return user, err
}

// findUsers returns a list of users. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsers(ctx echo.Context, filter keygo.UserFilter) ([]keygo.User, int, error) {
	var users []keygo.User
	result := Tx(ctx).Find(&users)
	return users, len(users), result.Error
}

// createUser creates a new user. Sets the new database ID to user.ID and sets
// the timestamps to the current time.
func createUser(ctx echo.Context, user keygo.User) (keygo.User, error) {
	result := Tx(ctx).Create(&user)
	return user, result.Error
}

// updateUser updates fields on a user object.
func updateUser(ctx echo.Context, id uuid.UUID, upd keygo.UserUpdate) (keygo.User, error) {
	user, err := findUserByID(ctx, id)
	if err != nil {
		return keygo.User{}, err
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
	user := keygo.User{ID: id}
	result := Tx(ctx).Delete(&user)
	return result.Error
}
