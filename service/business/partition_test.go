package business_test

import (
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/internal/tests"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PartitionBusinessTestSuite struct {
	tests.BaseTestSuite
}

func (p *PartitionBusinessTestSuite) TestSyncPartitionOnHydra() {
	// Test cases
	testCases := []struct {
		name        string
		shouldError bool
	}{
		{
			name:        "Sync partition on Hydra",
			shouldError: false,
		},
	}

	p.WithTestDependancies(p.T(), func(t *testing.T, dep *testdef.DependancyOption) {

		svc, ctx := p.CreateService(t, dep)

		hydraHost, cleanupContainer, err := p.CreateHydraContainer(t)
		defer cleanupContainer()
		require.NoError(t, err)

		cfg, ok := svc.Config().(*config.PartitionConfig)
		if ok {
			cfg.Oauth2ServiceAdminURI = hydraHost
		}

		tenantRepo := repository.NewTenantRepository(svc)
		partitionRepo := repository.NewPartitionRepository(svc)

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				tenant := models.Tenant{
					Name:        "default",
					Description: "Test",
				}

				err = tenantRepo.Save(ctx, &tenant)
				require.NoError(t, err)

				partition := &models.Partition{
					Name:        "test partition",
					Description: "",
					BaseModel: frame.BaseModel{
						TenantID: tenant.GetID(),
					},
				}

				err = partitionRepo.Save(ctx, partition)
				require.NoError(t, err)

				// Execute
				err = business.SyncPartitionOnHydra(ctx, svc, partition)

				// Verify
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err, "Could not sync this partition")
				}
			})
		}
	})
}

// TestPartitionBusiness runs the partition business test suite
func TestPartitionBusiness(t *testing.T) {
	suite.Run(t, new(PartitionBusinessTestSuite))
}
