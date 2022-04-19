package handlers

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) ListTenant(req *partitionV1.SearchRequest, stream partitionV1.PartitionService_ListTenantServer) error {
	tenantBusiness := business.NewTenantBusiness(stream.Context(), prtSrv.Service)
	err := tenantBusiness.ListTenant(stream.Context(), req, stream)
	if err != nil {
		log.Printf(" ListTenant -- could not list tenants %+v", err)
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (prtSrv *PartitionServer) CreateTenant(ctx context.Context, req *partitionV1.TenantRequest) (*partitionV1.TenantObject, error) {

	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.CreateTenant(ctx, req)
	if err != nil {
		log.Printf(" CreateTenant -- could not create a new tenant %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return tenant, nil
}
