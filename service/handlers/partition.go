package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/pitabwire/frame"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PartitionServer struct {
	Service *frame.Service
	partitionv1.UnimplementedPartitionServiceServer
}

func (prtSrv *PartitionServer) ListPartition(
	req *partitionv1.ListPartitionRequest,
	stream partitionv1.PartitionService_ListPartitionServer) error {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	err := partitionBusiness.ListPartition(stream.Context(), req, stream)
	if err != nil {
		logger.WithError(err).Debug(" could not list partition")
		return status.Errorf(codes.Internal, err.Error())
	}
	return nil
}

func (prtSrv *PartitionServer) CreatePartition(
	ctx context.Context,
	req *partitionv1.CreatePartitionRequest) (*partitionv1.CreatePartitionResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.CreatePartition(ctx, req)
	if err != nil {
		logger.WithError(err).Debug(" could not create a new partition")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreatePartitionResponse{Data: partition}, nil
}

func (prtSrv *PartitionServer) GetPartition(
	ctx context.Context,
	req *partitionv1.GetPartitionRequest) (*partitionv1.GetPartitionResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.GetPartition(ctx, req)
	if err != nil {
		logger.WithError(err).Debug(" could not obtain the specified partition")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.GetPartitionResponse{Data: partition}, nil
}

func (prtSrv *PartitionServer) UpdatePartition(
	ctx context.Context,
	req *partitionv1.UpdatePartitionRequest) (*partitionv1.UpdatePartitionResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.UpdatePartition(ctx, req)
	if err != nil {
		logger.WithError(err).Debug(" could not update existing partition")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.UpdatePartitionResponse{Data: partition}, nil
}

func (prtSrv *PartitionServer) CreatePartitionRole(
	ctx context.Context,
	req *partitionv1.CreatePartitionRoleRequest) (*partitionv1.CreatePartitionRoleResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.CreatePartitionRole(ctx, req)
	if err != nil {
		logger.WithError(err).Debug("could not create a new partition role")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreatePartitionRoleResponse{Data: partition}, nil
}

func (prtSrv *PartitionServer) ListPartitionRoles(
	ctx context.Context,
	req *partitionv1.ListPartitionRoleRequest) (*partitionv1.ListPartitionRoleResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	partition, err := partitionBusiness.ListPartitionRoles(ctx, req)
	if err != nil {
		logger.WithError(err).Debug(" could not obtain the list of partition roles")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return partition, nil
}

func (prtSrv *PartitionServer) RemovePartitionRole(
	ctx context.Context,
	req *partitionv1.RemovePartitionRoleRequest) (*partitionv1.RemovePartitionRoleResponse, error) {
	logger := prtSrv.Service.L()
	partitionBusiness := business.NewPartitionBusiness(prtSrv.Service)
	err := partitionBusiness.RemovePartitionRole(ctx, req)
	if err != nil {
		logger.WithError(err).Debug(" could not remove the specified partition role")
		return &partitionv1.RemovePartitionRoleResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemovePartitionRoleResponse{
		Succeeded: true,
	}, nil
}
