package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (s *Server) tenantsListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != "Admin" {
		return c.JSON(http.StatusNotFound, AuthError{Error: "not found"})
	}

	tenants, n, err := s.TenantService.FindTenants(c, app.TenantFilter{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d tenants", n)

	return c.JSON(http.StatusOK, tenants)
}
