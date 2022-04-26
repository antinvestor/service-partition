package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type PartitionServer struct {
	Service *frame.Service
	partitionv1.UnimplementedPartitionServiceServer
}

func (prtSrv *PartitionServer) ListPartition(
	req *partitionv1.SearchRequest,
	stream partitionv1.PartitionService_ListPartitionServer) error {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	err := partitionBusiness.ListPartition(stream.Context(), req, stream)
	if err != nil {
		log.Printf(" ListPartition -- could not list partition %+v", err)
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (prtSrv *PartitionServer) CreatePartition(
	ctx context.Context,
	req *partitionv1.PartitionCreateRequest) (*partitionv1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.CreatePartition(ctx, req)
	if err != nil {
		log.Printf(" CreatePartition -- could not create a new partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) GetPartition(
	ctx context.Context,
	req *partitionv1.PartitionGetRequest) (*partitionv1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.GetPartition(ctx, req)
	if err != nil {
		log.Printf(" GetPartition -- could not obtain the specified partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) UpdatePartition(
	ctx context.Context,
	req *partitionv1.PartitionUpdateRequest) (*partitionv1.PartitionObject, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.UpdatePartition(ctx, req)
	if err != nil {
		log.Printf(" UpdatePartition -- could not update existing partition %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) CreatePartitionRole(
	ctx context.Context,
	req *partitionv1.PartitionRoleCreateRequest) (*partitionv1.PartitionRoleObject, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.CreatePartitionRole(ctx, req)
	if err != nil {
		log.Printf(" CreatePartitionRole -- could not create a new partition role %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) ListPartitionRoles(
	ctx context.Context,
	req *partitionv1.PartitionRoleListRequest) (*partitionv1.PartitionRoleListResponse, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.ListPartitionRoles(ctx, req)
	if err != nil {
		log.Printf(" ListPartitionRoles -- could not obtain the list of partition roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) RemovePartitionRole(
	ctx context.Context,
	req *partitionv1.PartitionRoleRemoveRequest) (*partitionv1.RemoveResponse, error) {
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	err := partitionBusiness.RemovePartition(ctx, req)
	if err != nil {
		log.Printf(" RemovePartitionRole -- could not remove the specified partition role %+v", err)
		return &partitionv1.RemoveResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveResponse{
		Succeeded: true,
	}, nil
}
