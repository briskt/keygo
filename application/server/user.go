package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo"
)

func (s *Server) userHandler(c echo.Context) error {
	user := keygo.CurrentUser(c)

	id := c.Param("id")
	if id != user.ID {
		return c.JSON(http.StatusNotFound, AuthError{Error: "not found"})
	}

	return c.JSON(http.StatusOK, user)
}
