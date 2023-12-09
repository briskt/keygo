package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) Test_tenantsCreateHandler() {
	f := ts.createUserFixture()
	user := f.Users[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)

	tests := []struct {
		name       string
		actor      db.User
		wantStatus int
	}{
		{
			name:       "not a valid user",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot create a tenant",
			actor:      user,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can create a tenant",
			actor:      admin,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			input := app.TenantCreateInput{Name: "new tenant"}
			body, status := ts.request(http.MethodPost, "/api/tenants", tt.actor.Email, input)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var gotTenant app.Tenant
			ts.NoError(json.Unmarshal(body, &gotTenant))
			ts.Equal(input.Name, gotTenant.Name, "incorrect Tenant Name, body: \n%s", body)

			dbTenant, err := db.FindTenantByID(ts.ctx, gotTenant.ID)
			ts.NoError(err)
			ts.Equal(input.Name, dbTenant.Name, "incorrect Tenant Name in db")
		})
	}
}

func (ts *TestSuite) Test_tenantsGetHandler() {
	f := ts.createUserFixture()
	user := f.Users[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)

	tenant := ts.createTenantFixture().Tenants[0]

	tests := []struct {
		name       string
		actor      db.User
		wantStatus int
	}{
		{
			name:       "not a valid user",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot access a tenant",
			actor:      user,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can access a tenant",
			actor:      admin,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/tenants/"+tenant.ID, tt.actor.Email, nil)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var gotTenant app.Tenant
			ts.NoError(json.Unmarshal(body, &gotTenant))
			ts.Equal(tenant.ID, gotTenant.ID, "incorrect Tenant data, body: \n%s", body)
		})
	}
}

func (ts *TestSuite) Test_tenantsListHandler() {
	f := ts.createUserFixture()
	user := f.Users[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)

	tenant := ts.createTenantFixture().Tenants[0]

	tests := []struct {
		name       string
		actor      db.User
		wantStatus int
	}{
		{
			name:       "not a valid user",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot access a tenant",
			actor:      user,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can access a tenant",
			actor:      admin,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/tenants", tt.actor.Email, nil)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var tenants []app.Tenant
			ts.NoError(json.Unmarshal(body, &tenants))
			ts.Len(tenants, 1)
			ts.Equal(tenant.ID, tenants[0].ID, "incorrect Tenant ID, body: \n%s", body)
		})
	}
}

func (ts *TestSuite) Test_tenantsUsersCreateHandler() {
	user := ts.createUserFixture().Users[0]

	admin := ts.createUserFixture().Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)

	tenant := ts.createTenantFixture().Tenants[0]

	tests := []struct {
		name       string
		actor      db.User
		wantStatus int
	}{
		{
			name:       "not a valid user",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot access a tenant",
			actor:      user,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can access a tenant",
			actor:      admin,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			input := app.TenantUserCreateInput{Email: "tenant_user@example.com"}
			path := fmt.Sprintf("/api/tenants/%s/users", tenant.ID)
			body, status := ts.request(http.MethodPost, path, tt.actor.Email, input)

			// Assertions
			ts.Equal(tt.wantStatus, status, "incorrect http status, body: \n%s", body)

			if tt.wantStatus != http.StatusOK {
				return
			}

			var user db.User
			ts.NoError(json.Unmarshal(body, &user))
			ts.Equal(input.Email, user.Email, "incorrect user Email, body: \n%s", body)
			ts.Equal(tenant.ID, *user.TenantID, "incorrect user TenantID, body: \n%s", body)
		})
	}
}
