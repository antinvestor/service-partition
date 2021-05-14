package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
)

type TenantRepository interface {
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	Save(ctx context.Context, tenant *models.Tenant) error
	Delete(ctx context.Context, id string) error
}

type PartitionRepository interface {
	GetByID(ctx context.Context, id string) (*models.Partition, error)
	GetChildren(ctx context.Context, id string) ([]*models.Partition, error)
	Save(ctx context.Context, partition *models.Partition) error
	Delete(ctx context.Context, id string) error

	GetRoles(ctx context.Context, id string) ([]*models.PartitionRole, error)
	SaveRole(ctx context.Context, role *models.PartitionRole) error
	RemoveRole(ctx context.Context, partitionRoleId string) error
}

type PageRepository interface {
	GetByID(ctx context.Context, id string) (*models.Page, error)
	GetByPartitionAndName(ctx context.Context, partitionId string, name string) (*models.Page, error)
	Save(ctx context.Context, partition *models.Page) error
	Delete(ctx context.Context, id string) error
}

type AccessRepository interface {
	GetByID(ctx context.Context, id string) (*models.Access, error)
	GetByPartitionAndProfile(ctx context.Context, partitionId string, profile string) (*models.Access, error)
	Save(ctx context.Context, access *models.Access) error
	Delete(ctx context.Context, id string) error

	GetRoles(ctx context.Context, accessId string) ([]*models.AccessRole, error)
	SaveRole(ctx context.Context, role *models.AccessRole) error
	RemoveRole(ctx context.Context, accessRoleId string) error
}
