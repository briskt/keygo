package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (s *Server) usersListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusOK, []app.User{})
	}

	users, err := db.FindUsers(c, app.UserFilter{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d users", len(users))

	return c.JSON(http.StatusOK, users)
}

func (s *Server) userHandler(c echo.Context) error {
	user := app.CurrentUser(c)

	id := c.Param("id")
	if id != user.ID {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	return c.JSON(http.StatusOK, user)
}
