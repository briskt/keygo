package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/internal/mock"
)

func (ts *TestSuite) Test_GetUser() {
	fakeToken := app.Token{
		ID: "1",
		Auth: app.Auth{
			User: app.User{
				ID:    "1",
				Email: "test@example.com",
			},
		},
		PlainText: "12345",
		ExpiresAt: time.Now().Add(time.Minute),
	}
	ts.server.TokenService.(*mock.TokenService).Init([]app.Token{fakeToken})

	req := httptest.NewRequest(http.MethodGet, "/api/users/"+fakeToken.Auth.User.ID, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+fakeToken.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := ioutil.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var user app.User
	ts.NoError(json.Unmarshal(body, &user))
	ts.Equal(fakeToken.Auth.User, user, "incorrect user data, body: \n%s", body)
}

func (ts *TestSuite) Test_GetUserList() {
	fakeToken := app.Token{
		ID: "1",
		Auth: app.Auth{
			User: app.User{
				Email: "test@example.com",
				Role:  "Admin",
			},
		},
		PlainText: "12345",
		ExpiresAt: time.Now().Add(time.Minute),
	}
	ts.server.TokenService.(*mock.TokenService).Init([]app.Token{fakeToken})
	ts.server.UserService.CreateUser(ts.ctx, fakeToken.Auth.User)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+fakeToken.PlainText)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := ioutil.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusOK, res.Code, "incorrect http status, body: \n%s", body)

	var users []app.User
	ts.NoError(json.Unmarshal(body, &users))
	ts.Equal([]app.User{fakeToken.Auth.User}, users, "incorrect user data, body: \n%s", body)
}
