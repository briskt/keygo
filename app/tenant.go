package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

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

// Validate returns an error if the tenant contains invalid fields.
// This only performs basic validation.
func (u *TenantCreate) Validate() error {
	if len(u.Name) < 3 {
		return Errorf(ERR_INVALID, "Tenant name is required")
	}
	return nil
}

// TenantFilter is a filter passed to FindTenants()
type TenantFilter struct {
	// Filtering fields.
	ID     *string `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// TenantUpdate is a set of fields to be updated via UpdateTenant()
type TenantUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
