package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
	"github.com/schparky/keygo/db"
)

func RegisterUserRoutes(e *echo.Echo) {
	// Route => handler
	e.GET("/users", usersHandler)
}

func usersHandler(c echo.Context) error {
	s := db.NewUserService(db.Tx(c))
	u, _, err := s.FindUsers(keygo.UserFilter{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, u)
}
