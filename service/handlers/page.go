package handlers

import (
	"context"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/service/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (prtSrv *PartitionServer) CreatePage(ctx context.Context, req *partitionv1.CreatePageRequest) (*partitionv1.CreatePageResponse, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	page, err := pageBusiness.CreatePage(ctx, req)
	if err != nil {
		log.Printf(" CreatePage -- could not create a new page %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.CreatePageResponse{Data: page}, nil
}

func (prtSrv *PartitionServer) GetPage(ctx context.Context, req *partitionv1.GetPageRequest) (*partitionv1.GetPageResponse, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	page, err := pageBusiness.GetPage(ctx, req)
	if err != nil {
		log.Printf(" GetPage -- could not get page %+v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.GetPageResponse{Data: page}, nil
}

func (prtSrv *PartitionServer) RemovePage(ctx context.Context, req *partitionv1.RemovePageRequest) (*partitionv1.RemovePageResponse, error) {
	pageBusiness := business.NewPageBusiness(ctx, prtSrv.Service)
	err := pageBusiness.RemovePage(ctx, req)
	if err != nil {
		log.Printf(" RemovePage -- could not remove page %+v", err)
		return &partitionv1.RemovePageResponse{
			Succeeded: false,
		}, status.Errorf(codes.Internal, err.Error())
	}
	return &partitionv1.RemovePageResponse{
		Succeeded: true,
	}, nil
}
