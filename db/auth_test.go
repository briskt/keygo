package db_test

import (
	"testing"
	"time"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) TestAuthService_CreateAuth() {
	s := db.NewAuthService()
	now := time.Now()

	newUser := app.User{
		FirstName: "Joe",
		Email:     "joe@example.com",
	}

	existingUser := ts.CreateUser(app.User{FirstName: "Clark", Email: "clark@example.com"})

	tests := []struct {
		name        string
		auth        app.Auth
		wantErr     string
		wantErrCode string
	}{
		{
			name: "validation error",
			auth: app.Auth{
				User: newUser,
			},
			wantErr:     "provider required",
			wantErrCode: app.ERR_INVALID,
		},
		{
			name: "new user",
			auth: app.Auth{
				Provider:   "provider1",
				ProviderID: "xyz1234",
				User:       newUser,
			},
		},
		{
			name: "existing user",
			auth: app.Auth{
				Provider:   "provider2",
				ProviderID: "1",
				User:       existingUser,
			},
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			// Create new record and check generated fields
			newAuth, err := s.CreateAuth(ts.ctx, tt.auth)

			if tt.wantErr != "" {
				ts.Error(err, "didn't get expected error")
				ts.Equal(tt.wantErr, app.ErrorMessage(err), "unexpected error")
				ts.Equal(tt.wantErrCode, app.ErrorCode(err), "unexpected error code")
				return
			}

			ts.NoError(err)
			ts.False(newAuth.ID == "", "ID is not set")
			ts.WithinDuration(now, newAuth.CreatedAt, time.Second, "CreatedAt is not set")
			ts.WithinDuration(now, newAuth.UpdatedAt, time.Second, "UpdatedAt is not set")
			ts.False(newAuth.UserID == "", "UserID is not set")
			ts.Equal(newAuth.User.ID, newAuth.UserID, "User.ID is not set")

			// Query database and compare
			fromDB, err := s.FindAuthByID(ts.ctx, newAuth.ID)
			ts.NoError(err)
			ts.SameAuth(newAuth, fromDB)

			userDB, err := db.NewUserService().FindUserByID(ts.ctx, newAuth.UserID)
			ts.NoError(err, "user not found")
			ts.Equal(tt.auth.User.Email, userDB.Email, "user email not correct")
		})
	}
}

// SameAuth verifies two Auth objects are the same except for the timestamps
func (ts *TestSuite) SameAuth(expected app.Auth, actual app.Auth, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}

// CreateAuth creates an auth in the database. Fatal on error.
func (ts *TestSuite) CreateAuth() app.Auth {
	ts.T().Helper()

	auth := app.Auth{Provider: "a", ProviderID: "a", User: app.User{FirstName: "a", Email: "a"}}
	newAuth, err := db.NewAuthService().CreateAuth(ts.ctx, auth)
	if err != nil {
		ts.Fail("failed to create auth: " + err.Error())
	}
	return newAuth
}
