package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
	"github.com/briskt/keygo/db"
)

func (s *Server) tenantsCreateHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	var input app.TenantCreateInput
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	tenant, err := db.CreateTenant(c, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("created tenant (name %q, id %q)", tenant.Name, tenant.ID)

	return c.JSON(http.StatusOK, tenant)
}

func (s *Server) tenantsListHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	tenants, err := db.FindTenants(c, app.TenantFilter{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, AuthError{Error: err.Error()})
	}

	s.Logger.Infof("found %d tenants", len(tenants))

	return c.JSON(http.StatusOK, tenants)
}

func (s *Server) tenantsGetHandler(c echo.Context) error {
	user := app.CurrentUser(c)

	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	id := c.Param("id")
	tenant, err := db.FindTenantByID(c, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	t, err := db.ConvertTenant(c, tenant)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, t)
}

func (s *Server) tenantsUsersCreateHandler(c echo.Context) error {
	user := app.CurrentUser(c)
	if user.Role != app.UserRoleAdmin {
		return echo.NewHTTPError(http.StatusNotFound, AuthError{Error: "not found"})
	}

	var input app.TenantUserCreateInput
	err := (&echo.DefaultBinder{}).BindBody(c, &input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}

	tenantID := c.Param("id")
	tenantUser, err := db.CreateTenantUser(c, tenantID, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	s.Logger.Infof("created tenant user (email %q, id %q)", tenantUser.Email, tenantUser.ID)

	return c.JSON(http.StatusOK, tenantUser)
}
