package main

import (
	"buf.build/go/protovalidate"
	"context"
	"fmt"
	"github.com/antinvestor/apis/go/common"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/handlers"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/queue"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/pitabwire/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	serviceName := "service_partition"
	ctx := context.Background()

	cfg, err := frame.ConfigLoadWithOIDC[config.PartitionConfig](ctx)
	if err != nil {
		util.Log(ctx).WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewServiceWithContext(ctx, serviceName, frame.WithConfig(&cfg))
	logger := service.Log(ctx)

	serviceOptions := []frame.Option{frame.WithDatastore()}

	if cfg.DoDatabaseMigrate() {

		service.Init(ctx, serviceOptions...)
		err = service.MigrateDatastore(ctx, cfg.GetDatabaseMigrationPath(),
			models.Tenant{}, models.Partition{}, models.PartitionRole{},
			models.Access{}, models.AccessRole{}, models.Page{})

		if err != nil {
			logger.WithError(err).Fatal("could not migrate successfully")
		}
		return
	}

	err = service.RegisterForJwt(ctx)
	if err != nil {
		logger.WithError(err).Fatal("could not register for jwt")
		return
	}

	jwtAudience := cfg.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	validator, err := protovalidate.New()
	if err != nil {
		logger.WithError(err).Fatal("could not load validator for proto messages")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			service.UnaryAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator),
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			service.StreamAuthInterceptor(jwtAudience, cfg.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
			recovery.StreamServerInterceptor(),
		),
	)

	implementation := &handlers.PartitionServer{
		Service: service,
	}

	partitionv1.RegisterPartitionServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.WithGRPCServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	proxyOptions := common.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", cfg.GrpcServerPort),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := partitionv1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		logger.WithError(err).Fatal("could not create proxy handler")
		return
	}

	proxyServerOpt := frame.WithHTTPHandler(proxyMux)
	serviceOptions = append(serviceOptions, proxyServerOpt)

	partitionSyncQueueHandler := queue.PartitionSyncQueueHandler{
		Service: service,
	}
	partitionSyncQueueURL := cfg.QueuePartitionSyncURL
	partitionSyncQueue := frame.WithRegisterSubscriber(cfg.PartitionSyncName, partitionSyncQueueURL,  &partitionSyncQueueHandler)
	partitionSyncQueueP := frame.WithRegisterPublisher(cfg.PartitionSyncName, partitionSyncQueueURL)

	serviceOptions = append(serviceOptions, partitionSyncQueue, partitionSyncQueueP)

	service.Init(ctx, serviceOptions...)

	if cfg.SynchronizePrimaryPartitions {
		service.AddPreStartMethod(business.ReQueuePrimaryPartitionsForSync)
	}

	logger.WithField("server http port", cfg.HTTPServerPort).
		WithField("server grpc port", cfg.GrpcServerPort).
		Info(" Initiating server operations")
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		logger.WithError(err).Fatal("could not run server")
	}

}
