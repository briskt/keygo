package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
	"github.com/schparky/keygo/internal/mock"
)

var (
	mockDB = map[string]*db.User{
		"jon@labstack.com": {ID: uuid.New(), Email: "jon@labstack.com"},
	}
	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com"}`
)

func (ts *TestSuite) Test_GetUser() {
	const clientID = "abc123"
	fakeToken := keygo.Token{
		Auth: keygo.Auth{
			User: keygo.User{
				Email: "test@example.com",
			},
		},
		PlainText: "12345",
		ExpiresAt: time.Now().Add(time.Minute),
	}
	ts.server.TokenService.(*mock.TokenService).Init([]keygo.Token{fakeToken}, []string{clientID})

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+clientID+fakeToken.PlainText)
	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := ioutil.ReadAll(res.Body)
	ts.NoError(err)

	// Assertions
	ts.Equal(http.StatusCreated, res.Code, "incorrect http status, body: \n%s", body)

	var user keygo.User
	ts.NoError(json.Unmarshal(body, &user))
	ts.Equal(fakeToken.Auth.User, user, "incorrect user data, body: \n%s", body)
}
