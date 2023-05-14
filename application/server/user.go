package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
)

func (s *Server) userHandler(c echo.Context) error {
	i := c.Get(keygo.ContextKeyToken)
	token, ok := i.(keygo.Token)
	if !ok {
		return c.JSON(http.StatusBadRequest, AuthError{Error: "no token"})
	}

	id := c.Param("id")
	if id != token.Auth.UserID.String() {
		return c.JSON(http.StatusNotFound, AuthError{Error: "not found"})
	}

	return c.JSON(http.StatusOK, token.Auth.User)
}
