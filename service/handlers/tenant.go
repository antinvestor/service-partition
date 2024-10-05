package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
)

func (prtSrv *PartitionServer) GetTenant(
	ctx context.Context,
	req *partitionv1.GetTenantRequest) (*partitionv1.GetTenantResponse, error) {
	logger := prtSrv.Service.L(ctx)
	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.GetTenant(ctx, req.GetId())
	if err != nil {
		logger.Debug("could not obtain the specified tenant")
		return nil, prtSrv.toApiError(err)
	}
	return &partitionv1.GetTenantResponse{Data: tenant}, nil
}

func (prtSrv *PartitionServer) ListTenant(req *partitionv1.ListTenantRequest, stream partitionv1.PartitionService_ListTenantServer) error {
	ctx := stream.Context()
	logger := prtSrv.Service.L(ctx)
	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	err := tenantBusiness.ListTenant(ctx, req, stream)
	if err != nil {
		logger.Debug("could not list tenants")
		return prtSrv.toApiError(err)
	}
	return nil
}

func (prtSrv *PartitionServer) CreateTenant(ctx context.Context, req *partitionv1.CreateTenantRequest) (*partitionv1.CreateTenantResponse, error) {
	logger := prtSrv.Service.L(ctx)
	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.CreateTenant(ctx, req)
	if err != nil {
		logger.Debug("could not create a new tenant")
		return nil, prtSrv.toApiError(err)
	}
	return &partitionv1.CreateTenantResponse{Data: tenant}, nil
}
