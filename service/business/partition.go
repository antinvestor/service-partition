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

	properties := frame.DBPropertiesToMap(partitionModel.Properties)

	return &partitionV1.PartitionObject{
		PartitionId: partitionModel.ID,
		TenantId:    partitionModel.TenantID,
		ParentId:    partitionModel.ParentID,
		Name:        partitionModel.Name,
		Description: partitionModel.Description,
		Properties:  properties,
		State:       common.STATE(partitionModel.State),
	}
}

func toApiPartitionRole(partitionModel *models.PartitionRole) *partitionV1.PartitionRoleObject {

	properties := frame.DBPropertiesToMap(partitionModel.Properties)

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


	partition := &models.Partition{
		ParentID:    request.GetParentId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Properties:  frame.DBPropertiesFromMap(request.GetProperties()),
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = pb.partitionRepo.Save(ctx, partition)
	if err != nil {
		return nil, err
	}

	err = pb.service.Publish(ctx, config.QueuePartitionSyncName, partition)
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

func ReQueuePrimaryPartitionsForSync(service *frame.Service) {
	ctx := context.Background()
	partitionRepository := repository.NewPartitionRepository(service)
	partition, err := partitionRepository.GetByID(ctx, "c2f4j7au6s7f91uqnokg")
	if err != nil {
		log.Printf(" ReQueuePrimaryPartitionsForSync -- could not get default system partition because :%v+", err)
		return
	}

	err = service.Publish(ctx, config.QueuePartitionSyncName, partition)
	if err != nil {
		log.Printf(" ReQueuePrimaryPartitionsForSync -- could not publish because :%v+", err)
		return
	}

	partition, err = partitionRepository.GetByID(ctx, "9bsv0s3pbdv002o80qhg")
	if err != nil {
		log.Printf(" ReQueuePrimaryPartitionsForSync -- could not get test partition because :%v+", err)
		return
	}

	err = service.Publish(ctx, config.QueuePartitionSyncName, partition)
	if err != nil {
		log.Printf(" ReQueuePrimaryPartitionsForSync -- could not publish because :%v+", err)
		return
	}
}

func SyncPartitionOnHydra(ctx context.Context, service *frame.Service, partition *models.Partition) error {

	hydraUrl := fmt.Sprintf("%s%s", frame.GetEnv(config.EnvOauth2ServiceAdminUri, ""), "/clients")
	hydraIDUrl := fmt.Sprintf("%s/%s", hydraUrl, partition.ID)

	if partition.DeletedAt.Valid {

		//	We need to delete this partition on hydra as well
		_, _, err := service.InvokeRestService(ctx, http.MethodDelete, hydraIDUrl, make(map[string]interface{}), nil)
		return err
	}

	status, result, err := service.InvokeRestService(ctx, http.MethodGet, hydraIDUrl, make(map[string]interface{}), nil)
	if err != nil {
		return err
	}

	httpMethod := http.MethodPost
	if status == 200 {
		//	We need to update this partition on hydra as well as it already exists
		httpMethod = http.MethodPut
		hydraUrl = hydraIDUrl
	}

	logoUri := ""
	if val, ok := partition.Properties["logo_uri"]; ok {
		logoUri = val.(string)
	}

	var audienceList []string
	if val, ok := partition.Properties["audience"]; ok {
		audIfc := val.([]interface{})
		for _, v := range audIfc {
			aud := v.(string)
			audienceList = append(audienceList, aud)
		}
	} else {
		audienceList = []string{}
	}

	var uriList []string
	if val, ok := partition.Properties["redirect_uris"]; ok {
		redirectUri, ok := val.(string)
		if ok {
			uriList = strings.Split(redirectUri, ",")

		} else {

			redirectUris, ok := val.([]interface{})
			if ok {
				for _, v := range redirectUris {
					uriList = append(uriList, fmt.Sprintf("%v", v))
				}
			} else {
				log.Printf(" SyncPartitionOnHydra -- The required redirect uri list is invalid %v", val)
			}
		}
	}

	payload := map[string]interface{}{
		"client_name":                partition.Name,
		"grant_types":                []string{"authorization_code", "refresh_token"},
		"response_types":             []string{"token", "id_token", "code"},
		"scope":                      "openid offline offline_access profile contact",
		"redirect_uris":              uriList,
		"logo_uri":                   logoUri,
		"audience":                   audienceList,
		"token_endpoint_auth_method": "none",
	}

	if httpMethod == http.MethodPost {
		payload["client_id"] = partition.ID
	}

	status, result, err = service.InvokeRestService(ctx, httpMethod, hydraUrl, payload, nil)
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

	if partition.Properties == nil {
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
