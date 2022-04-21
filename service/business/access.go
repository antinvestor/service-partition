package business

import (
	"context"
	"errors"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"log"
	"strings"
)

type AccessBusiness interface {
	GetAccess(ctx context.Context, request *partitionV1.AccessGetRequest) (*partitionV1.AccessObject, error)
	RemoveAccess(ctx context.Context, request *partitionV1.AccessRemoveRequest) error
	CreateAccess(ctx context.Context, request *partitionV1.AccessCreateRequest) (*partitionV1.AccessObject, error)

	RemoveAccessRole(ctx context.Context, request *partitionV1.AccessRoleRemoveRequest) error
	ListAccessRoles(
		ctx context.Context,
		request *partitionV1.AccessRoleListRequest) (*partitionV1.AccessRoleListResponse, error)
	CreateAccessRole(
		ctx context.Context,
		request *partitionV1.AccessRoleCreateRequest) (*partitionV1.AccessRoleObject, error)
}

func NewAccessBusiness(ctx context.Context, service *frame.Service) AccessBusiness {
	accessRepo := repository.NewAccessRepository(service)
	partitionRepo := repository.NewPartitionRepository(service)

	return &accessBusiness{
		service:       service,
		accessRepo:    accessRepo,
		partitionRepo: partitionRepo,
	}
}

type accessBusiness struct {
	service       *frame.Service
	accessRepo    repository.AccessRepository
	partitionRepo repository.PartitionRepository
}

func toAPIAccess(
	partitionObject *partitionV1.PartitionObject,
	accessModel *models.Access) (*partitionV1.AccessObject, error) {

	if partitionObject == nil {
		return nil, errors.New("no partition exists for this access")
	}

	return &partitionV1.AccessObject{
		AccessId:  accessModel.GetID(),
		ProfileId: accessModel.ProfileID,
		Partition: partitionObject,
	}, nil
}

func toAPIAccessRole(partitionRoleObj *partitionV1.PartitionRoleObject, accessRoleModel *models.AccessRole) *partitionV1.AccessRoleObject {

	return &partitionV1.AccessRoleObject{
		AccessRoleId: accessRoleModel.GetID(),
		AccessId:     accessRoleModel.AccessID,
		Role:         partitionRoleObj,
	}
}

func (ab *accessBusiness) GetAccess(
	ctx context.Context,
	request *partitionV1.AccessGetRequest) (*partitionV1.AccessObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	var access *models.Access
	if request.GetAccessId() != "" {
		access, err = ab.accessRepo.GetByID(ctx, request.GetAccessId())
		if err != nil {
			return nil, err
		}

		partition, err := ab.partitionRepo.GetByID(ctx, access.PartitionID)
		if err != nil {
			return nil, err
		}

		partitionObject := toAPIPartition(partition)

		return toAPIAccess(partitionObject, access)
	}

	access, err = ab.accessRepo.GetByPartitionAndProfile(ctx, request.GetPartitionId(), request.GetProfileId())
	if err != nil {
		return nil, err
	}

	partition, err := ab.partitionRepo.GetByID(ctx, access.PartitionID)
	if err != nil {
		return nil, err
	}

	partitionObject := toAPIPartition(partition)

	return toAPIAccess(partitionObject, access)
}

func (ab *accessBusiness) RemoveAccess(
	ctx context.Context,
	request *partitionV1.AccessRemoveRequest) error {

	err := request.Validate()
	if err != nil {
		return err
	}

	err = ab.accessRepo.Delete(ctx, request.GetAccessId())
	if err != nil {
		return err
	}

	return nil
}

func (ab *accessBusiness) CreateAccess(
	ctx context.Context,
	request *partitionV1.AccessCreateRequest) (*partitionV1.AccessObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	log.Printf(" CreateAccess -- supplied request %+v", request)

	partition, err := ab.partitionRepo.GetByID(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	access, err := ab.accessRepo.GetByPartitionAndProfile(ctx, request.GetPartitionId(), request.GetProfileId())
	if err != nil {

		if !strings.Contains(err.Error(), "record not found") {
			return nil, err
		}
	} else {
		partitionObject := toAPIPartition(partition)

		return toAPIAccess(partitionObject, access)
	}

	access = &models.Access{
		ProfileID: request.GetProfileId(),
		BaseModel: frame.BaseModel{
			TenantID:    partition.TenantID,
			PartitionID: partition.GetID(),
		},
	}

	err = ab.accessRepo.Save(ctx, access)
	if err != nil {
		return nil, err
	}

	log.Printf(" CreateAccess -- final access created is  %+v", access)
	partitionObject := toAPIPartition(partition)

	return toAPIAccess(partitionObject, access)
}

func (ab *accessBusiness) ListAccessRoles(
	ctx context.Context,
	request *partitionV1.AccessRoleListRequest) (*partitionV1.AccessRoleListResponse, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	accessRoleList, err := ab.accessRepo.GetRoles(ctx, request.GetAccessId())
	if err != nil {
		return nil, err
	}

	parititionRoleIDs := make([]string, 0)

	for _, accessR := range accessRoleList {
		parititionRoleIDs = append(parititionRoleIDs, accessR.PartitionRoleID)
	}

	partitionRoles, err := ab.partitionRepo.GetRolesByID(ctx, parititionRoleIDs...)
	if err != nil {
		return nil, err
	}

	partitionRoleIDMap := make(map[string]*partitionV1.PartitionRoleObject)
	for _, partitionRole := range partitionRoles {
		partitionRoleIDMap[partitionRole.ID] = toAPIPartitionRole(partitionRole)
	}

	response := make([]*partitionV1.AccessRoleObject, 0)

	for _, acc := range accessRoleList {
		response = append(response, toAPIAccessRole(partitionRoleIDMap[acc.PartitionRoleID], acc))
	}

	return &partitionV1.AccessRoleListResponse{
		Role: response,
	}, nil
}

func (ab *accessBusiness) RemoveAccessRole(
	ctx context.Context,
	request *partitionV1.AccessRoleRemoveRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	err = ab.accessRepo.RemoveRole(ctx, request.GetAccessRoleId())
	if err != nil {
		return err
	}

	return nil
}

func (ab *accessBusiness) CreateAccessRole(
	ctx context.Context,
	request *partitionV1.AccessRoleCreateRequest) (*partitionV1.AccessRoleObject, error) {

	err := request.Validate()
	if err != nil {
		return nil, err
	}

	access, err := ab.accessRepo.GetByID(ctx, request.GetAccessId())
	if err != nil {
		return nil, err
	}

	partitionRoles, err := ab.partitionRepo.GetRolesByID(ctx, request.GetPartitionRoleId())
	if err != nil {
		return nil, err
	}

	accessRole := &models.AccessRole{
		AccessID:        access.GetID(),
		PartitionRoleID: partitionRoles[0].GetID(),
	}

	err = ab.accessRepo.SaveRole(ctx, accessRole)
	if err != nil {
		return nil, err
	}

	partitionRoleObj := toAPIPartitionRole(partitionRoles[0])
	return toAPIAccessRole(partitionRoleObj, accessRole), nil
}
