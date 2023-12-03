package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (ts *TestSuite) Test_tenantsCreateHandler() {
	f := ts.createUserFixture()
	admin := f.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	token := f.Tokens[0]

	input := app.TenantCreateInput{Name: "new tenant"}
	j, _ := json.Marshal(&input)
	req := httptest.NewRequest(http.MethodPost, "/api/tenants", bytes.NewReader(j))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var gotTenant app.Tenant
	ts.NoError(json.Unmarshal(body, &gotTenant))
	ts.Equal(input.Name, gotTenant.Name, "incorrect Tenant Name, body: \n%s", body)

	dbTenant, err := db.FindTenantByID(ts.ctx, gotTenant.ID)
	ts.NoError(err)
	ts.Equal(input.Name, dbTenant.Name, "incorrect Tenant Name in db")

	// TODO: test error response
}

func (ts *TestSuite) Test_tenantsGetHandler() {
	f := ts.createUserFixture()
	admin := f.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	token := f.Tokens[0]
	tenant := ts.createTenantFixture().Tenants[0]

	req := httptest.NewRequest(http.MethodGet, "/api/tenants/"+tenant.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var gotTenant app.Tenant
	ts.NoError(json.Unmarshal(body, &gotTenant))
	ts.Equal(tenant.ID, gotTenant.ID, "incorrect Tenant data, body: \n%s", body)

	// TODO: test error response
}

func (ts *TestSuite) Test_tenantsListHandler() {
	f := ts.createUserFixture()
	admin := f.Users[0]
	admin.Role = app.UserRoleAdmin
	ts.NoError(db.Tx(ts.ctx).Save(&admin).Error)
	token := f.Tokens[0]
	tenant := ts.createTenantFixture().Tenants[0]

	req := httptest.NewRequest(http.MethodGet, "/api/tenants", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var tenants []app.Tenant
	ts.NoError(json.Unmarshal(body, &tenants))
	ts.Len(tenants, 1)
	ts.Equal(tenant.ID, tenants[0].ID, "incorrect Tenant ID, body: \n%s", body)

	// TODO: test error response
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
