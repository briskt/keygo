package app

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Tenant struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate returns an error if the tenant contains invalid fields.
// This only performs basic validation.
func (u *Tenant) Validate() error {
	if len(u.Name) < 3 {
		return Errorf(ERR_INVALID, "Tenant name is required")
	}
	return nil
}

// TenantService represents a service for managing tenants
type TenantService interface {
	// FindTenantByID retrieves a tenant by ID
	FindTenantByID(echo.Context, string) (Tenant, error)

	// FindTenants retrieves a list of tenants by filter
	FindTenants(echo.Context, TenantFilter) ([]Tenant, int, error)

	// CreateTenant creates a new tenant
	CreateTenant(echo.Context, Tenant) (Tenant, error)

	// UpdateTenant updates a tenant object
	UpdateTenant(echo.Context, string, TenantUpdate) (Tenant, error)

	// DeleteTenant permanently deletes a tenant and all child objects
	DeleteTenant(echo.Context, string) error
}

// TenantFilter represents a filter passed to FindTenants()
type TenantFilter struct {
	// Filtering fields.
	ID     *string `json:"id"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// TenantUpdate represents a set of fields to be updated via UpdateTenant()
type TenantUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
