package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (prtSrv *PartitionServer) GetTenant(
	ctx context.Context,
	req *partitionv1.GetTenantRequest) (*partitionv1.GetTenantResponse, error) {
	logger := prtSrv.Service.L()
	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.GetTenant(ctx, req.GetId())
	if err != nil {
		logger.Debug("could not obtain the specified tenant")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.GetTenantResponse{Data: tenant}, nil
}

func (prtSrv *PartitionServer) ListTenant(req *partitionv1.ListTenantRequest, stream partitionv1.PartitionService_ListTenantServer) error {
	logger := prtSrv.Service.L()
	tenantBusiness := business.NewTenantBusiness(stream.Context(), prtSrv.Service)
	err := tenantBusiness.ListTenant(stream.Context(), req, stream)
	if err != nil {
		logger.Debug("could not list tenants")
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (prtSrv *PartitionServer) CreateTenant(ctx context.Context, req *partitionv1.CreateTenantRequest) (*partitionv1.CreateTenantResponse, error) {
	logger := prtSrv.Service.L()
	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.CreateTenant(ctx, req)
	if err != nil {
		logger.Debug("could not create a new tenant")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreateTenantResponse{Data: tenant}, nil
}
