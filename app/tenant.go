package app

import (
	"time"
)

const minTenantNameLength = 3

// Tenant is the full model that identifies an app Tenant
type Tenant struct {
	ID        string
	Name      string
	UserIDs   []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TenantCreateInput is a set of fields to define a new tenant for CreateTenant()
type TenantCreateInput struct {
	Name string
}

// Validate returns an error if the struct contains invalid information
func (tc *TenantCreateInput) Validate() error {
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

// TenantUpdateInput is a set of fields to be updated via UpdateTenant()
type TenantUpdateInput struct {
	Name *string
}

// Validate returns an error if the struct contains invalid information
func (tu *TenantUpdateInput) Validate() error {
	if tu.Name != nil && *tu.Name == "" {
		return Errorf(ERR_INVALID, "Tenant name is required")
	}
	if tu.Name != nil && len(*tu.Name) < minTenantNameLength {
		return Errorf(ERR_INVALID, "Tenant name must be at least %d characters", minTenantNameLength)
	}
	return nil
}

// TenantUserCreateInput is a set of fields to define a new tenant user for CreateTenantUser()
type TenantUserCreateInput struct {
	Name string
}
