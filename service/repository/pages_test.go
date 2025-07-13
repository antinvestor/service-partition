package repository_test

import (
	"github.com/antinvestor/service-partition/internal/tests"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type PageTestSuite struct {
	tests.BaseTestSuite
}

func (suite *PageTestSuite) TestGetByPartitionAndName() {
	// Test cases
	testCases := []struct {
		name        string
		pageName    string
		shouldError bool
	}{
		{
			name:        "Get page by partition and name",
			pageName:    "test",
			shouldError: false,
		},
	}
	
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		pageRepo := repository.NewPageRepository(svc)
		tenantRepo := repository.NewTenantRepository(svc)
		partitionRepo := repository.NewPartitionRepository(svc)
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				tenant := models.Tenant{
					Name:        "default",
					Description: "Test",
				}

				err := tenantRepo.Save(ctx, &tenant)
				require.NoError(t, err)

				partition := models.Partition{
					Name:        "Test Partition",
					Description: "Test partition description",
					BaseModel: frame.BaseModel{
						TenantID: tenant.GetID(),
					},
				}

				err = partitionRepo.Save(ctx, &partition)
				require.NoError(t, err)

				page := models.Page{
					Name: tc.pageName,
					HTML: "<div></div>",
					BaseModel: frame.BaseModel{
						TenantID:    tenant.GetID(),
						PartitionID: partition.GetID(),
					},
				}

				err = pageRepo.Save(ctx, &page)
				require.NoError(t, err)

				// Execute
				savedPage, err := pageRepo.GetByPartitionAndName(ctx, partition.GetID(), page.Name)
				
				// Verify
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, partition.GetID(), savedPage.PartitionID, "Page partition id should match parent partition id")
					assert.Equal(t, page.GetID(), savedPage.GetID(), "Page ID should match saved page ID")
				}
			})
		}
	})
}

func (suite *PageTestSuite) TestSave() {
	// Test cases
	testCases := []struct {
		name        string
		pageName    string
		html        string
		shouldError bool
	}{
		{
			name:        "Save page",
			pageName:    "test",
			html:        "<div></div>",
			shouldError: false,
		},
	}
	
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		pageRepo := repository.NewPageRepository(svc)
		tenantRepo := repository.NewTenantRepository(svc)
		partitionRepo := repository.NewPartitionRepository(svc)
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				tenant := models.Tenant{
					Name:        "default",
					Description: "Test",
				}

				err := tenantRepo.Save(ctx, &tenant)
				require.NoError(t, err)

				partition := models.Partition{
					Name:        "Test Partition",
					Description: "Test partition description",
					BaseModel: frame.BaseModel{
						TenantID: tenant.GetID(),
					},
				}

				err = partitionRepo.Save(ctx, &partition)
				require.NoError(t, err)

				page := models.Page{
					Name: tc.pageName,
					HTML: tc.html,
					BaseModel: frame.BaseModel{
						TenantID:    tenant.GetID(),
						PartitionID: partition.GetID(),
					},
				}

				// Execute
				err = pageRepo.Save(ctx, &page)
				
				// Verify
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					
					savedPage, err := pageRepo.GetByID(ctx, page.GetID())
					assert.NoError(t, err)
					assert.Equal(t, partition.GetID(), savedPage.PartitionID, "Page partition id should match parent partition id")
					assert.Equal(t, tc.pageName, savedPage.Name, "Page name should match")
					assert.Equal(t, tc.html, savedPage.HTML, "Page HTML should match")
				}
			})
		}
	})
}

func (suite *PageTestSuite) TestDelete() {
	// Test cases
	testCases := []struct {
		name        string
		shouldError bool
	}{
		{
			name:        "Delete page",
			shouldError: false,
		},
	}
	
	suite.WithTestDependancies(suite.T(), func(t *testing.T, dep *testdef.DependancyOption) {
		svc, ctx := suite.CreateService(t, dep)
		pageRepo := repository.NewPageRepository(svc)
		tenantRepo := repository.NewTenantRepository(svc)
		partitionRepo := repository.NewPartitionRepository(svc)
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Setup
				tenant := models.Tenant{
					Name:        "default",
					Description: "Test",
				}

				err := tenantRepo.Save(ctx, &tenant)
				require.NoError(t, err)

				partition := models.Partition{
					Name:        "Test Partition",
					Description: "Test partition description",
					BaseModel: frame.BaseModel{
						TenantID: tenant.GetID(),
					},
				}

				err = partitionRepo.Save(ctx, &partition)
				require.NoError(t, err)

				page := models.Page{
					Name: "test",
					HTML: "<div></div>",
					BaseModel: frame.BaseModel{
						TenantID:    tenant.GetID(),
						PartitionID: partition.GetID(),
					},
				}

				err = pageRepo.Save(ctx, &page)
				require.NoError(t, err)

				// Execute
				err = pageRepo.Delete(ctx, page.GetID())
				
				// Verify
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					
					deletedPage, err := pageRepo.GetByID(ctx, page.GetID())
					assert.Error(t, err, "Should get an error when fetching deleted page")
					assert.True(t, strings.Contains(err.Error(), "record not found"), "Error should mention 'record not found'")
					assert.Empty(t, deletedPage.GetID(), "Deleted page ID should be empty")
				}
			})
		}
	})
}

// TestPageRepository runs the page repository test suite
func TestPageRepository(t *testing.T) {
	suite.Run(t, new(PageTestSuite))
}
