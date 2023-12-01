package business_test

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"reflect"
	"testing"
)

func Test_extractProperties(t *testing.T) {
	type args struct {
		props datatypes.JSONMap
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := frame.DBPropertiesToMap(tt.args.props); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tenantBusiness_CreateTenant(t1 *testing.T) {

	ctx := context.Background()

	type fields struct {
		service    *frame.Service
		tenantRepo repository.TenantRepository
	}
	type args struct {
		ctx     context.Context
		request *partitionv1.CreateTenantRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *partitionv1.TenantObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := business.NewTenantBusinessWithRepo(ctx, tt.fields.service, tt.fields.tenantRepo)
			got, err := t.CreateTenant(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t1.Errorf("CreateTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("CreateTenant() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tenantBusiness_GetTenant(t1 *testing.T) {

	ctx := context.Background()
	type fields struct {
		service    *frame.Service
		tenantRepo repository.TenantRepository
	}
	type args struct {
		ctx      context.Context
		tenantId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *partitionv1.TenantObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := business.NewTenantBusinessWithRepo(ctx, tt.fields.service, tt.fields.tenantRepo)
			got, err := t.GetTenant(tt.args.ctx, tt.args.tenantId)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetTenant() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toApiTenant(t *testing.T) {
	type args struct {
		tenantModel *models.Tenant
	}
	tests := []struct {
		name string
		args args
		want *partitionv1.TenantObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := business.ToApiTenant(tt.args.tenantModel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToApiTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toModelTenant(t *testing.T) {
	type args struct {
		tenantApi *partitionv1.TenantObject
	}
	tests := []struct {
		name string
		args args
		want *models.Tenant
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := business.ToModelTenant(tt.args.tenantApi); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToModelTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}
