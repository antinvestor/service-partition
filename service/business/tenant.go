package business

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type TenantBusiness interface {
	GetTenant(ctx context.Context, tenantId string) (*partitionV1.TenantObject, error)
	CreateTenant(ctx context.Context, request *partitionV1.TenantRequest) (*partitionV1.TenantObject, error)
}

func NewTenantBusiness(ctx context.Context, service *frame.Service) TenantBusiness {
	tenantRepo := repository.NewTenantRepository(service)
	return &tenantBusiness{
		service:    service,
		tenantRepo: tenantRepo,
	}
}

type tenantBusiness struct {
	service    *frame.Service
	tenantRepo repository.TenantRepository
}



func toApiTenant(tenantModel *models.Tenant) *partitionV1.TenantObject {

	properties := frame.DBPropertiesToMap(tenantModel.Properties)

	return &partitionV1.TenantObject{
		TenantId:    tenantModel.ID,
		Description: tenantModel.Description,
		Properties:  properties,
	}
}

func toModelTenant(tenantApi *partitionV1.TenantObject) *models.Tenant {

	return &models.Tenant{
		Description: tenantApi.GetDescription(),
		Properties:  frame.DBPropertiesFromMap(tenantApi.GetProperties()),
	}
}

func (t *tenantBusiness) GetTenant(ctx context.Context, tenantId string) (*partitionV1.TenantObject, error) {

	//err := request.Validate()
	//if err != nil {
	//	return nil, err
	//}


	tenant, err := t.tenantRepo.GetByID(ctx, tenantId)
	if err != nil {
		return nil, err
	}

	return toApiTenant(tenant), nil
}

func (t *tenantBusiness) CreateTenant(ctx context.Context, request *partitionV1.TenantRequest) (*partitionV1.TenantObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}


	jsonMap := make(datatypes.JSONMap)
	for k, v := range request.GetProperties() {
		jsonMap[k] = v
	}

	tenantModel := &models.Tenant{
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Properties:  jsonMap,
	}

	err = t.tenantRepo.Save(ctx, tenantModel)
	if err != nil {
		return nil, err
	}

	return toApiTenant(tenantModel), nil
}
