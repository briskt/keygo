package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	Id        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Role      string    `db:"role"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func RegisterUserRoutes(e *echo.Echo) {
	// Route => handler
	e.GET("/users", usersHandler)
}

func usersHandler(c echo.Context) error {
	tmp := c.Get("tx")
	tx, ok := tmp.(*gorm.DB)
	if !ok {
		panic("no transaction found in context")
	}

	u := make([]User, 0)
	result := tx.Find(&u)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}

	return c.JSON(http.StatusOK, u)
}
