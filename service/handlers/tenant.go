package handlers

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) CreateTenant(ctx context.Context, req *partitionV1.TenantRequest) (*partitionV1.TenantObject, error) {

	tenantBusiness := business.NewTenantBusiness(ctx, prtSrv.Service)
	tenant, err := tenantBusiness.CreateTenant(ctx, req)
	if err != nil {
		log.Printf(" CreateTenant -- could not create a new tenant %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return tenant, nil
}
