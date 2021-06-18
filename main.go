package main

import (
	"context"
	"fmt"
	partitionV1 "github.com/antinvestor/service-partition-api"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/business"
	"github.com/antinvestor/service-partition/service/handlers"
	"github.com/antinvestor/service-partition/service/models"
	"github.com/antinvestor/service-partition/service/queue"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"

	"github.com/pitabwire/frame"

	"github.com/grpc-ecosystem/go-grpc-middleware"
)

func main() {

	serviceName := "service_partition"
	ctx := context.Background()

	datasource := frame.GetEnv(config.EnvDatabaseUrl, "postgres://ant:@nt@localhost/service_partition")
	mainDb := frame.Datastore(ctx, datasource, false)

	readOnlydatasource := frame.GetEnv(config.EnvReplicaDatabaseUrl, datasource)
	readDb := frame.Datastore(ctx, readOnlydatasource, true)

	service := frame.NewService(serviceName, mainDb, readDb)

	isMigration, err := strconv.ParseBool(frame.GetEnv(config.EnvMigrate, "false"))
	if err != nil {
		isMigration = false
	}

	stdArgs := os.Args[1:]
	if (len(stdArgs) > 0 && stdArgs[0] == "migrate") || isMigration {

		migrationPath := frame.GetEnv(config.EnvMigrationPath, "./migrations/0001")
		err := service.MigrateDatastore(ctx, migrationPath,
			models.Tenant{}, models.Partition{}, models.PartitionRole{},
			models.Access{}, models.AccessRole{}, models.Page{})

		if err != nil {
			log.Fatalf("main -- Could not migrate successfully because : %v", err)
		}
		return

	}

	var serviceOptions []frame.Option

	jwtAudience := frame.GetEnv(config.EnvOauth2JwtVerifyAudience, serviceName)
	jwtIssuer := frame.GetEnv(config.EnvOauth2JwtVerifyIssuer, "")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			frame.UnaryAuthInterceptor(jwtAudience, jwtIssuer),
		)),
	)

	implementation := &handlers.PartitionServer{
		Service: service,
	}

	partitionV1.RegisterPartitionServiceServer(grpcServer, implementation)

	grpcServerOpt := frame.GrpcServer(grpcServer)
	serviceOptions = append(serviceOptions, grpcServerOpt)

	partitionSyncQueueHandler := queue.PartitionSyncQueueHandler{
		Service: service,
	}
	partitionSyncQueueURL := frame.GetEnv(config.EnvQueuePartitionSync, fmt.Sprintf("mem://%s", config.QueuePartitionSyncName))
	partitionSyncQueue := frame.RegisterSubscriber(config.QueuePartitionSyncName, partitionSyncQueueURL, 2, &partitionSyncQueueHandler)
	partitionSyncQueueP := frame.RegisterPublisher(config.QueuePartitionSyncName, partitionSyncQueueURL)

	serviceOptions = append(serviceOptions, partitionSyncQueue, partitionSyncQueueP)

	service.Init(serviceOptions...)

	serverPort := frame.GetEnv(config.EnvServerPort, "7003")

	service.AddPreStartMethod(business.ReQueuePrimaryPartitionsForSync)

	log.Printf(" main -- Initiating server operations on : %s", serverPort)
	err = implementation.Service.Run(ctx, fmt.Sprintf(":%v", serverPort))
	if err != nil {
		log.Fatalf("main -- Could not run Server : %v", err)
	}

}
