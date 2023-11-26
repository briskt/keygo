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
		Role:  app.UserRoleAdmin,
	}
	createdUser, err := db.CreateUser(ts.ctx, fakeUserCreate)
	ts.NoError(err)

	fakeToken := app.Token{
		ID: "1",
		User: app.User{
			ID:    createdUser.ID,
			Email: "test@example.com",
			Role:  app.UserRoleAdmin,
		},
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	newToken, err := db.CreateToken(ts.ctx, app.TokenCreateInput{
		UserID:    createdUser.ID,
		AuthID:    createdUser.ID,
		ExpiresAt: fakeToken.ExpiresAt,
	})
	ts.NoError(err)
	if err != nil {
		return Fixtures{}
	}
	fakeToken.PlainText = newToken.PlainText

	u, err := db.ConvertUser(ts.ctx, createdUser)
	ts.NoError(err)

	return Fixtures{
		Users:  []app.User{u},
		Tokens: []app.Token{fakeToken},
	}
}

func (ts *TestSuite) createTenantFixture() Fixtures {
	fakeTenantCreate := app.TenantCreateInput{
		Name: "Test Tenant",
	}
	createdTenant, err := db.CreateTenant(ts.ctx, fakeTenantCreate)
	ts.NoError(err)

	t, err := db.ConvertTenant(ts.ctx, createdTenant)
	ts.NoError(err)

	return Fixtures{
		Tenants: []app.Tenant{t},
	}
}

type Fixtures struct {
	Tenants []app.Tenant
	Tokens  []app.Token
	Users   []app.User
}

func deleteAll(c echo.Context, i any) {
	result := db.Tx(c).Where("TRUE").Delete(i)
	if result.Error != nil {
		panic(fmt.Sprintf("failed to delete all %T: %s", i, result.Error))
	}
}
