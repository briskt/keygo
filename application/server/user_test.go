package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

var (
	mockDB = map[string]*db.User{
		"jon@labstack.com": {ID: uuid.New(), Email: "jon@labstack.com"},
	}
	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com"}`
)

func (ts *TestSuite) Test_GetUser() {
	const clientID = "abc123"
	auth := ts.CreateAuth()
	newToken, err := db.NewTokenService().CreateToken(ts.ctx, auth.ID, clientID)
	ts.NoError(err)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+clientID+newToken.PlainText)
	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)

	// Assertions
	ts.Equal(http.StatusCreated, res.Code, "incorrect http status")

	body, err := ioutil.ReadAll(res.Body)
	var user keygo.User
	ts.NoError(json.Unmarshal(body, &user))
	ts.Equal(auth.User, user, "incorrect user data")
}

// CreateAuth creates an auth in the database. Fatal on error.
func (ts *TestSuite) CreateAuth() keygo.Auth {
	ts.T().Helper()

	auth := keygo.Auth{Provider: "a", ProviderID: "a", User: keygo.User{FirstName: "a", Email: "a"}}
	newAuth, err := db.NewAuthService().CreateAuth(ts.ctx, auth)
	if err != nil {
		ts.Fail("failed to create auth: " + err.Error())
	}
	return newAuth
}
