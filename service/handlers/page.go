package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) CreatePage(ctx context.Context, req *partitionv1.PageCreateRequest) (*partitionv1.PageObject, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	page, err := pageBusiness.CreatePage(ctx, req)
	if err != nil {
		log.Printf(" CreatePage -- could not create a new page %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return page, nil
}

func (prtSrv *PartitionServer) GetPage(ctx context.Context, req *partitionv1.PageGetRequest) (*partitionv1.PageObject, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	page, err := pageBusiness.GetPage(ctx, req)
	if err != nil {
		log.Printf(" GetPage -- could not get page %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return page, nil
}

func (prtSrv *PartitionServer) RemovePage(ctx context.Context, req *partitionv1.PageRemoveRequest) (*partitionv1.RemoveResponse, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	err := pageBusiness.RemovePage(ctx, req)
	if err != nil {
		log.Printf(" RemovePage -- could not remove page %+v", err)
		return &partitionv1.RemoveResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemoveResponse{
		Succeeded: true,
	}, nil
}
