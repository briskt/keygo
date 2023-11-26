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

func (u *Tenant) BeforeCreate(_ *gorm.DB) error {
	u.ID = newID()
	return nil
}

// FindTenants retrieves a list of tenants by filter. Also returns total count of
// matching tenants which may differ from returned results if filter.Limit is specified.
func FindTenants(ctx echo.Context, _ app.TenantFilter) ([]Tenant, error) {
	var tenants []Tenant
	result := Tx(ctx).Find(&tenants)
	if result.Error != nil {
		return []Tenant{}, result.Error
	}

	return tenants, nil
}

// CreateTenant creates a new tenant.
func CreateTenant(ctx echo.Context, input app.TenantCreateInput) (Tenant, error) {
	if err := input.Validate(); err != nil {
		return Tenant{}, err
	}

	newTenant := Tenant{
		Name: input.Name,
	}
	err := Tx(ctx).Create(&newTenant).Error
	if err != nil {
		return Tenant{}, err
	}

	return newTenant, nil
}

// UpdateTenant updates a tenant object.
func UpdateTenant(ctx echo.Context, id string, input app.TenantUpdateInput) (Tenant, error) {
	if err := input.Validate(); err != nil {
		return Tenant{}, err
	}

	tenant, err := FindTenantByID(ctx, id)
	if err != nil {
		return Tenant{}, err
	}

	if input.Name != nil {
		tenant.Name = *input.Name
	}

	result := Tx(ctx).Save(&tenant)
	if err != nil {
		return Tenant{}, result.Error
	}

	return tenant, nil
}

// DeleteTenant permanently deletes a tenant and all child objects
func DeleteTenant(ctx echo.Context, id string) error {
	result := Tx(ctx).Where("id = ?", id).Delete(&Tenant{})
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// CreateTenantUser permanently deletes a tenant and all child objects
func CreateTenantUser(ctx echo.Context, input app.TenantUserCreateInput) error {
	return nil
}

// FindTenantByID is a function to fetch a tenant by ID.
func FindTenantByID(ctx echo.Context, id string) (Tenant, error) {
	var tenant Tenant
	result := Tx(ctx).First(&tenant, "id = ?", id)
	return tenant, result.Error
}

func ConvertTenant(c echo.Context, t Tenant) (app.Tenant, error) {
	tenant := app.Tenant{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	users, err := FindUsers(c, app.UserFilter{TenantID: &t.ID})
	if err != nil {
		return app.Tenant{}, err
	}
	tenant.UserIDs = make([]string, len(users))
	for i, user := range users {
		tenant.UserIDs[i] = user.ID
	}
	return tenant, nil
}
