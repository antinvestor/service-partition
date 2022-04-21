package business

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/pitabwire/frame"
	"gorm.io/datatypes"
	"reflect"
	"testing"
)

func getTestService(name string, ctx context.Context) *frame.Service {
	mainDb := frame.Datastore(ctx, "postgres://partition:secret@localhost:5423/partitiondatabase?sslmode=disable", false)
	return frame.NewService(name, mainDb)
}

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
	type fields struct {
		service    *frame.Service
		tenantRepo repository.TenantRepository
	}
	type args struct {
		ctx     context.Context
		request *partitionV1.TenantRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *partitionV1.TenantObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tenantBusiness{
				service:    tt.fields.service,
				tenantRepo: tt.fields.tenantRepo,
			}
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
		want    *partitionV1.TenantObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tenantBusiness{
				service:    tt.fields.service,
				tenantRepo: tt.fields.tenantRepo,
			}
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
		want *partitionV1.TenantObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toApiTenant(tt.args.tenantModel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toApiTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toModelTenant(t *testing.T) {
	type args struct {
		tenantApi *partitionV1.TenantObject
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
			if got := toModelTenant(tt.args.tenantApi); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toModelTenant() = %v, want %v", got, tt.want)
			}
		})
	}
}
