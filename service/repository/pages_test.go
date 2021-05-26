package repository

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/pitabwire/frame"
	"strings"
	"testing"
)


func TestPageRepository_GetByPartitionAndName(t *testing.T) {

	ctx := context.Background()
	srv := getTestService("Page Srv", ctx)

	tenantRepo := NewTenantRepository(srv)
	tenant := models.Tenant{
		Name: "Save T",
		Description: "Test",
	}

	err := tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := NewPartitionRepository(srv)
	partition := models.Partition{
		Name: "",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = partitionRepo.Save(ctx, &partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	pageRepo := NewPageRepository(srv)
	page := models.Page{
		Name: "test",
		Html: "<div></div>",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = pageRepo.Save(ctx, &page)
	if err != nil {
		t.Errorf("There was an error saving page role : %v", err)
		return
	}


	savedPage, err := pageRepo.GetByPartitionAndName(ctx, partition.GetID(), page.Name)
	if err != nil {
		t.Errorf("There was an error getting saved page : %v", err)
		return
	}

	if savedPage.PartitionID != partition.GetID() ||  savedPage.GetID() != page.GetID(){
		t.Errorf("Page role partition id: %v should match parent partition id: %v", savedPage.PartitionID, partition.GetID())
		return
	}

}

func TestPageRepository_Save(t *testing.T) {
	ctx := context.Background()
	srv := getTestService("Page Srv", ctx)

	tenantRepo := NewTenantRepository(srv)
	tenant := models.Tenant{
		Name: "Save T",
		Description: "Test",
	}

	err := tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := NewPartitionRepository(srv)
	partition := models.Partition{
		Name: "",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = partitionRepo.Save(ctx, &partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	pageRepo := NewPageRepository(srv)
	page := models.Page{
		Name: "test",
		Html: "<div></div>",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = pageRepo.Save(ctx, &page)
	if err != nil {
		t.Errorf("There was an error saving page role : %v", err)
		return
	}


	savedPage, err := pageRepo.GetByID(ctx, page.GetID())
	if err != nil {
		t.Errorf("There was an error getting saved page : %v", err)
		return
	}

	if savedPage.PartitionID != partition.GetID() {
		t.Errorf("Page role partition id: %v should match parent partition id: %v", savedPage.PartitionID, partition.GetID())
		return
	}

}

func TestPageRepository_Delete(t *testing.T) {

	ctx := context.Background()
	srv := getTestService("Page Srv", ctx)

	tenantRepo := NewTenantRepository(srv)
	tenant := models.Tenant{
		Name: "Save T",
		Description: "Test",
	}

	err := tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := NewPartitionRepository(srv)
	partition := models.Partition{
		Name: "",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = partitionRepo.Save(ctx, &partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	pageRepo := NewPageRepository(srv)
	page := models.Page{
		Name: "test",
		Html: "<div></div>",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = pageRepo.Save(ctx, &page)
	if err != nil {
		t.Errorf("There was an error saving page role : %v", err)
		return
	}


	err = pageRepo.Delete(ctx, page.GetID())
	if err != nil {
		t.Errorf("There was an error deleting saved page : %v", err)
		return
	}

	deletedPage, err := pageRepo.GetByID(ctx, page.GetID())
	if err != nil && !strings.Contains(err.Error(), "record not found") {

		t.Errorf("There was an error getting deleted page : %v", err)
		return
	}

	if deletedPage != nil && deletedPage.ID != ""{
		t.Errorf("Page : %v is supposed to be nil but somehow it exists  ", deletedPage)
		return
	}

}

