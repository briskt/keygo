package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/schparky/keygo/internal/mock"
	"github.com/schparky/keygo/server"
)

// TestSuite contains common setup and configuration for tests
type TestSuite struct {
	suite.Suite
	*require.Assertions
	server *server.Server
	ctx    echo.Context
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())
}

func Test_RunSuite(t *testing.T) {
	svr := server.New()
	svr.AuthService = mock.NewAuthService()
	svr.UserService = mock.NewUserService()
	svr.TokenService = mock.NewTokenService()
	suite.Run(t, &TestSuite{
		server: svr,
		ctx:    testContext(),
	})
}

func testContext() echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	return ctx
}
