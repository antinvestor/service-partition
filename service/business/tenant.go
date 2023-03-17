package business

import (
	"context"
	partitionv1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type TenantBusiness interface {
	GetTenant(ctx context.Context, tenantId string) (*partitionv1.TenantObject, error)
	CreateTenant(ctx context.Context, request *partitionv1.TenantRequest) (*partitionv1.TenantObject, error)
	ListTenant(ctx context.Context, request *partitionv1.SearchRequest, stream partitionv1.PartitionService_ListTenantServer) error
}

func NewTenantBusiness(ctx context.Context, service *frame.Service) TenantBusiness {
	tenantRepo := repository.NewTenantRepository(service)

	return NewTenantBusinessWithRepo(ctx, service, tenantRepo)
}
func NewTenantBusinessWithRepo(_ context.Context, service *frame.Service, repo repository.TenantRepository) TenantBusiness {

	return &tenantBusiness{
		service:    service,
		tenantRepo: repo,
	}
}

type tenantBusiness struct {
	service    *frame.Service
	tenantRepo repository.TenantRepository
}

func ToApiTenant(tenantModel *models.Tenant) *partitionv1.TenantObject {

	properties := frame.DBPropertiesToMap(tenantModel.Properties)

	return &partitionv1.TenantObject{
		TenantId:    tenantModel.ID,
		Description: tenantModel.Description,
		Properties:  properties,
	}
}

func ToModelTenant(tenantApi *partitionv1.TenantObject) *models.Tenant {

	return &models.Tenant{
		Description: tenantApi.GetDescription(),
		Properties:  frame.DBPropertiesFromMap(tenantApi.GetProperties()),
	}
}

func (t *tenantBusiness) GetTenant(ctx context.Context, tenantId string) (*partitionv1.TenantObject, error) {

	//err := request.Validate()
	//if err != nil {
	//	return nil, err
	//}

	tenant, err := t.tenantRepo.GetByID(ctx, tenantId)
	if err != nil {
		return nil, err
	}

	return ToApiTenant(tenant), nil
}

func (t *tenantBusiness) CreateTenant(ctx context.Context, request *partitionv1.TenantRequest) (*partitionv1.TenantObject, error) {

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

	return ToApiTenant(tenantModel), nil
}

func (t *tenantBusiness) ListTenant(ctx context.Context, request *partitionv1.SearchRequest, stream partitionv1.PartitionService_ListTenantServer) error {

	err := request.Validate()
	if err != nil {
		return err
	}

	tenantList, err := t.tenantRepo.GetByQuery(ctx, request.GetQuery(), request.GetCount(), request.GetPage())
	if err != nil {
		return err
	}

	for _, tenant := range tenantList {

		err = stream.Send(ToApiTenant(tenant))
		if err != nil {
			return err
		}

	}

	return nil
}
