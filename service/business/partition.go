package business

import (
	"context"
	"github.com/antinvestor/apis/common"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
)

type PartitionBusiness interface {
	GetPartition(ctx context.Context, request *partitionV1.PartitionGetRequest) (*partitionV1.PartitionObject, error)
	RemovePartition(ctx context.Context, request *partitionV1.PartitionRoleRemoveRequest) error
	CreatePartition(ctx context.Context, request *partitionV1.PartitionCreateRequest) (*partitionV1.PartitionObject, error)
	UpdatePartition(ctx context.Context, request *partitionV1.PartitionUpdateRequest) (*partitionV1.PartitionObject, error)

	RemovePartitionRole(ctx context.Context, request *partitionV1.PartitionRoleRemoveRequest) error
	ListPartitionRoles(ctx context.Context, request *partitionV1.PartitionRoleListRequest) (*partitionV1.PartitionRoleListResponse, error)
	CreatePartitionRole(ctx context.Context, request *partitionV1.PartitionRoleCreateRequest) (*partitionV1.PartitionRoleObject, error)
}

func NewPartitionBusiness(ctx context.Context, service *frame.Service) PartitionBusiness {
	tenantRepository := repository.NewTenantRepository(service)
	partitionRepository := repository.NewPartitionRepository(service)

	return &partitionBusiness{
		service:       service,
		partitionRepo: partitionRepository,
		tenantRepo:    tenantRepository,
	}
}

type partitionBusiness struct {
	service       *frame.Service
	tenantRepo    repository.TenantRepository
	partitionRepo repository.PartitionRepository
}

func toApiPartition(partitionModel *models.Partition) *partitionV1.PartitionObject {

	properties := extractProperties(partitionModel.Properties)

	return &partitionV1.PartitionObject{
		TenantId:    partitionModel.TenantID,
		ParentId:    partitionModel.ParentID,
		Name:        partitionModel.Name,
		Description: partitionModel.Description,
		Properties:  properties,
		State:       common.STATE(partitionModel.State),
	}
}

func toApiPartitionRole(partitionModel *models.PartitionRole) *partitionV1.PartitionRoleObject {

	properties := extractProperties(partitionModel.Properties)

	return &partitionV1.PartitionRoleObject{
		PartitionId: partitionModel.PartitionID,
		Name:        partitionModel.Name,
		Properties:  properties,
	}
}

func (pb *partitionBusiness) GetPartition(ctx context.Context, request *partitionV1.PartitionGetRequest) (*partitionV1.PartitionObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}


	partition, err := pb.partitionRepo.GetByID(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	return toApiPartition(partition), nil
}

func (pb *partitionBusiness) RemovePartition(ctx context.Context, request *partitionV1.PartitionRoleRemoveRequest) error {

	err := request.Validate()
	if err != nil {
		return err
	}


	err = pb.partitionRepo.Delete(ctx, request.GetPartitionRoleId())
	if err != nil {
		return err
	}

	return nil
}

func (pb *partitionBusiness) CreatePartition(ctx context.Context, request *partitionV1.PartitionCreateRequest) (*partitionV1.PartitionObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	tenant, err := pb.tenantRepo.GetByID(ctx, request.GetTenantId())
	if err != nil {
		return nil, err
	}

	jsonMap := make(datatypes.JSONMap)
	for k, v := range request.GetProperties() {
		jsonMap[k] = v
	}

	partition := &models.Partition{
		ParentID:    request.GetParentId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Properties:  jsonMap,
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = pb.partitionRepo.Save(ctx, partition)
	if err != nil {
		return nil, err
	}

	return toApiPartition(partition), nil
}

func (pb *partitionBusiness) UpdatePartition(ctx context.Context, request *partitionV1.PartitionUpdateRequest) (*partitionV1.PartitionObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	partition, err := pb.partitionRepo.GetByID(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	jsonMap := partition.Properties
	for k, v := range request.GetProperties() {
		jsonMap[k] = v
	}

	partition.Name = request.GetName()
	partition.Description = request.GetDescription()
	partition.Properties = jsonMap

	err = pb.partitionRepo.Save(ctx, partition)
	if err != nil {
		return nil, err
	}

	return toApiPartition(partition), nil
}

func (pb *partitionBusiness) ListPartitionRoles(ctx context.Context, request *partitionV1.PartitionRoleListRequest) (*partitionV1.PartitionRoleListResponse, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	partitionRoleList, err := pb.partitionRepo.GetRoles(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	response := make([]*partitionV1.PartitionRoleObject, 0)

	for _, pat := range partitionRoleList {
		response = append(response, toApiPartitionRole(pat))
	}

	return &partitionV1.PartitionRoleListResponse{
		Role: response,
	}, nil
}

func (pb *partitionBusiness) RemovePartitionRole(ctx context.Context, request *partitionV1.PartitionRoleRemoveRequest) error {

	err := request.Validate()
	if err != nil {
		return err
	}

	err = pb.partitionRepo.RemoveRole(ctx, request.GetPartitionRoleId())
	if err != nil {
		return err
	}

	return nil
}

func (pb *partitionBusiness) CreatePartitionRole(ctx context.Context, request *partitionV1.PartitionRoleCreateRequest) (*partitionV1.PartitionRoleObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	partition, err := pb.partitionRepo.GetByID(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	jsonMap := make(datatypes.JSONMap)
	for k, v := range request.GetProperties() {
		jsonMap[k] = v
	}

	partitionRole := &models.PartitionRole{
		Name:       request.GetName(),
		Properties: jsonMap,
		BaseModel: frame.BaseModel{
			PartitionID: partition.PartitionID,
			TenantID:    partition.TenantID,
		},
	}

	err = pb.partitionRepo.SaveRole(ctx, partitionRole)
	if err != nil {
		return nil, err
	}

	return toApiPartitionRole(partitionRole), nil
}
