package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
)

func (s *Server) registerUserRoutes() {
	e := s.Echo
	e.GET("/user", userHandler)
}

func userHandler(c echo.Context) error {
	i := c.Get(keygo.ContextKeyToken)
	token, ok := i.(keygo.Token)
	if !ok {
		return c.JSON(http.StatusBadRequest, "no token")
	}

	return c.JSON(http.StatusCreated, token.Auth.User)
}
