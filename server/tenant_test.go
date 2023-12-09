package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) Test_tenantsCreateHandler() {
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
	}{
		{
			name:       "not a valid token",
			token:      "x",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot create a tenant",
			token:      userToken.PlainText,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can create a tenant",
			token:      adminToken.PlainText,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			input := app.TenantCreateInput{Name: "new tenant"}
			body, status := ts.request(http.MethodPost, "/api/tenants", tt.token, input)

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
	userToken := f.Tokens[0]

	f2 := ts.createUserFixture()
	admin := f2.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	adminToken := f2.Tokens[0]

	tenant := ts.createTenantFixture().Tenants[0]

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "not a valid token",
			token:      "x",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "a user cannot access a tenant",
			token:      userToken.PlainText,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "admin can access a tenant",
			token:      adminToken.PlainText,
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		ts.T().Run(tt.name, func(t *testing.T) {
			body, status := ts.request(http.MethodGet, "/api/tenants/"+tenant.ID, tt.token, nil)

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
	f := ts.createUserFixture()
	admin := f.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	token := f.Tokens[0]
	tenant := ts.createTenantFixture().Tenants[0]

	input := app.TenantUserCreateInput{Email: "tenant_user@example.com"}
	j, _ := json.Marshal(&input)
	req := httptest.NewRequest(http.MethodPost, "/api/tenants/"+tenant.ID+"/users", bytes.NewReader((j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var user db.User
	ts.NoError(json.Unmarshal(body, &user))
	ts.Equal(input.Email, user.Email, "incorrect user Email, body: \n%s", body)
	ts.Equal(tenant.ID, *user.TenantID, "incorrect user TenantID, body: \n%s", body)

	// TODO: test error response
}
