package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) CreateAccess(ctx context.Context, req *partitionv1.AccessCreateRequest) (*partitionv1.AccessObject, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	access, err := accessBusiness.CreateAccess(ctx, req)
	if err != nil {
		log.Printf(" CreateAccess -- could not create new access %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return access, nil
}
func (prtSrv *PartitionServer) GetAccess(ctx context.Context, req *partitionv1.AccessGetRequest) (*partitionv1.AccessObject, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	access, err := accessBusiness.GetAccess(ctx, req)
	if err != nil {
		log.Printf(" GetAccess -- could not get access %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return access, nil
}
func (prtSrv *PartitionServer) RemoveAccess(ctx context.Context, req *partitionv1.AccessRemoveRequest) (*partitionv1.RemoveResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	err := accessBusiness.RemoveAccess(ctx, req)
	if err != nil {
		log.Printf(" RemoveAccess -- could not remove access %+v", err)
		return &partitionv1.RemoveResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveResponse{
		Succeeded: true,
	}, nil
}
func (prtSrv *PartitionServer) CreateAccessRole(ctx context.Context, req *partitionv1.AccessRoleCreateRequest) (*partitionv1.AccessRoleObject, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	accessRole, err := accessBusiness.CreateAccessRole(ctx, req)
	if err != nil {
		log.Printf(" CreateAccessRole -- could not create new access roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return accessRole, nil
}
func (prtSrv *PartitionServer) ListAccessRoles(ctx context.Context, req *partitionv1.AccessRoleListRequest) (*partitionv1.AccessRoleListResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	accessRoleList, err := accessBusiness.ListAccessRoles(ctx, req)
	if err != nil {
		log.Printf(" ListAccessRoles -- could not get list of access roles %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return accessRoleList, nil
}
func (prtSrv *PartitionServer) RemoveAccessRole(ctx context.Context, req *partitionv1.AccessRoleRemoveRequest) (*partitionv1.RemoveResponse, error) {
	accessBusiness := business.NewAccessBusiness(ctx, prtSrv.Service)
	err := accessBusiness.RemoveAccessRole(ctx, req)
	if err != nil {
		log.Printf(" RemoveAccessRole -- could not remove access role %+v", err)
		return &partitionv1.RemoveResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveResponse{
		Succeeded: true,
	}, nil
}
