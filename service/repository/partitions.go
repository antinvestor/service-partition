package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
)

type partitionRepository struct {
	service *frame.Service
}

func (pr *partitionRepository) GetByID(ctx context.Context, id string) (*models.Partition, error) {
	partition := &models.Partition{}
	err := pr.service.DB(ctx, true).First(partition, "id = ?", id).Error
	return partition, err
}

func (pr *partitionRepository) GetChildren(ctx context.Context, id string) ([]*models.Partition, error) {
	childPartition := make([]*models.Partition, 0)
	err := pr.service.DB(ctx, true).Find(&childPartition, "parent_id = ?", id).Error
	return childPartition, err
}

func (pr *partitionRepository) Save(ctx context.Context, partition *models.Partition) error {
	return pr.service.DB(ctx, false).Save(partition).Error
}

func (pr *partitionRepository) Delete(ctx context.Context, id string) error {

	partition, err := pr.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return pr.service.DB(ctx, false).Delete(partition).Error
}

func (pr *partitionRepository) GetRoles(ctx context.Context, partitionId string) ([]*models.PartitionRole, error) {
	partitionRoles := make([]*models.PartitionRole, 0)
	err := pr.service.DB(ctx, true).Find(&partitionRoles, "partition_id = ?", partitionId).Error
	return partitionRoles, err
}

func (pr *partitionRepository) SaveRole(ctx context.Context, role *models.PartitionRole) error {
	return pr.service.DB(ctx, false).Save(role).Error
}

func (pr *partitionRepository) RemoveRole(ctx context.Context, partitionRoleId string) error {
	return pr.service.DB(ctx, false).Where("id = ?", partitionRoleId).Delete(&models.PartitionRole{}).Error
}

func NewPartitionRepository(service *frame.Service) PartitionRepository {
	partitionRepository := partitionRepository{
		service: service,
	}
	return &partitionRepository
}
