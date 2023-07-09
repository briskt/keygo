package server_test

import (
	"testing"
)

func (ts *TestSuite) TestServer_findOrCreateUser() {
	f := ts.createUserFixture()
	user := f.Users[0]

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
