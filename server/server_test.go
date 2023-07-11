package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/internal/mock"
	"github.com/briskt/keygo/server"
)

// TestSuite contains common setup and configuration for tests
type TestSuite struct {
	suite.Suite
	*require.Assertions

	server          *server.Server
	ctx             echo.Context
	mockUserService *mock.UserService
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())
	ts.server.UserService.(*mock.UserService).DeleteAllUsers()
	ts.server.TokenService.(*mock.TokenService).DeleteAllTokens()
}

func Test_RunSuite(t *testing.T) {
	mus := mock.NewUserService()
	s := app.DataServices{
		TenantService: nil,
		TokenService:  mock.NewTokenService(),
		UserService:   &mus,
	}
	svr := server.New(server.WithDataServices(s))
	suite.Run(t, &TestSuite{
		server:          svr,
		ctx:             testContext(),
		mockUserService: &mus,
	})
}

func testContext() echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	return ctx
}

func (ts *TestSuite) createUserFixture() Fixtures {
	fakeUserCreate := app.UserCreate{
		Email: "test@example.com",
		Role:  "Admin",
	}
	createdUser, err := ts.server.UserService.CreateUser(ts.ctx, fakeUserCreate)
	ts.NoError(err)

	fakeToken := app.Token{
		ID: "1",
		User: app.User{
			ID:    createdUser.ID,
			Email: "test@example.com",
			Role:  "Admin",
		},
		PlainText: "12345",
		ExpiresAt: time.Now().Add(time.Minute),
	}
	ts.server.TokenService.(*mock.TokenService).Init([]app.Token{fakeToken})

	return Fixtures{
		Users:  []app.User{createdUser},
		Tokens: []app.Token{fakeToken},
	}
}

type Fixtures struct {
	Users  []app.User
	Tokens []app.Token
}
