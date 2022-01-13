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
	i := c.Get("token")
	token := i.(keygo.Token)

	return c.JSON(http.StatusOK, token.Auth.User)
}
