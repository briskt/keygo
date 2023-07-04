package db_test

import (
	"github.com/briskt/keygo/app"
)

func (ts *TestSuite) TestTokenService_CreateToken() {
	user, err := ts.UserService.CreateUser(ts.ctx, app.UserCreate{Email: "a@b.com"})
	ts.NoError(err)

	// Create new record and check generated fields
	newToken, err := ts.TokenService.CreateToken(ts.ctx, app.TokenCreate{AuthID: "a", UserID: user.ID})

	ts.NoError(err)
	ts.NotEmpty(newToken.ID, "ID is not set")
	ts.NotEmpty(newToken.PlainText, "expected Token")
	ts.NotZero(newToken.ExpiresAt, "expected ExpiredAt")
	ts.NotZero(newToken.CreatedAt, "expected CreatedAt")
	ts.NotZero(newToken.UpdatedAt, "expected UpdatedAt")

	// Query database and compare
	fromDB, err := ts.TokenService.FindToken(ts.ctx, newToken.PlainText)
	ts.NoError(err, "couldn't find created token %s", newToken.PlainText)
	ts.Equal(newToken.ID, fromDB.ID)
	ts.Equal(newToken.UserID, fromDB.UserID)

	// Expect validation error
	_, err = ts.TokenService.CreateToken(ts.ctx, app.TokenCreate{})
	ts.Error(err, "expected validation error")
	ts.Equal(app.ERR_INVALID, app.ErrorCode(err))
}
