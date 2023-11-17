package mock

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/briskt/keygo/app"
)

type TenantService struct {
	Tenants map[string]app.Tenant

	FindTenantsFn func(ctx echo.Context, filter app.TenantFilter) ([]app.Tenant, int, error)
}

// Ensure service implements interface.
var _ app.TenantService = (*TenantService)(nil)

func NewTenantService() TenantService {
	return TenantService{
		Tenants: map[string]app.Tenant{},
	}
}

func (m *TenantService) DeleteAllTenants() {
	m.Tenants = map[string]app.Tenant{}
}

func (m *TenantService) FindTenantByID(context echo.Context, id string) (app.Tenant, error) {
	// TODO: decide if this is a better API than using the ID in TenantFilter passed to FindTenants
	t, ok := m.Tenants[id]
	if !ok {
		return app.Tenant{}, fmt.Errorf("no Tenant found by ID %q", id)
	}
	return t, nil
}

func (m *TenantService) FindTenants(context echo.Context, filter app.TenantFilter) ([]app.Tenant, int, error) {
	if m.FindTenantsFn != nil {
		return m.FindTenantsFn(context, filter)
	}
	var Tenants []app.Tenant
	for _, t := range m.Tenants {
		if filter.Name != nil && *filter.Name != t.Name {
			continue
		}
		if filter.ID != nil && *filter.ID != t.ID {
			continue
		}

		Tenants = append(Tenants, t)
	}
	return Tenants, len(Tenants), nil
}

func (m *TenantService) CreateTenant(context echo.Context, input app.TenantCreateInput) (app.Tenant, error) {
	if err := input.Validate(); err != nil {
		return app.Tenant{}, err
	}
	now := time.Now()
	Tenant := app.Tenant{
		ID:        newID(),
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.Tenants[Tenant.ID] = Tenant
	return Tenant, nil
}

func (m *TenantService) UpdateTenant(context echo.Context, id string, input app.TenantUpdateInput) (app.Tenant, error) {
	if err := input.Validate(); err != nil {
		return app.Tenant{}, err
	}
	panic("implement mock TenantService UpdateTenant")
}

func (m *TenantService) DeleteTenant(context echo.Context, id string) error {
	panic("implement mock TenantService DeleteTenant")
}

func (m *TenantService) CreateTenantUser(context echo.Context, input app.TenantUserCreateInput) error {
	panic("implement mock TenantService CreateTenantUser")
}
