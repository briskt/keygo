package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

const minTenantNameLength = 3

// TenantService is a service for managing tenants
type TenantService interface {
	// FindTenantByID retrieves a tenant by ID
	FindTenantByID(ctx echo.Context, id string) (Tenant, error)

	// FindTenants retrieves a list of tenants by filter
	FindTenants(ctx echo.Context, tenantFilter TenantFilter) ([]Tenant, int, error)

	// CreateTenant creates a new tenant
	CreateTenant(ctx echo.Context, tenantCreate TenantCreate) (Tenant, error)

	// UpdateTenant updates a tenant object
	UpdateTenant(ctx echo.Context, id string, tenantUpdate TenantUpdate) (Tenant, error)

	// DeleteTenant permanently deletes a tenant and all child objects
	DeleteTenant(ctx echo.Context, id string) error
}

// Tenant is the full model that identifies an app Tenant
type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TenantCreate is a set of fields to define a new user for CreateTenant()
type TenantCreate struct {
	Name string
}

// Validate returns an error if the struct contains invalid information
func (tc *TenantCreate) Validate() error {
	if tc.Name == "" {
		return Errorf(ERR_INVALID, "Tenant name is required")
	}
	if len(tc.Name) < minTenantNameLength {
		return Errorf(ERR_INVALID, "Tenant name must be at least %d characters", minTenantNameLength)
	}
	return nil
}

// TenantFilter is a filter passed to FindTenants()
type TenantFilter struct {
	// Filtering fields.
	ID   *string
	Name *string

	// Restrict to subset of results.
	Offset int
	Limit  int
}

// TenantUpdate is a set of fields to be updated via UpdateTenant()
type TenantUpdate struct {
	Name *string
}

// Validate returns an error if the struct contains invalid information
func (tu *TenantUpdate) Validate() error {
	if tu.Name != nil && *tu.Name == "" {
		return Errorf(ERR_INVALID, "Tenant name is required")
	}
	if tu.Name != nil && len(*tu.Name) < minTenantNameLength {
		return Errorf(ERR_INVALID, "Tenant name must be at least %d characters", minTenantNameLength)
	}
	return nil
}
