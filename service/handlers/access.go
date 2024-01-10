package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) CreateAccess(ctx context.Context, req *partitionv1.CreateAccessRequest) (*partitionv1.CreateAccessResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	access, err := accessBusiness.CreateAccess(ctx, req)
	if err != nil {
		log.Printf(" CreateAccess -- could not create new access %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreateAccessResponse{Data: access}, nil
}
func (prtSrv *PartitionServer) GetAccess(ctx context.Context, req *partitionv1.GetAccessRequest) (*partitionv1.GetAccessResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	access, err := accessBusiness.GetAccess(ctx, req)
	if err != nil {
		log.Printf(" GetAccess -- could not get access %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.GetAccessResponse{Data: access}, nil
}
func (prtSrv *PartitionServer) RemoveAccess(ctx context.Context, req *partitionv1.RemoveAccessRequest) (*partitionv1.RemoveAccessResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	err := accessBusiness.RemoveAccess(ctx, req)
	if err != nil {
		log.Printf(" RemoveAccess -- could not remove access %+v", err)
		return &partitionv1.RemoveAccessResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveAccessResponse{
		Succeeded: true,
	}, nil
}
func (prtSrv *PartitionServer) CreateAccessRole(ctx context.Context, req *partitionv1.CreateAccessRoleRequest) (*partitionv1.CreateAccessRoleResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	accessRole, err := accessBusiness.CreateAccessRole(ctx, req)
	if err != nil {
		log.Printf(" CreateAccessRole -- could not create new access roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreateAccessRoleResponse{Data: accessRole}, nil
}
func (prtSrv *PartitionServer) ListAccessRoles(ctx context.Context, req *partitionv1.ListAccessRoleRequest) (*partitionv1.ListAccessRoleResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	accessRoleList, err := accessBusiness.ListAccessRoles(ctx, req)
	if err != nil {
		log.Printf(" ListAccessRoles -- could not get list of access roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return accessRoleList, nil
}
func (prtSrv *PartitionServer) RemoveAccessRole(ctx context.Context, req *partitionv1.RemoveAccessRoleRequest) (*partitionv1.RemoveAccessRoleResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	err := accessBusiness.RemoveAccessRole(ctx, req)
	if err != nil {
		log.Printf(" RemoveAccessRole -- could not remove access role %+v", err)
		return &partitionv1.RemoveAccessRoleResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveAccessRoleResponse{
		Succeeded: true,
	}, nil
}
