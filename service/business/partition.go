package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antinvestor/apis/common"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"log"
	"net/http"
	"strings"
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

	partitionAsByte, err := json.Marshal(partition)
	if err != nil {
		return nil, err
	}

	err = pb.service.Publish(ctx, config.QueuePartitionSyncName, partitionAsByte)
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




func SyncPartitionOnHydra(ctx context.Context, partition *models.Partition) error {

	service := frame.FromContext(ctx)
	if service == nil {
		return errors.New("no service was found in the context provided")
	}

	hydraUrl := fmt.Sprintf("%s%s", frame.GetEnv("HYDRA_URL", ""), "/clients")

	logoUri := ""
	if val, ok := partition.Properties["logo_uri"]; ok {
		logoUri = val.(string)
	}
	redirectUri := ""
	if val, ok := partition.Properties["request_uris"]; ok {
		redirectUri = val.(string)

		redirectUri = strings.Replace(redirectUri, "[","", 1)
		redirectUri = strings.Replace(redirectUri, "]","", 1)
	}



	payload := map[string]interface{}{
		"client_id": partition.ID,
		"client_name":    partition.Name,
		"grant_types":  []string{"authorization_code", "refresh_token"},
		"token_endpoint_auth_method":   "none",
		"response_types": []string{"token", "id_token", "code"},
		"scope": "openid,offline_access,profile,contact",
		"request_uris": strings.Split(redirectUri, ","),
		"logo_uri": logoUri,
	}

	status, result, err := service.InvokeRestService(ctx, http.MethodPost, hydraUrl, payload, nil)
	if err != nil {
		return err
	}

	if status > 299 || status < 200 {
		return errors.New(fmt.Sprintf(" invalid response status %d had message %s", status, string(result)))
	}

	var response map[string]interface{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return err
	}

	log.Printf("Returned response is : %v", response)

	if partition.Properties == nil{
		partition.Properties = make(datatypes.JSONMap)
	}
	for k, v := range response {
		partition.Properties[k] = v
	}

	partitionRepository := repository.NewPartitionRepository(service)
	err = partitionRepository.Save(ctx, partition)
	if err != nil {
		return err
	}

	return nil
}
