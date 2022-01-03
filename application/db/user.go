package db

import (
	"time"

	"github.com/gofrs/uuid"
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
type UserService struct {
	db *gorm.DB
}

// NewUserService returns a new instance of UserService.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// FindUserByID retrieves a user by ID along with their associated auth objects.
func (s *UserService) FindUserByID(id uuid.UUID) (keygo.User, error) {
	return findUserByID(s.db, id)
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(filter keygo.UserFilter) ([]keygo.User, int, error) {
	return findUsers(s.db, filter)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(user keygo.User) error {
	_, err := createUser(s.db, user)
	return err
}

// UpdateUser updates a user object.
func (s *UserService) UpdateUser(id uuid.UUID, upd keygo.UserUpdate) (keygo.User, error) {
	user, err := updateUser(s.db, id, upd)
	return user, err
}

// DeleteUser permanently deletes a user and all child objects
func (s *UserService) DeleteUser(id uuid.UUID) error {
	if err := deleteUser(s.db, id); err != nil {
		return err
	}
	return nil
}

// findUserByID is a helper function to fetch a user by ID.
func findUserByID(tx *gorm.DB, id uuid.UUID) (keygo.User, error) {
	var user keygo.User
	result := tx.First(&user, id)
	return user, result.Error
}

// findUserByEmail is a helper function to fetch a user by email.
func findUserByEmail(tx *gorm.DB, email string) (keygo.User, error) {
	var user keygo.User
	result := tx.Where("email = ?", email).First(&user)
	return user, result.Error
}

// findUsers returns a list of users. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsers(tx *gorm.DB, filter keygo.UserFilter) ([]keygo.User, int, error) {
	var users []keygo.User
	result := tx.Find(&users)
	return users, len(users), result.Error
}

// createUser creates a new user. Sets the new database ID to user.ID and sets
// the timestamps to the current time.
func createUser(tx *gorm.DB, user keygo.User) (keygo.User, error) {
	result := tx.Create(&user)
	return user, result.Error
}

// updateUser updates fields on a user object.
func updateUser(tx *gorm.DB, id uuid.UUID, upd keygo.UserUpdate) (keygo.User, error) {
	user := keygo.User{
		ID:        id,
		Email:     *upd.Email,
		FirstName: *upd.FirstName,
	}
	result := tx.Save(&user)
	return user, result.Error
}

// deleteUser permanently removes a user by ID.
func deleteUser(tx *gorm.DB, id uuid.UUID) error {
	user := keygo.User{ID: id}
	result := tx.Delete(&user)
	return result.Error
}
