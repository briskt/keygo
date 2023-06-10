package db_test

import (
	"github.com/google/uuid"

	"github.com/briskt/keygo"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) TestUserService_CreateUser() {
	s := db.NewUserService()

	u := keygo.User{
		FirstName: "susy",
		LastName:  "smith",
		Email:     "susy@example.com",
	}

	// Create new record and check generated fields
	newUser, err := s.CreateUser(ts.ctx, u)
	ts.NoError(err)
	ts.False(newUser.ID == uuid.Nil, "ID is not set")
	ts.False(newUser.CreatedAt.IsZero(), "expected CreatedAt")
	ts.False(newUser.UpdatedAt.IsZero(), "expected UpdatedAt")

	// Query database and compare
	fromDB, err := s.FindUserByID(ts.ctx, newUser.ID)
	ts.NoError(err)
	ts.SameUser(newUser, fromDB)

	// Expect a validation error
	_, err = s.CreateUser(ts.ctx, keygo.User{})
	ts.Error(err)
	ts.Equal(keygo.ErrorCode(err), keygo.ERR_INVALID)
	ts.Equal(`FirstName required`, keygo.ErrorMessage(err))
}

// SameUser verifies two User objects are the same except for the timestamps
func (ts *TestSuite) SameUser(expected keygo.User, actual keygo.User, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}

// CreateUser creates a user in the database. Fatal on error.
func (ts *TestSuite) CreateUser(user keygo.User) keygo.User {
	ts.T().Helper()
	newUser, err := db.NewUserService().CreateUser(ts.ctx, user)
	if err != nil {
		ts.Fail("failed to create user: " + err.Error())
	}
	return newUser
}
