package db_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

func (ts *TestSuite) TestAuthService_CreateAuth() {
	s := db.NewAuthService()
	now := time.Now()

	newUser := keygo.User{
		FirstName: "Joe",
		Email:     "joe@example.com",
	}

	existingUser := ts.CreateUser(keygo.User{FirstName: "Clark", Email: "clark@example.com"})

	tests := []struct {
		name        string
		auth        keygo.Auth
		wantErr     string
		wantErrCode string
	}{
		{
			name: "validation error",
			auth: keygo.Auth{
				User: newUser,
			},
			wantErr:     "provider required",
			wantErrCode: keygo.ERR_INVALID,
		},
		{
			name: "new user",
			auth: keygo.Auth{
				Provider:   "provider1",
				ProviderID: "xyz1234",
				User:       newUser,
			},
		},
		{
			name: "existing user",
			auth: keygo.Auth{
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
				ts.Equal(tt.wantErr, keygo.ErrorMessage(err), "unexpected error")
				ts.Equal(tt.wantErrCode, keygo.ErrorCode(err), "unexpected error code")
				return
			}

			ts.NoError(err)
			ts.False(newAuth.ID == uuid.Nil, "ID is not set")
			ts.WithinDuration(now, newAuth.CreatedAt, time.Second, "CreatedAt is not set")
			ts.WithinDuration(now, newAuth.UpdatedAt, time.Second, "UpdatedAt is not set")
			ts.False(newAuth.UserID == uuid.Nil, "UserID is not set")
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
func (ts *TestSuite) SameAuth(expected keygo.Auth, actual keygo.Auth, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}

// MustCreateAuth creates an auth in the database. Fatal on error.
func (ts *TestSuite) CreateAuth() keygo.Auth {
	ts.T().Helper()

	auth := keygo.Auth{Provider: "a", ProviderID: "a", User: keygo.User{FirstName: "a", Email: "a"}}
	newAuth, err := db.NewAuthService().CreateAuth(ts.ctx, auth)
	if err != nil {
		ts.Fail("failed to create auth: " + err.Error())
	}
	return newAuth
}
