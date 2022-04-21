package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
)

type accessRepository struct {
	service *frame.Service
}

func (ar *accessRepository) GetByID(ctx context.Context, id string) (*models.Access, error) {
	access := &models.Access{}
	err := ar.service.DB(ctx, true).First(access, " accesses.id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return access, nil
}

func (ar *accessRepository) GetByPartitionAndProfile(ctx context.Context, partitionId string, profileId string) (*models.Access, error) {
	access := &models.Access{}
	err := ar.service.DB(ctx, true).First(access, " partition_id = ? AND profile_id = ?", partitionId, profileId).Error
	if err != nil {
		return nil, err
	}

	return access, nil
}

func (ar *accessRepository) Save(ctx context.Context, access *models.Access) error {
	return ar.service.DB(ctx, false).Save(access).Error
}

func (ar *accessRepository) Delete(ctx context.Context, id string) error {

	err := ar.service.DB(ctx, false).Where(" access_id = ?", id).Delete(&models.AccessRole{}).Error
	if err != nil {
		return err
	}

	return ar.service.DB(ctx, false).Where(" id = ?", id).Delete(&models.Access{}).Error
}

func (ar *accessRepository) GetRoles(ctx context.Context, accessId string) ([]*models.AccessRole, error) {
	accessRoles := make([]*models.AccessRole, 0)
	err := ar.service.DB(ctx, true).
		Find(&accessRoles, " access_id = ?", accessId).Error

	return accessRoles, err
}

func (ar *accessRepository) SaveRole(ctx context.Context, role *models.AccessRole) error {
	return ar.service.DB(ctx, false).Save(role).Error
}

func (ar *accessRepository) RemoveRole(ctx context.Context, accessRoleId string) error {
	return ar.service.DB(ctx, false).Where(" id = ?", accessRoleId).Delete(&models.AccessRole{}).Error
}

func NewAccessRepository(service *frame.Service) AccessRepository {
	partitionRepository := accessRepository{
		service: service,
	}
	return &partitionRepository
}
