package db_test

import (
	"github.com/google/uuid"

	"github.com/briskt/keygo"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) TestTokenService_CreateToken() {
	s := db.NewTokenService()

	auth := ts.CreateAuth()

	// Create new record and check generated fields
	newToken, err := s.CreateToken(ts.ctx, auth.ID)

	ts.NoError(err)
	ts.False(newToken.ID == uuid.Nil, "ID is not set")
	ts.NotEmpty(newToken.PlainText, "expected Token")
	ts.False(newToken.ExpiresAt.IsZero(), "expected ExpiredAt")
	ts.False(newToken.CreatedAt.IsZero(), "expected CreatedAt")
	ts.False(newToken.UpdatedAt.IsZero(), "expected UpdatedAt")

	// Query database and compare
	fromDB, err := s.FindToken(ts.ctx, newToken.PlainText)
	ts.NoError(err, "couldn't find created token %s", newToken.PlainText)
	fromDB.PlainText = newToken.PlainText
	fromDB.LastLoginAt = newToken.LastLoginAt
	fromDB.ExpiresAt = newToken.ExpiresAt
	ts.SameToken(newToken, fromDB)

	// Expect validation error
	_, err = s.CreateToken(ts.ctx, uuid.UUID{})
	ts.Error(err, "expected validation error")
	ts.Equal(keygo.ERR_INVALID, keygo.ErrorCode(err))
	ts.Equal(`AuthID required.`, keygo.ErrorMessage(err), "unexpected error")
}

// SameToken verifies two Token objects are the same except for the timestamps
func (ts *TestSuite) SameToken(expected keygo.Token, actual keygo.Token, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}
