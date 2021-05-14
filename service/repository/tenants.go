package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
)

type tenantRepository struct {
	service *frame.Service
}

func (tr *tenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error){
	tenant := &models.Tenant{}
	err := tr.service.DB(ctx, true).First(tenant, "id = ?", id).Error
	return tenant, err
}

func (tr *tenantRepository) Save(ctx context.Context, tenant *models.Tenant) error{
	return tr.service.DB(ctx, false).Save(tenant).Error
}

func (tr *tenantRepository) Delete(ctx context.Context, id string) error{

	tenant, err := tr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return tr.service.DB(ctx, false).Delete(tenant).Error
}

func NewTenantRepository( service *frame.Service )  TenantRepository{
	tenantRepository := tenantRepository{
		service: service,
	}
	return  &tenantRepository
}
