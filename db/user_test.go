package db_test

import (
	"testing"
	"time"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) TestUserService_CreateUser() {
	s := db.NewUserService()

	u := app.UserCreate{
		FirstName: "susy",
		LastName:  "smith",
		Email:     "susy@example.com",
	}

	// Create new record and check generated fields
	newUser, err := s.CreateUser(ts.ctx, u)
	ts.NoError(err)
	ts.False(newUser.ID == "", "ID is not set")
	ts.False(newUser.CreatedAt.IsZero(), "expected CreatedAt")
	ts.False(newUser.UpdatedAt.IsZero(), "expected UpdatedAt")

	// Query database and compare
	fromDB, err := s.FindUserByID(ts.ctx, newUser.ID)
	ts.NoError(err)
	ts.SameUser(newUser, fromDB)

	// Expect a validation error
	_, err = s.CreateUser(ts.ctx, app.UserCreate{})
	ts.Error(err)
	ts.Equal(app.ErrorCode(err), app.ERR_INVALID)
	ts.Equal(`Email required`, app.ErrorMessage(err))
}

// SameUser verifies two User objects are the same except for the timestamps
func (ts *TestSuite) SameUser(expected app.User, actual app.User, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}

// CreateUser creates a user in the database. Fatal on error.
func (ts *TestSuite) CreateUser(user app.UserCreate) app.User {
	ts.T().Helper()
	newUser, err := db.NewUserService().CreateUser(ts.ctx, user)
	if err != nil {
		ts.Fail("failed to create user: " + err.Error())
	}
	return newUser
}

func (ts *TestSuite) Test_FindUserByID() {
	s := db.NewUserService()

	user, err := s.CreateUser(ts.ctx, app.UserCreate{FirstName: "joe", Email: "joe@example.com"})
	ts.NoError(err)

	found, err := s.FindUserByID(ts.ctx, user.ID)
	ts.NoError(err)

	ts.Equal(user.ID, found.ID)
	ts.Equal(user.FirstName, found.FirstName)
	ts.Equal(user.Email, found.Email)
}

func (ts *TestSuite) Test_FindUsers() {
	s := db.NewUserService()

	joeEmail := "joe@example.com"
	joe, err := s.CreateUser(ts.ctx, app.UserCreate{FirstName: "joe", Email: joeEmail})
	ts.NoError(err)
	sally, err := s.CreateUser(ts.ctx, app.UserCreate{FirstName: "sally", Email: "sally@example.com"})
	ts.NoError(err)

	notFindableEmail := "nobody@example.com"

	tests := []struct {
		name      string
		filter    app.UserFilter
		wantError bool
		wantUsers []string
	}{
		{
			name:      "empty filter",
			filter:    app.UserFilter{},
			wantError: false,
			wantUsers: []string{joe.ID, sally.ID},
		},
		{
			name:      "filter by email",
			filter:    app.UserFilter{Email: &joeEmail},
			wantError: false,
			wantUsers: []string{joe.ID},
		},
		{
			name:      "no results",
			filter:    app.UserFilter{Email: &notFindableEmail},
			wantError: false,
			wantUsers: []string{},
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			found, n, err := s.FindUsers(ts.ctx, tt.filter)
			if tt.wantError {
				ts.Error(err)
				return
			}
			ts.NoError(err)
			ts.Equal(n, len(found))
			foundIDs := make([]string, len(found))
			for i := range found {
				foundIDs[i] = found[i].ID
			}
			ts.EqualValues(tt.wantUsers, foundIDs)
		})
	}
}

func (ts *TestSuite) Test_TouchLastLoginAt() {
	s := db.NewUserService()

	joe, err := s.CreateUser(ts.ctx, app.UserCreate{FirstName: "joe", Email: "joe@example.com"})
	ts.NoError(err)

	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	err = ts.DB.Exec("update users set last_login_at = ? where id = ?", yesterday, joe.ID).Error
	ts.NoError(err)

	err = s.TouchLastLoginAt(ts.ctx, joe.ID)
	ts.NoError(err)

	found, err := s.FindUserByID(ts.ctx, joe.ID)
	ts.NoError(err)

	ts.WithinDuration(now, *found.LastLoginAt, time.Second)
}
