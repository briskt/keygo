package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) Test_GetUser() {
	f := ts.createUserFixture()
	user := f.Users[0]
	userToken := f.Tokens[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	adminToken := f2.Tokens[0]

	tests := []struct {
		name       string
		token      string
		userID     string
		wantStatus int
	}{
		{
			name:       "not a valid token",
			token:      "x",
			userID:     user.ID,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "not a valid ID",
			token:      adminToken.PlainText,
			userID:     "x",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "non-admin cannot access other users",
			token:      userToken.PlainText,
			userID:     admin.ID,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "a user can access their own record",
			token:      userToken.PlainText,
			userID:     user.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin can access other users",
			token:      adminToken.PlainText,
			userID:     user.ID,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/users/"+tt.userID, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, "Bearer "+tt.token)

			res := httptest.NewRecorder()
			ts.server.ServeHTTP(res, req)
			body, err := io.ReadAll(res.Body)
			ts.NoError(err)

			// Assertions
			ts.Equal(tt.wantStatus, res.Code, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var gotUser app.User
			ts.NoError(json.Unmarshal(body, &gotUser))
			ts.Equal(user.ID, gotUser.ID, "incorrect user data, body: \n%s", body)
		})
	}
}

func (ts *TestSuite) Test_GetUserList() {
	f := ts.createUserFixture()
	user := f.Users[0]
	token := f.Tokens[0]

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var users []app.User
	ts.NoError(json.Unmarshal(body, &users))
	ts.Len(users, 1)
	ts.Equal(user.ID, users[0].ID, "incorrect user ID, body: \n%s", body)
	ts.Equal(user.Email, users[0].Email, "incorrect user email, body: \n%s", body)
	ts.Equal(user.Role, users[0].Role, "incorrect user role, body: \n%s", body)

	// TODO: test error response
}
