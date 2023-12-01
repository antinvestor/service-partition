package business

import (
	"context"
	"errors"
	partitionv1 "github.com/antinvestor/apis/partition/v1"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"log"
	"strings"
)

type AccessBusiness interface {
	GetAccess(ctx context.Context, request *partitionv1.GetAccessRequest) (*partitionv1.AccessObject, error)
	RemoveAccess(ctx context.Context, request *partitionv1.RemoveAccessRequest) error
	CreateAccess(ctx context.Context, request *partitionv1.CreateAccessRequest) (*partitionv1.AccessObject, error)

	RemoveAccessRole(ctx context.Context, request *partitionv1.RemoveAccessRoleRequest) error
	ListAccessRoles(
		ctx context.Context,
		request *partitionv1.ListAccessRoleRequest) (*partitionv1.ListAccessRoleResponse, error)
	CreateAccessRole(
		ctx context.Context,
		request *partitionv1.CreateAccessRoleRequest) (*partitionv1.AccessRoleObject, error)
}

func NewAccessBusiness(_ context.Context, service *frame.Service) AccessBusiness {
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
	partitionObject *partitionv1.PartitionObject,
	accessModel *models.Access) (*partitionv1.AccessObject, error) {

	if partitionObject == nil {
		return nil, errors.New("no partition exists for this access")
	}

	return &partitionv1.AccessObject{
		AccessId:  accessModel.GetID(),
		ProfileId: accessModel.ProfileID,
		Partition: partitionObject,
	}, nil
}

func toAPIAccessRole(partitionRoleObj *partitionv1.PartitionRoleObject, accessRoleModel *models.AccessRole) *partitionv1.AccessRoleObject {

	return &partitionv1.AccessRoleObject{
		AccessRoleId: accessRoleModel.GetID(),
		AccessId:     accessRoleModel.AccessID,
		Role:         partitionRoleObj,
	}
}

func (ab *accessBusiness) GetAccess(
	ctx context.Context,
	request *partitionv1.GetAccessRequest) (*partitionv1.AccessObject, error) {

	var err error
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
	request *partitionv1.RemoveAccessRequest) error {

	err := ab.accessRepo.Delete(ctx, request.GetAccessId())
	if err != nil {
		return err
	}

	return nil
}

func (ab *accessBusiness) CreateAccess(
	ctx context.Context,
	request *partitionv1.CreateAccessRequest) (*partitionv1.AccessObject, error) {

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
	request *partitionv1.ListAccessRoleRequest) (*partitionv1.ListAccessRoleResponse, error) {

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

	partitionRoleIDMap := make(map[string]*partitionv1.PartitionRoleObject)
	for _, partitionRole := range partitionRoles {
		partitionRoleIDMap[partitionRole.ID] = toAPIPartitionRole(partitionRole)
	}

	response := make([]*partitionv1.AccessRoleObject, 0)

	for _, acc := range accessRoleList {
		response = append(response, toAPIAccessRole(partitionRoleIDMap[acc.PartitionRoleID], acc))
	}

	return &partitionv1.ListAccessRoleResponse{
		Role: response,
	}, nil
}

func (ab *accessBusiness) RemoveAccessRole(
	ctx context.Context,
	request *partitionv1.RemoveAccessRoleRequest) error {

	err := ab.accessRepo.RemoveRole(ctx, request.GetAccessRoleId())
	if err != nil {
		return err
	}

	return nil
}

func (ab *accessBusiness) CreateAccessRole(
	ctx context.Context,
	request *partitionv1.CreateAccessRoleRequest) (*partitionv1.AccessRoleObject, error) {

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
