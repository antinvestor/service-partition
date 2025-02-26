package main

import (
	"context"
	"fmt"
	apis "github.com/antinvestor/apis/go/common"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"log"
)

func main() {

	ctx := context.Background()
	partitionCli, err := partitionv1.NewPartitionsClient(ctx,
		apis.WithEndpoint("localhost:50051"),
		apis.WithTokenEndpoint("https://oauth2.chamamobile.com/oauth2/token"),
		apis.WithTokenUsername("dc8e8598-89d5-4983-818e-50fbb6933498"),
		apis.WithTokenPassword("b1QQBzBM04CDvEX9heDx"),
		apis.WithAudiences("service_profile", "service_partition", "service_files", "service_notification", "service_payment", "service_stawi"))
	if err != nil {
		log.Fatalf("Failed to create partition client: %v", err)
	}

	log.Println("Partition client token : ", partitionCli.GetInfo())

	partition, err := partitionCli.GetPartition(ctx, "9bsv0s0hijjg02qks6kg")
	if err != nil {
		log.Fatalf("Failed to get partition: %v", err)
	}

	fmt.Println("Successfully retrieved partition:", partition)
}
