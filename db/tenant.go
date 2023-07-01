package db

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/briskt/keygo/app"
)

type Tenant struct {
	ID        string `gorm:"primaryKey;type:string"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   gorm.DeletedAt
}

func (u *Tenant) BeforeCreate(tx *gorm.DB) error {
	u.ID = newID()
	return nil
}

// Ensure service implements interface.
var _ app.TenantService = (*TenantService)(nil)

// TenantService is a service for managing tenants.
type TenantService struct{}

// NewTenantService returns a new instance of TenantService.
func NewTenantService() *TenantService {
	return &TenantService{}
}

// FindTenantByID retrieves a tenant by ID along with their associated auth objects.
func (s *TenantService) FindTenantByID(ctx echo.Context, id string) (app.Tenant, error) {
	tenant, err := findTenantByID(ctx, id)
	if err != nil {
		return app.Tenant{}, err
	}
	return exportTenant(tenant), nil
}

// FindTenants retrieves a list of tenants by filter. Also returns total count of
// matching tenants which may differ from returned results if filter.Limit is specified.
func (s *TenantService) FindTenants(ctx echo.Context, filter app.TenantFilter) ([]app.Tenant, int, error) {
	var tenants []Tenant
	result := Tx(ctx).Find(&tenants)
	if result.Error != nil {
		return []app.Tenant{}, 0, result.Error
	}
	appTenants := make([]app.Tenant, len(tenants))
	for i := range tenants {
		appTenants[i] = exportTenant(tenants[i])
	}
	return appTenants, len(tenants), nil
}

// CreateTenant creates a new tenant.
func (s *TenantService) CreateTenant(ctx echo.Context, input app.TenantCreate) (app.Tenant, error) {
	if err := input.Validate(); err != nil {
		return app.Tenant{}, err
	}

	newTenant := Tenant{
		Name: input.Name,
	}
	result := Tx(ctx).Create(&newTenant)

	return exportTenant(newTenant), result.Error
}

// UpdateTenant updates a tenant object.
func (s *TenantService) UpdateTenant(ctx echo.Context, id string, input app.TenantUpdate) (app.Tenant, error) {
	tenant, err := findTenantByID(ctx, id)
	if err != nil {
		return app.Tenant{}, err
	}

	if input.Name != nil {
		tenant.Name = *input.Name
	}

	result := Tx(ctx).Save(&tenant)
	if err != nil {
		return app.Tenant{}, result.Error
	}

	return exportTenant(tenant), nil
}

// DeleteTenant permanently deletes a tenant and all child objects
func (s *TenantService) DeleteTenant(ctx echo.Context, id string) error {
	result := Tx(ctx).Where("id = ?", id).Delete(&Tenant{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// findTenantByID is a helper function to fetch a tenant by ID.
func findTenantByID(ctx echo.Context, id string) (Tenant, error) {
	var tenant Tenant
	result := Tx(ctx).First(&tenant, "id = ?", id)
	return tenant, result.Error
}

func exportTenant(t Tenant) app.Tenant {
	return app.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
