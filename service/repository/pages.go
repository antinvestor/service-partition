package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
)

type pageRepository struct {
	service *frame.Service
}

func (pgr *pageRepository) GetByID(ctx context.Context, id string) (*models.Page, error){
	page := &models.Page{}
	err := pgr.service.DB(ctx, true).First(page, "id = ?", id).Error
	return page, err
}

func (pgr *pageRepository) GetByPartitionAndName(ctx context.Context, partitionId string, name string) (*models.Page, error){
	page := &models.Page{}
	err := pgr.service.DB(ctx, true).First(page, "partition_id = ? AND name = ?", partitionId, name).Error
	return page, err
}

func (pgr *pageRepository) Save(ctx context.Context, page *models.Page) error{
	return pgr.service.DB(ctx, false).Save(page).Error
}

func (pgr *pageRepository) Delete(ctx context.Context, id string) error{
	return pgr.service.DB(ctx, false).Where("id = ?", id).Delete(&models.Page{}).Error
}

func NewPageRepository( service *frame.Service )  PageRepository{
	tenantRepository := pageRepository{
		service: service,
	}
	return  &tenantRepository
}

