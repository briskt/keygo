package db_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

// TestSuite contains common setup and configuration for tests
type TestSuite struct {
	suite.Suite
	*require.Assertions
	ctx echo.Context
	DB  *gorm.DB
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())
	ts.NoError(ts.DB.Exec("TRUNCATE TABLE tenants CASCADE").Error)
}

func Test_RunSuite(t *testing.T) {
	dbConnection := db.OpenDB()
	suite.Run(t, &TestSuite{
		ctx: testContext(dbConnection),
		DB:  dbConnection,
	})
}

func testContext(tx *gorm.DB) echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.Set(app.ContextKeyTx, tx)
	return ctx
}
