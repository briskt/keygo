package server_test

import (
	"encoding/json"
	"net/http"
	"testing"

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
			body, status := ts.request(http.MethodGet, "/api/users/"+tt.userID, tt.token, nil)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

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
	userToken := f.Tokens[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	adminToken := f2.Tokens[0]

	tests := []struct {
		name       string
		token      string
		wantStatus int
		want       int
	}{
		{
			name:       "not a valid token",
			token:      "x",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "non-admin gets an empty list of users",
			token:      userToken.PlainText,
			wantStatus: http.StatusOK,
			want:       0,
		},
		{
			name:       "admin can list users",
			token:      adminToken.PlainText,
			wantStatus: http.StatusOK,
			want:       2,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/users", tt.token, nil)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var users []app.User
			ts.NoError(json.Unmarshal(body, &users))
			ts.Equal(tt.want, len(users), "got the wrong number of users")
		})
	}
}
