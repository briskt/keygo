package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (s *Server) usersListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != "Admin" {
		return c.JSON(http.StatusNotFound, AuthError{Error: "not found"})
	}

	users, n, err := s.UserService.FindUsers(c, app.UserFilter{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d users", n)

	return c.JSON(http.StatusOK, users)
}

func (s *Server) userHandler(c echo.Context) error {
	user := app.CurrentUser(c)

	id := c.Param("id")
	if id != user.ID {
		return c.JSON(http.StatusNotFound, AuthError{Error: "not found"})
	}

	return c.JSON(http.StatusOK, user)
}
