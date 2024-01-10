package business

import (
	"context"
	"encoding/json"
	"fmt"
	commonv1 "github.com/antinvestor/apis/go/common/v1"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
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
	GetPartition(ctx context.Context, request *partitionv1.GetPartitionRequest) (*partitionv1.PartitionObject, error)
	CreatePartition(
		ctx context.Context,
		request *partitionv1.CreatePartitionRequest) (*partitionv1.PartitionObject, error)
	UpdatePartition(
		ctx context.Context,
		request *partitionv1.UpdatePartitionRequest) (*partitionv1.PartitionObject, error)
	ListPartition(
		ctx context.Context,
		request *partitionv1.ListPartitionRequest,
		stream partitionv1.PartitionService_ListPartitionServer) error

	RemovePartitionRole(ctx context.Context, request *partitionv1.RemovePartitionRoleRequest) error
	ListPartitionRoles(
		ctx context.Context,
		request *partitionv1.ListPartitionRoleRequest) (*partitionv1.ListPartitionRoleResponse, error)
	CreatePartitionRole(
		ctx context.Context,
		request *partitionv1.CreatePartitionRoleRequest) (*partitionv1.PartitionRoleObject, error)
}

func NewPartitionBusiness(service *frame.Service) PartitionBusiness {
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

func toAPIPartition(partitionModel *models.Partition) *partitionv1.PartitionObject {
	properties := frame.DBPropertiesToMap(partitionModel.Properties)

	return &partitionv1.PartitionObject{
		Id:          partitionModel.ID,
		TenantId:    partitionModel.TenantID,
		ParentId:    partitionModel.ParentID,
		Name:        partitionModel.Name,
		Description: partitionModel.Description,
		Properties:  properties,
		State:       commonv1.STATE(partitionModel.State),
	}
}

func toAPIPartitionRole(partitionModel *models.PartitionRole) *partitionv1.PartitionRoleObject {
	properties := frame.DBPropertiesToMap(partitionModel.Properties)

	return &partitionv1.PartitionRoleObject{
		PartitionId: partitionModel.PartitionID,
		Name:        partitionModel.Name,
		Properties:  properties,
	}
}

func (pb *partitionBusiness) ListPartition(ctx context.Context, request *partitionv1.ListPartitionRequest, stream partitionv1.PartitionService_ListPartitionServer) error {

	partitionList, err := pb.partitionRepo.GetByQuery(ctx, request.GetQuery(), uint32(request.GetCount()), uint32(request.GetPage()))
	if err != nil {
		return err
	}

	var responseObjects []*partitionv1.PartitionObject
	for _, partition := range partitionList {

		responseObjects = append(responseObjects, toAPIPartition(partition))
	}

	err = stream.Send(&partitionv1.ListPartitionResponse{Data: responseObjects})
	if err != nil {
		return err
	}

	return nil
}

func (pb *partitionBusiness) GetPartition(
	ctx context.Context,
	request *partitionv1.GetPartitionRequest) (*partitionv1.PartitionObject, error) {

	partition, err := pb.partitionRepo.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return toAPIPartition(partition), nil
}

func (pb *partitionBusiness) CreatePartition(
	ctx context.Context,
	request *partitionv1.CreatePartitionRequest) (*partitionv1.PartitionObject, error) {

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

	partitionConfig := pb.service.Config().(*config.PartitionConfig)

	err = pb.service.Publish(ctx, partitionConfig.PartitionSyncName, partition)
	if err != nil {
		return nil, err
	}

	return toAPIPartition(partition), nil
}

func (pb *partitionBusiness) UpdatePartition(
	ctx context.Context,
	request *partitionv1.UpdatePartitionRequest) (*partitionv1.PartitionObject, error) {

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

	return toAPIPartition(partition), nil
}

func (pb *partitionBusiness) ListPartitionRoles(
	ctx context.Context,
	request *partitionv1.ListPartitionRoleRequest,
) (*partitionv1.ListPartitionRoleResponse, error) {

	partitionRoleList, err := pb.partitionRepo.GetRoles(ctx, request.GetPartitionId())
	if err != nil {
		return nil, err
	}

	response := make([]*partitionv1.PartitionRoleObject, 0)

	for _, pat := range partitionRoleList {
		response = append(response, toAPIPartitionRole(pat))
	}

	return &partitionv1.ListPartitionRoleResponse{
		Role: response,
	}, nil
}

func (pb *partitionBusiness) RemovePartitionRole(
	ctx context.Context,
	request *partitionv1.RemovePartitionRoleRequest,
) error {

	err := pb.partitionRepo.RemoveRole(ctx, request.GetPartitionRoleId())
	if err != nil {
		return err
	}

	return nil
}

func (pb *partitionBusiness) CreatePartitionRole(
	ctx context.Context,
	request *partitionv1.CreatePartitionRoleRequest) (
	*partitionv1.PartitionRoleObject, error) {

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

	return toAPIPartitionRole(partitionRole), nil
}

func ReQueuePrimaryPartitionsForSync(service *frame.Service) {
	ctx := context.Background()
	partitionRepository := repository.NewPartitionRepository(service)
	partitionConfig := service.Config().(*config.PartitionConfig)

	partitionList, err := partitionRepository.GetByQuery(ctx, "", 100, 0)
	if err != nil {
		log.Printf(" ReQueuePrimaryPartitionsForSync -- could not get default system partition because :%v+", err)

		return
	}

	for _, partition := range partitionList {
		err = service.Publish(ctx, partitionConfig.PartitionSyncName, partition)
		if err != nil {
			log.Printf(" ReQueuePrimaryPartitionsForSync -- could not publish because :%v+", err)

			return
		}
	}
}

func SyncPartitionOnHydra(ctx context.Context, service *frame.Service, partition *models.Partition) error {

	partitionConfig := service.Config().(*config.PartitionConfig)

	hydraURL := fmt.Sprintf("%s%s", partitionConfig.GetOauth2ServiceAdminURI(), "/admin/clients")
	httpMethod := http.MethodPost

	clientID, ok := partition.Properties["client_id"]
	if ok {

		hydraIDUrl := fmt.Sprintf("%s/%s", hydraURL, clientID)

		if partition.DeletedAt.Valid { //	We need to delete this partition on hydra as well
			_, _, err := service.InvokeRestService(ctx, http.MethodDelete, hydraIDUrl, make(map[string]interface{}), nil)
			return err
		}

		status, _, err := service.InvokeRestService(
			ctx,
			http.MethodGet,
			hydraIDUrl,
			make(map[string]interface{}),
			nil)
		if err != nil {
			return err
		}

		if status == 200 {
			//	We need to update this partition on hydra as well as it already exists
			httpMethod = http.MethodPut
			hydraURL = hydraIDUrl
		}

	}

	logoURI := ""
	if val, ok := partition.Properties["logo_uri"]; ok {
		logoURI = val.(string)
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
		redirectURI, ok := val.(string)
		if ok {
			uriList = strings.Split(redirectURI, ",")

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
		"response_types":             []string{"token", "id_token", "code", "token id_token", "token code id_token"},
		"scope":                      "openid offline offline_access profile contact",
		"redirect_uris":              uriList,
		"logo_uri":                   logoURI,
		"audience":                   audienceList,
		"token_endpoint_auth_method": "none",
	}

	status, result, err := service.InvokeRestService(ctx, httpMethod, hydraURL, payload, nil)
	if err != nil {
		return err
	}

	if status > 299 || status < 200 {
		return fmt.Errorf(" invalid response status %d had message %s", status, string(result))
	}

	var response map[string]interface{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return err
	}

	clientId, ok := partition.Properties["client_id"].(string)
	if ok {
		partition.ClientID = clientId
	}

	if partition.Properties == nil {
		partition.Properties = make(datatypes.JSONMap)
	}
	for k, v := range response {
		partition.Properties[k] = v
	}

	partitionRepository := repository.NewPartitionRepository(service)
	err = partitionRepository.Save(ctx, partition)
	if err != nil { //nolint:wsl
		return err
	}

	return nil
}
