package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	server *server.Server
	ctx    echo.Context
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())
}

func Test_RunSuite(t *testing.T) {
	s := app.DataServices{
		TenantService: nil,
		TokenService:  mock.NewTokenService(),
		UserService:   mock.NewUserService(),
	}
	svr := server.New(server.WithDataServices(s))
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
