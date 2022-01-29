package db_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

func (ts *TestSuite) TestUserService_CreateUser() {
	ts.T().Run("OK", func(t *testing.T) {
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
		ts.False(newUser.UpdatedAt.IsZero(), "expected CreatedAt")

		// Query database and compare
		fromDB, err := s.FindUserByID(ts.ctx, newUser.ID)
		ts.NoError(err)
		ts.SameUser(newUser, fromDB)
	})

	// Ensure an error is returned if user's name is not set.
	ts.T().Run("ErrNameRequired", func(t *testing.T) {
		ctx := testContext(ts.DB)

		s := db.NewUserService()
		if _, err := s.CreateUser(ctx, keygo.User{}); err == nil {
			t.Fatal("expected error")
		} else if keygo.ErrorCode(err) != keygo.ERR_INVALID || keygo.ErrorMessage(err) != `FirstName required.` {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

// SameUser verifies two User objects are the same except for the timestamps
func (ts *TestSuite) SameUser(expected keygo.User, actual keygo.User, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}
