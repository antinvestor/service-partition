package handlers

import (
	"context"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type PartitionServer struct {
	Service *frame.Service
	partitionV1.UnimplementedPartitionServiceServer
}

func (prtSrv *PartitionServer) CreatePartition(ctx context.Context, req *partitionV1.PartitionCreateRequest) (*partitionV1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	partition, err := partitionBusiness.CreatePartition(ctx, req)
	if err != nil {
		log.Printf(" CreatePartition -- could not create a new partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) GetPartition(ctx context.Context, req *partitionV1.PartitionGetRequest) (*partitionV1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	partition, err := partitionBusiness.GetPartition(ctx, req)
	if err != nil {
		log.Printf(" GetPartition -- could not obtain the specified partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}
func (prtSrv *PartitionServer) UpdatePartition(ctx context.Context, req *partitionV1.PartitionUpdateRequest) (*partitionV1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	partition, err := partitionBusiness.UpdatePartition(ctx, req)
	if err != nil {
		log.Printf(" UpdatePartition -- could not update existing partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}
func (prtSrv *PartitionServer) CreatePartitionRole(ctx context.Context, req *partitionV1.PartitionRoleCreateRequest) (*partitionV1.PartitionRoleObject, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	partition, err := partitionBusiness.CreatePartitionRole(ctx, req)
	if err != nil {
		log.Printf(" CreatePartitionRole -- could not create a new partition role %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}
func (prtSrv *PartitionServer) ListPartitionRoles(ctx context.Context, req *partitionV1.PartitionRoleListRequest) (*partitionV1.PartitionRoleListResponse, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	partition, err := partitionBusiness.ListPartitionRoles(ctx, req)
	if err != nil {
		log.Printf(" ListPartitionRoles -- could not obtain the list of partition roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}
func (prtSrv *PartitionServer) RemovePartitionRole(ctx context.Context, req *partitionV1.PartitionRoleRemoveRequest) (*partitionV1.RemoveResponse, error) {
	partitionBusiness := business.NewPartitionBusiness(ctx, prtSrv.Service)
	err := partitionBusiness.RemovePartition(ctx, req)
	if err != nil {
		log.Printf(" RemovePartitionRole -- could not remove the specified partition role %+v", err)
		return &partitionV1.RemoveResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionV1.RemoveResponse{
		Succeeded: true,
	}, nil
}


