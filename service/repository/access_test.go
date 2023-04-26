package repository_test

import (
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/antinvestor/service-partition/testsutil"
	"github.com/pitabwire/frame"
	"strings"
	"testing"
)

func TestAccessRepository_Save(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Access Srv")
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

	accessRepo := repository.NewAccessRepository(srv)
	access := models.Access{
		ProfileID: "profile",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = accessRepo.Save(ctx, &access)
	if err != nil {
		t.Errorf("There was an error saving access : %v", err)
		return
	}

	savedAccess, err := accessRepo.GetByID(ctx, access.GetID())
	if err != nil {
		t.Errorf("There was an error getting saved access : %+v", err)
		return
	}

	if savedAccess.PartitionID != partition.GetID() {
		t.Errorf("Access partition id: %v should match parent partition id: %v", savedAccess.PartitionID, partition.GetID())
		return
	}
}

func TestAccessRepository_GetByPartitionAndProfile(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Access Srv")
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
		Name:        "Partition",
		Description: "Partition details",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = partitionRepo.Save(ctx, &partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	accessRepo := repository.NewAccessRepository(srv)
	access := models.Access{
		ProfileID: "profile_j",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = accessRepo.Save(ctx, &access)
	if err != nil {
		t.Errorf("There was an error saving access : %v", err)
		return
	}

	savedAccess, err := accessRepo.GetByPartitionAndProfile(ctx, partition.GetID(), "profile_j")
	if err != nil {
		t.Errorf("There was an error getting saved access : %+v", err)
		return
	}

	if savedAccess.PartitionID != partition.GetID() {
		t.Errorf("Access partition id: %v should match parent partition id: %v", savedAccess.PartitionID, partition.GetID())
		return
	}
}

func TestAccessRepository_SaveRole(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Access Srv")
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

	accessRepo := repository.NewAccessRepository(srv)
	access := models.Access{
		ProfileID: "profile_j",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = accessRepo.Save(ctx, &access)
	if err != nil {
		t.Errorf("There was an error saving access : %v", err)
		return
	}

	accessRole := models.AccessRole{
		AccessID:        access.GetID(),
		PartitionRoleID: partitionRole.GetID(),
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = accessRepo.SaveRole(ctx, &accessRole)
	if err != nil {
		t.Errorf("There was an error saving access role: %v", err)
		return
	}

	savedAccessRoles, err := accessRepo.GetRoles(ctx, access.GetID())
	if err != nil {
		t.Errorf("There was an error getting saved access : %v", err)
		return
	}

	if len(savedAccessRoles) != 1 {
		t.Errorf("We should have only one access role saved")
		return
	}

	if savedAccessRoles[0].PartitionID != partition.GetID() || savedAccessRoles[0].AccessID != access.GetID() {
		t.Errorf("Partition role partition id: %v should match parent partition id: %v", savedAccessRoles[0].PartitionID, partition.GetID())
		return
	}

}

func TestAccessRepository_RemoveRole(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Access Srv")
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

	accessRepo := repository.NewAccessRepository(srv)
	access := models.Access{
		ProfileID: "profile_j",
		BaseModel: frame.BaseModel{
			TenantID:    tenant.GetID(),
			PartitionID: partition.GetID(),
		},
	}

	err = accessRepo.Save(ctx, &access)
	if err != nil {
		t.Errorf("There was an error saving access : %v", err)
		return
	}

	accessRole := models.AccessRole{
		AccessID:        access.GetID(),
		PartitionRoleID: partitionRole.GetID(),
	}

	err = accessRepo.SaveRole(ctx, &accessRole)
	if err != nil {
		t.Errorf("There was an error saving access role: %v", err)
		return
	}

	err = accessRepo.RemoveRole(ctx, accessRole.GetID())
	if err != nil {
		t.Errorf("There was an error deleting saved access role : %v", err)
		return
	}

	deletedAccessRoles, err := accessRepo.GetRoles(ctx, access.GetID())
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		t.Errorf("There was an error getting saved access : %v", err)
		return
	}

	if len(deletedAccessRoles) != 0 {
		t.Errorf("There should be no access role now")
		return
	}

}
