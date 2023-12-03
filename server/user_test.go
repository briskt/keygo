package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (ts *TestSuite) Test_GetUser() {
	f := ts.createUserFixture()
	user := f.Users[0]
	token := f.Tokens[0]

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+user.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var gotUser app.User
	ts.NoError(json.Unmarshal(body, &gotUser))
	ts.Equal(user.ID, gotUser.ID, "incorrect user data, body: \n%s", body)

	// TODO: test error response
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
