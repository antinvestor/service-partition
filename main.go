package main

import (
	"fmt"
	"github.com/antinvestor/apis/go/common"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/handlers"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/queue"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidateinterceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pitabwire/frame"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	serviceName := "service_partition"
	partitionConfig, err := frame.ConfigFromEnv[config.PartitionConfig]()
	if err != nil {
		logrus.WithError(err).Fatal("could not process configs")
		return
	}

	ctx, service := frame.NewService(serviceName, frame.Config(&partitionConfig))
	logger := service.L(ctx)

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
			logging.UnaryServerInterceptor(frame.LoggingInterceptor(logger), frame.GetLoggingOptions()...),
			service.UnaryAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.UnaryServerInterceptor(validator),
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(frame.LoggingInterceptor(logger), frame.GetLoggingOptions()...),
			service.StreamAuthInterceptor(jwtAudience, partitionConfig.Oauth2JwtVerifyIssuer),
			protovalidateinterceptor.StreamServerInterceptor(validator),
			recovery.StreamServerInterceptor(),
		),
	)

	implementation := &handlers.PartitionServer{
		Service: service,
	}

	partitionv1.RegisterPartitionServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	proxyOptions := common.ProxyOptions{
		GrpcServerEndpoint: fmt.Sprintf("localhost:%s", partitionConfig.GrpcServerPort),
		GrpcServerDialOpts: []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}

	proxyMux, err := partitionv1.CreateProxyHandler(ctx, proxyOptions)
	if err != nil {
		logger.WithError(err).Fatal("could not create proxy handler")
		return
	}

	proxyServerOpt := frame.HttpHandler(proxyMux)
	serviceOptions = append(serviceOptions, proxyServerOpt)

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
