package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (ts *TestSuite) Test_GetTenant() {
	f := ts.createUserFixture()
	token := f.Tokens[0]
	tenant := ts.createTenantFixture().Tenants[0]

	ts.mockTokenService.FindTokenFn = func(_ echo.Context, raw string) (app.Token, error) {
		return token, nil
	}
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

func (ts *TestSuite) Test_GetTenantList() {
	f := ts.createUserFixture()
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

	var Tenants []app.Tenant
	ts.NoError(json.Unmarshal(body, &Tenants))
	ts.Len(Tenants, 1)
	ts.Equal(tenant.ID, Tenants[0].ID, "incorrect Tenant ID, body: \n%s", body)

	// TODO: test error response
}
