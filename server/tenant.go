package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

func (s *Server) tenantsCreateHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != "Admin" {
		return echo.NewHTTPError(http.StatusUnauthorized, AuthError{Error: "not an authorized user"})
	}

	var input app.TenantCreate
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	tenant, err := s.TenantService.CreateTenant(c, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("created tenant (name %q, id %q)", tenant.Name, tenant.ID)

	return c.JSON(http.StatusOK, tenant)
}

func (s *Server) tenantsListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != "Admin" {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	tenants, n, err := s.TenantService.FindTenants(c, app.TenantFilter{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d tenants", n)

	return c.JSON(http.StatusOK, tenants)
}
