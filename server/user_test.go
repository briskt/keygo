package server_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) Test_GetUser() {
	user := ts.createUserFixture(app.UserRoleBasic)
	admin := ts.createUserFixture(app.UserRoleAdmin)

	tests := []struct {
		name       string
		actor      db.User
		userID     string
		wantStatus int
	}{
		{
			name:       "not a valid user",
			userID:     user.ID,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "not a valid ID",
			actor:      admin,
			userID:     "x",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "non-admin cannot access other users",
			actor:      user,
			userID:     admin.ID,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "a user can access their own record",
			actor:      user,
			userID:     user.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "admin can access other users",
			actor:      admin,
			userID:     user.ID,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/users/"+tt.userID, tt.actor.Email, nil)

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
	user := ts.createUserFixture(app.UserRoleBasic)
	admin := ts.createUserFixture(app.UserRoleAdmin)

	tests := []struct {
		name       string
		actor      db.User
		wantStatus int
		want       int
	}{
		{
			name:       "not a valid user",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "non-admin gets an empty list of users",
			actor:      user,
			wantStatus: http.StatusOK,
			want:       0,
		},
		{
			name:       "admin can list users",
			actor:      admin,
			wantStatus: http.StatusOK,
			want:       2,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/users", tt.actor.Email, nil)

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
