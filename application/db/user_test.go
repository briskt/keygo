package db_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

func TestUserService_CreateUser(t *testing.T) {
	// Ensure user can be created.
	t.Run("OK", func(t *testing.T) {
		tx := MustOpenDB(t)

		s := db.NewUserService()

		u := keygo.User{
			FirstName: "susy",
			LastName:  "smith",
			Email:     "susy@example.com",
		}

		ctx := testContext(tx)

		// Create new user & verify ID and timestamps are set.
		newUser, err := s.CreateUser(ctx, u)
		if err != nil {
			t.Fatal(err)
		} else if newUser.ID == uuid.Nil {
			t.Fatalf("ID is not set: %v", newUser)
		} else if newUser.CreatedAt.IsZero() {
			fmt.Printf("-user=%v\n", u)
			t.Fatalf("expected created at: %v", newUser)
		} else if newUser.UpdatedAt.IsZero() {
			t.Fatalf("expected updated at: %v", newUser)
		}

		// Create second user with email.
		u2 := keygo.User{FirstName: "jane", Email: "jane@example.com"}
		if newUser, err := s.CreateUser(ctx, u2); err != nil {
			t.Fatal(err)
		} else if newUser.ID == uuid.Nil {
			t.Fatalf("ID is not set")
		}

		// Fetch user from database & compare.
		if other, err := s.FindUserByID(ctx, newUser.ID); err != nil {
			t.Fatal(err)
		} else if !compareUsers(newUser, other) {
			t.Fatalf("mismatch: %#v != %#v", newUser, other)
		}
	})

	// Ensure an error is returned if user's name is not set.
	t.Run("ErrNameRequired", func(t *testing.T) {
		ctx := testContext(db.DB)

		s := db.NewUserService()
		if _, err := s.CreateUser(ctx, keygo.User{}); err == nil {
			t.Fatal("expected error")
		} else if keygo.ErrorCode(err) != keygo.ERR_INVALID || keygo.ErrorMessage(err) != `FirstName required.` {
			t.Fatalf("unexpected error: %#v", err)
		}
	})
}

func compareUsers(user keygo.User, other keygo.User) bool {
	other.CreatedAt = user.CreatedAt
	other.UpdatedAt = user.UpdatedAt
	return reflect.DeepEqual(user, other)
}

// MustCreateUser creates a user in the database. Fatal on error.
func MustCreateUser(tb testing.TB, ctx echo.Context, user keygo.User) (keygo.User, echo.Context) {
	tb.Helper()
	if _, err := db.NewUserService().CreateUser(ctx, user); err != nil {
		tb.Fatal(err)
	}
	return user, keygo.NewContextWithUser(ctx, user)
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func testContext(tx *gorm.DB) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("tx", tx)
	return ctx
}

func testContextWithUser(tx *gorm.DB, user keygo.User) echo.Context {
	return keygo.NewContextWithUser(testContext(tx), user)
}
