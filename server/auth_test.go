package server_test

import (
	"testing"

	"github.com/briskt/keygo/app"
)

func (ts *TestSuite) TestServer_findOrCreateUser() {
	user := ts.createUserFixture(app.UserRoleBasic)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing user",
			email:   user.Email,
			wantErr: false,
		},
		{
			name:    "new user",
			email:   "joe@example.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			got, err := ts.server.FindOrCreateUser(ts.ctx, tt.email)
			ts.NoError(err)
			ts.Equal(tt.email, got.Email)
		})
	}
}
