package business

import (
	"context"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"os"
	"testing"
)

func TestPartitionBusiness_SyncPartitionOnHydra(t *testing.T) {

	ctx := context.Background()
	service := getTestService("Partition Srv", ctx)

	err := os.Setenv(config.EnvOauth2ServiceAdminUri, "http://localhost:4445")
	if err != nil {
		t.Errorf("There was an error setting HYDRA_URL : %v", err)
		return
	}

	tenantRepo := repository.NewTenantRepository(service)
	tenant := models.Tenant{
		Name:        "default",
		Description: "Test",
	}

	err = tenantRepo.Save(ctx, &tenant)
	if err != nil {
		t.Errorf("There was an error saving tenant : %v", err)
		return
	}

	partitionRepo := repository.NewPartitionRepository(service)
	partition := &models.Partition{
		Name:        "test partition",
		Description: "",
		BaseModel: frame.BaseModel{
			TenantID: tenant.GetID(),
		},
	}

	err = partitionRepo.Save(ctx, partition)
	if err != nil {
		t.Errorf("There was an error saving partition : %v", err)
		return
	}

	err = SyncPartitionOnHydra(ctx, service, partition)
	if err != nil {
		t.Errorf("Could not sync this partition : %v", err)
		return
	}

}
