package main

import (
	"fmt"
	partitionv1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/handlers"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/queue"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
)

func main() {

	serviceName := "service_partition"
	var partitionConfig config.PartitionConfig
	err := frame.ConfigProcess("", &partitionConfig)
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.Config(&partitionConfig))
	logger := service.L()

	serviceOptions := []frame.Option{frame.Datastore(ctx)}

	if partitionConfig.DoDatabaseMigrate() {

		service.Init(serviceOptions...)
		err := service.MigrateDatastore(ctx, partitionConfig.GetDatabaseMigrationPath(),
			models.Tenant{}, models.Partition{}, models.PartitionRole{},
			models.Access{}, models.AccessRole{}, models.Page{})

		if err != nil {
			logger.WithError(err).Fatalf("could not migrate successfully")
		}
		return
	}

	err = service.RegisterForJwt(ctx)
	if err != nil {
		logger.WithError(err).Fatal("could not register for jwt")
		return
	}

	jwtAudience := partitionConfig.Oauth2JwtVerifyAudience
	if jwtAudience == "" {
		jwtAudience = serviceName
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			service.UnaryAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer),
		)),
		grpc.StreamInterceptor(service.StreamAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer)),
	)

	implementation := &handlers.PartitionServer{
		Service: service,
	}

	partitionv1.RegisterPartitionServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	partitionSyncQueueHandler := queue.PartitionSyncQueueHandler{
		Service: service,
	}
	partitionSyncQueueURL := partitionConfig.QueuePartitionSyncURL
	partitionSyncQueue := frame.RegisterSubscriber(partitionConfig.PartitionSyncName, partitionSyncQueueURL, 2, &partitionSyncQueueHandler)
	partitionSyncQueueP := frame.RegisterPublisher(partitionConfig.PartitionSyncName, partitionSyncQueueURL)

	serviceOptions = append(serviceOptions, partitionSyncQueue, partitionSyncQueueP)

	service.Init(serviceOptions...)

	serverPort := partitionConfig.ServerPort

	if partitionConfig.SynchronizePrimaryPartitions {
		service.AddPreStartMethod(business.ReQueuePrimaryPartitionsForSync)
	}

	logger.WithField("server port", serverPort).Info(" Initiating server operations")
	err = implementation.Service.Run(ctx, fmt.Sprintf(":%v", serverPort))
	if err != nil {
		logger.WithError(err).Fatal("could not run server")
	}

}
