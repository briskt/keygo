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

	server            *server.Server
	ctx               echo.Context
	mockTenantService *mock.TenantService // TODO: replace this with server.TenantService
	mockTokenService  *mock.TokenService
	mockUserService   *mock.UserService
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())
	ts.server.TenantService.(*mock.TenantService).DeleteAllTenants()
	ts.server.TokenService.(*mock.TokenService).DeleteAllTokens()
	ts.server.UserService.(*mock.UserService).DeleteAllUsers()
}

func Test_RunSuite(t *testing.T) {
	mockTenantService := mock.NewTenantService()
	mockTokenService := mock.NewTokenService()
	mockUserService := mock.NewUserService()
	mockTokenService.UpdateTokenFn = func(ctx echo.Context, id string, input app.TokenUpdateInput) error {
		return nil
	}

	s := app.DataServices{
		TenantService: &mockTenantService,
		TokenService:  &mockTokenService,
		UserService:   &mockUserService,
	}
	svr := server.New(server.WithDataServices(s))
	suite.Run(t, &TestSuite{
		server:            svr,
		ctx:               testContext(),
		mockTenantService: &mockTenantService,
		mockTokenService:  &mockTokenService,
		mockUserService:   &mockUserService,
	})
}

func testContext() echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	return ctx
}

func (ts *TestSuite) createUserFixture() Fixtures {
	fakeUserCreate := app.UserCreateInput{
		Email: "test@example.com",
		Role:  app.UserRoleAdmin,
	}
	createdUser, err := ts.server.UserService.CreateUser(ts.ctx, fakeUserCreate)
	ts.NoError(err)

	fakeToken := app.Token{
		ID: "1",
		User: app.User{
			ID:    createdUser.ID,
			Email: "test@example.com",
			Role:  app.UserRoleAdmin,
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

func (ts *TestSuite) createTenantFixture() Fixtures {
	fakeTenantCreate := app.TenantCreateInput{
		Name: "Test Tenant",
	}
	createdTenant, err := ts.server.TenantService.CreateTenant(ts.ctx, fakeTenantCreate)
	ts.NoError(err)

	return Fixtures{
		Tenants: []app.Tenant{createdTenant},
	}
}

type Fixtures struct {
	Tenants []app.Tenant
	Tokens  []app.Token
	Users   []app.User
}
