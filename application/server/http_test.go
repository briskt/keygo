package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
	"github.com/schparky/keygo/migrations"
	"github.com/schparky/keygo/server"
)

// TestSuite contains common setup and configuration for tests
type TestSuite struct {
	suite.Suite
	*require.Assertions
	server *server.Server
	ctx    echo.Context
	DB     *gorm.DB
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	if sqlDB, err := ts.DB.DB(); err != nil {
		panic(err.Error())
	} else {
		migrations.Fresh(sqlDB)
	}
	ts.Assertions = require.New(ts.T())
}

func Test_RunSuite(t *testing.T) {
	suite.Run(t, &TestSuite{
		server: server.New(),
		ctx:    testContext(db.DB),
		DB:     db.DB,
	})
}

func testContext(db *gorm.DB) echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.Set(keygo.ContextKeyTx, db)
	return ctx
}

// NewRequest creates a new HTTP request using the server's base URL and
// attaching a user session based on the context.
//func (ts *TestSuite) NewRequest(method, url string, body io.Reader) *http.Request {
//	r, err := http.NewRequest(method, s.URL()+url, body)
//	if err != nil {
//		tb.Fatal(err)
//	}
//
//	// Generate session cookie for user, if logged in.
//	if user := wtf.UserFromContext(ctx); user != nil {
//		data, err := s.MarshalSession(wtfhttp.Session{UserID: user.ID})
//		if err != nil {
//			tb.Fatal(err)
//		}
//		r.AddCookie(&http.Cookie{
//			Name:  wtfhttp.SessionCookieName,
//			Value: data,
//			Path:  "/",
//		})
//	}
//
//	return r
//}
