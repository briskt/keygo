package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/schparky/keygo"
)

func RegisterUserRoutes(e *echo.Echo) {
	// Route => handler
	e.GET("/user", userHandler)
}

func userHandler(c echo.Context) error {
	i := c.Get(keygo.ContextKeyToken)
	token, ok := i.(keygo.Token)
	if !ok {
		return c.JSON(http.StatusBadRequest, "no token")
	}

	return c.JSON(http.StatusOK, token.Auth.User)
}
