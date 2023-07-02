package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (s *Server) usersListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != "Admin" {
		return echo.NewHTTPError(http.StatusOK, []app.User{})
	}

	users, n, err := s.UserService.FindUsers(c, app.UserFilter{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d users", n)

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

func (s *Server) userTokensListHandler(c echo.Context) error {
	userID := c.Param("id")

	tokens, err := s.TokenService.ListTokensForUser(c, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d tokens for user %s", len(tokens), userID)
	return c.JSON(http.StatusOK, tokens)
}
