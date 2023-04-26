package repository_test

import (
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/antinvestor/service-partition/testsutil"
	"strings"
	"testing"
)

func TestTenantRepository_GetByID(t *testing.T) {
	ctx, srv, err := testsutil.GetTestService("Tenant Srv")
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
	}

	savedTenant, err := tenantRepo.GetByID(ctx, tenant.GetID())
	if err != nil {
		t.Errorf("There was an error getting tenant : %v", err)
	}

	if tenant.GetID() != savedTenant.GetID() {
		t.Errorf("The obtained tenant doesn't match")
	}

}

func TestTenantRepository_Save(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Tenant Srv")
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
	}

}

func TestTenantRepository_Delete(t *testing.T) {

	ctx, srv, err := testsutil.GetTestService("Tenant Srv")
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
	}

	err = tenantRepo.Delete(ctx, tenant.GetID())
	if err != nil {
		t.Errorf("There was an error deleting tenant : %v", err)
	}

	deletedTenant, err := tenantRepo.GetByID(ctx, tenant.GetID())
	if err != nil && !strings.Contains(err.Error(), "record not found") {

		t.Errorf("There was an error getting tenant : %v", err)
	}

	if deletedTenant != nil && deletedTenant.ID != "" {
		t.Errorf("Tenant %v is expected to be deleted", deletedTenant)
	}

}
