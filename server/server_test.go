package server_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

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
	tx     *gorm.DB
}

type Fixtures struct {
	Tenants []db.Tenant
	Users   []db.User
}

// SetupTest runs before every test function
func (ts *TestSuite) SetupTest() {
	ts.Assertions = require.New(ts.T())

	ts.NoError(ts.tx.Exec("TRUNCATE TABLE tenants CASCADE").Error)
}

func Test_RunSuite(t *testing.T) {
	tx := db.OpenDB()
	svr := server.New(server.WithDataBase(tx))
	ctx := testContext()
	ctx.Set(app.ContextKeyTx, tx)
	suite.Run(t, &TestSuite{
		server: svr,
		ctx:    ctx,
		tx:     tx,
	})
}

func testContext() echo.Context {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	return ctx
}

func (ts *TestSuite) createUserFixture(role string) db.User {
	fakeUserCreate := app.UserCreateInput{
		Email: fmt.Sprintf("test%s@example.com", RandStr(6)),
		Role:  role,
	}
	createdUser, err := db.CreateUser(ts.ctx, fakeUserCreate)
	ts.NoError(err)

	ts.createTokenFixture(createdUser.Email, createdUser.ID)

	return createdUser
}

func (ts *TestSuite) createTenantFixture() db.Tenant {
	fakeTenantCreate := app.TenantCreateInput{
		Name: "Test Tenant",
	}
	createdTenant, err := db.CreateTenant(ts.ctx, fakeTenantCreate)
	ts.NoError(err)

	return createdTenant
}

// RandStr generates a random string of length `n` containing uppercase, lowercase, and numbers
func RandStr(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}

func (ts *TestSuite) request(method, path, token string, input any) ([]byte, int) {
	var r io.Reader
	if input != nil {
		j, _ := json.Marshal(&input)
		r = bytes.NewReader(j)
	}
	req := httptest.NewRequest(method, path, r)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)

	res := httptest.NewRecorder()
	ts.server.ServeHTTP(res, req)
	body, err := io.ReadAll(res.Body)
	ts.NoError(err)
	return body, res.Code
}

func (ts *TestSuite) createTokenFixture(plainText, userID string) db.Token {
	token := db.Token{
		UserID:    userID,
		Hash:      fmt.Sprintf("%x", sha256.Sum256([]byte(plainText))),
		PlainText: plainText,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	err := ts.tx.Omit("User").Create(&token).Error
	if err != nil {
		panic("failed to create token fixture: " + err.Error())
	}
	return token
}
