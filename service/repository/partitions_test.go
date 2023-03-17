package repository_test

import (
	"context"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/antinvestor/service-partition/testsutil"
	"github.com/pitabwire/frame"
	"testing"
)

func TestPartitionRepository_GetByID(t *testing.T) {

	ctx := context.Background()
	srv, err := testsutil.GetTestService("Partition Srv", ctx)
	if err != nil {
		t.Errorf("There was an error getting service : %v", err)
		return
	}

	tenantRepo := repository.NewTenantRepository(srv)
	tenant := models.Tenant{
		Name:        "Save T",
		Description: "Test",
	}

	err = tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := repository.NewPartitionRepository(srv)
	partition := models.Partition{
		Name:        "",
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

	savedPartition, err := partitionRepo.GetByID(ctx, partition.GetID())
	if err != nil {
		t.Errorf("There was an error getting partition : %v", err)
		return
	}

	if partition.GetID() != savedPartition.GetID() {
		t.Errorf("The obtained partition doesn't match what was saved")
		return
	}

}

func TestPartitionRepository_GetChildren(t *testing.T) {

	ctx := context.Background()
	srv, err := testsutil.GetTestService("Partition Srv", ctx)
	if err != nil {
		t.Errorf("There was an error getting service : %v", err)
		return
	}

	tenantRepo := repository.NewTenantRepository(srv)
	tenant := models.Tenant{
		Name:        "Save T",
		Description: "Test",
	}

	err = tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := repository.NewPartitionRepository(srv)
	partition := models.Partition{
		Name:        "",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.ID,
		},
	}

	err = partitionRepo.Save(ctx, &partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	childPartition := models.Partition{
		Name:        "",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.ID,
		},
		ParentID: partition.GetID(),
	}

	err = partitionRepo.Save(ctx, &childPartition)
	if err != nil {
		t.Errorf("There was an error saving child partition : %v", err)
		return
	}

	childrentPartitions, err := partitionRepo.GetChildren(ctx, partition.GetID())
	if err != nil {
		t.Errorf("There was an error getting children partition : %v", err)
		return
	}

	if len(childrentPartitions) != 1 {
		t.Errorf("There should be only one child partition now")
		return
	}

	if childrentPartitions[0].ParentID != partition.GetID() {
		t.Errorf("Child partition parent id: %v should match parent partition id: %v", childrentPartitions[0].ParentID, partition.GetID())
		return
	}

}

func TestPartitionRepository_SaveRole(t *testing.T) {

	ctx := context.Background()
	srv, err := testsutil.GetTestService("Partition Srv", ctx)
	if err != nil {
		t.Errorf("There was an error getting service : %v", err)
		return
	}

	tenantRepo := repository.NewTenantRepository(srv)
	tenant := models.Tenant{
		Name:        "Save T",
		Description: "Test",
	}

	err = tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := repository.NewPartitionRepository(srv)
	partition := models.Partition{
		Name:        "",
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

	partitionRole := models.PartitionRole{
		Name: "",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = partitionRepo.SaveRole(ctx, &partitionRole)
	if err != nil {
		t.Errorf("There was an error saving partition role : %v", err)
		return
	}

	partitionRoles, err := partitionRepo.GetRoles(ctx, partition.GetID())
	if err != nil {
		t.Errorf("There was an error getting partition roles : %v", err)
		return
	}

	if len(partitionRoles) != 1 {
		t.Errorf("There should be only one partition role now")
		return
	}

	if partitionRoles[0].PartitionID != partition.GetID() {
		t.Errorf("Partition role partition id: %v should match parent partition id: %v", partitionRoles[0].PartitionID, partition.GetID())
		return
	}

}

func TestPartitionRepository_RemoveRole(t *testing.T) {

	ctx := context.Background()
	srv, err := testsutil.GetTestService("Partition Srv", ctx)
	if err != nil {
		t.Errorf("There was an error getting service : %v", err)
		return
	}

	tenantRepo := repository.NewTenantRepository(srv)
	tenant := models.Tenant{
		Name:        "Save T",
		Description: "Test",
	}

	err = tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := repository.NewPartitionRepository(srv)

	partition := models.Partition{
		Name:        "",
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

	partitionRole := models.PartitionRole{
		Name: "",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = partitionRepo.SaveRole(ctx, &partitionRole)
	if err != nil {
		t.Errorf("There was an error saving partition role : %v", err)
		return
	}

	partitionRoles, err := partitionRepo.GetRoles(ctx, partition.GetID())
	if err != nil {
		t.Errorf("There was an error getting partition roles : %v", err)
		return
	}

	if len(partitionRoles) != 1 {
		t.Errorf("There should be only one partition role now")
		return
	}

	err = partitionRepo.RemoveRole(ctx, partitionRoles[0].GetID())
	if err != nil {
		t.Errorf("There was an error removing partition roles : %v", err)
		return
	}

	partitionRoles, err = partitionRepo.GetRoles(ctx, partition.GetID())
	if err != nil {
		t.Errorf("There was an error getting partition roles : %v", err)
		return
	}

	if len(partitionRoles) != 0 {
		t.Errorf("There should be no partition role now")
		return
	}

}
