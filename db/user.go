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
	TenantID    *string
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deleted     gorm.DeletedAt
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.ID = newID()
	return nil
}

// FindUserByID retrieves a user by ID along with their associated auth objects.
func FindUserByID(ctx echo.Context, id string) (User, error) {
	user, err := findUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func FindUsers(ctx echo.Context, filter app.UserFilter) ([]User, error) {
	var users []User
	q := Tx(ctx)
	if filter.Email != nil {
		q = q.Where("email = ?", filter.Email)
	}
	if filter.TenantID != nil {
		q = q.Where("tenant_id = ?", filter.TenantID)
	}
	result := q.Find(&users)
	return users, result.Error
}

// CreateUser creates a new user.
func CreateUser(ctx echo.Context, userCreate app.UserCreateInput) (User, error) {
	if err := userCreate.Validate(); err != nil {
		return User{}, err
	}

	user := User{
		FirstName: userCreate.FirstName,
		LastName:  userCreate.LastName,
		Email:     userCreate.Email,
		AvatarURL: userCreate.AvatarURL,
		Role:      userCreate.Role,
	}

	// TODO: remove this when ready
	user.Role = app.UserRoleAdmin

	result := Tx(ctx).Create(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// UpdateUser updates a user object.
func UpdateUser(ctx echo.Context, id string, input app.UserUpdateInput) (User, error) {
	if err := input.Validate(); err != nil {
		return User{}, err
	}
	user, err := findUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	result := Tx(ctx).Save(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// DeleteUser permanently deletes a user and all child objects
func DeleteUser(ctx echo.Context, id string) error {
	result := Tx(ctx).Where("id = ?", id).Delete(&User{})
	return result.Error
}

// TouchLastLoginAt sets the LastLoginAt field to the current time
func TouchLastLoginAt(ctx echo.Context, id string) error {
	result := Tx(ctx).Model(&User{}).Where("id = ?", id).Update("last_login_at", time.Now())
	return result.Error
}

// findUserByID is a helper function to fetch a user by ID.
func findUserByID(ctx echo.Context, id string) (User, error) {
	var user User
	result := Tx(ctx).First(&user, "id = ?", id)
	return user, result.Error
}

func ConvertUser(_ echo.Context, u User) (app.User, error) {
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
