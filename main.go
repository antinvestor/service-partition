package main

import (
	partitionv1 "github.com/antinvestor/apis/partition/v1"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/handlers"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/queue"
	"github.com/bufbuild/protovalidate-go"
	protovalidate_interceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	validator, err := protovalidate.New()
	if err != nil {
		logger.WithError(err).Fatal("could not load validator for proto messages")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			service.UnaryAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer),
			recovery.UnaryServerInterceptor(),
			protovalidate_interceptor.UnaryServerInterceptor(validator),
		),
		grpc.ChainStreamInterceptor(
			service.StreamAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer),
			recovery.StreamServerInterceptor(),
			protovalidate_interceptor.StreamServerInterceptor(validator),
		),
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

	if partitionConfig.SynchronizePrimaryPartitions {
		service.AddPreStartMethod(business.ReQueuePrimaryPartitionsForSync)
	}

	logger.WithField("server http port", partitionConfig.HttpServerPort).
		WithField("server grpc port", partitionConfig.GrpcServerPort).
		Info(" Initiating server operations")
	err = implementation.Service.Run(ctx, "")
	if err != nil {
		logger.WithError(err).Fatal("could not run server")
	}

}
