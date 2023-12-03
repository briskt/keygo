package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
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

	deleteAll(ts.ctx, &db.Tenant{})
	deleteAll(ts.ctx, &db.Token{})
	deleteAll(ts.ctx, &db.User{})
}

func Test_RunSuite(t *testing.T) {
	db := db.OpenDB()
	svr := server.New(server.WithDataBase(db))
	ctx := testContext()
	ctx.Set(app.ContextKeyTx, db)
	suite.Run(t, &TestSuite{
		server: svr,
		ctx:    ctx,
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
	}
	createdUser, err := db.CreateUser(ts.ctx, fakeUserCreate)
	ts.NoError(err)

	newToken, err := db.CreateToken(ts.ctx, app.TokenCreateInput{
		UserID:    createdUser.ID,
		AuthID:    createdUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	ts.NoError(err)
	if err != nil {
		return Fixtures{}
	}

	return Fixtures{
		Users:  []db.User{createdUser},
		Tokens: []db.Token{newToken},
	}
}

func (ts *TestSuite) createTenantFixture() Fixtures {
	fakeTenantCreate := app.TenantCreateInput{
		Name: "Test Tenant",
	}
	createdTenant, err := db.CreateTenant(ts.ctx, fakeTenantCreate)
	ts.NoError(err)

	return Fixtures{
		Tenants: []db.Tenant{createdTenant},
	}
}

type Fixtures struct {
	Tenants []db.Tenant
	Tokens  []db.Token
	Users   []db.User
}

func deleteAll(c echo.Context, i any) {
	result := db.Tx(c).Where("TRUE").Delete(i)
	if result.Error != nil {
		panic(fmt.Sprintf("failed to delete all %T: %s", i, result.Error))
	}
}
