package db_test

import (
	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) TestTokenService_CreateToken() {
	s := db.NewTokenService()

	token := app.Token{AuthID: "a", User: app.User{FirstName: "a", Email: "a"}}

	// Create new record and check generated fields
	newToken, err := s.CreateToken(ts.ctx, token)

	ts.NoError(err)
	ts.False(newToken.ID == "", "ID is not set")
	ts.NotEmpty(newToken.PlainText, "expected Token")
	ts.False(newToken.ExpiresAt.IsZero(), "expected ExpiredAt")
	ts.False(newToken.CreatedAt.IsZero(), "expected CreatedAt")
	ts.False(newToken.UpdatedAt.IsZero(), "expected UpdatedAt")

	// Query database and compare
	fromDB, err := s.FindToken(ts.ctx, newToken.PlainText)
	ts.NoError(err, "couldn't find created token %s", newToken.PlainText)
	fromDB.PlainText = newToken.PlainText
	fromDB.LastUsedAt = newToken.LastUsedAt
	fromDB.ExpiresAt = newToken.ExpiresAt
	ts.SameToken(newToken, fromDB)

	// Expect validation error
	_, err = s.CreateToken(ts.ctx, app.Token{})
	ts.Error(err, "expected validation error")
	ts.Equal(app.ERR_INVALID, app.ErrorCode(err))
	ts.Equal(`AuthID required.`, app.ErrorMessage(err), "unexpected error")
}

// SameToken verifies two Token objects are the same except for the timestamps
func (ts *TestSuite) SameToken(expected app.Token, actual app.Token, msgAndArgs ...interface{}) {
	actual.CreatedAt = expected.CreatedAt
	actual.UpdatedAt = expected.UpdatedAt
	ts.Equal(expected, actual, msgAndArgs...)
}
