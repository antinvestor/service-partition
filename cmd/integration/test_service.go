package main

import (
	"context"
	"fmt"
	apis "github.com/antinvestor/apis/go/common"
	partitionv1 "github.com/antinvestor/apis/go/partition/v1"
	"log"
)

// Unary interceptor to log metadata and request body

func main() {

	ctx := context.Background()

	//endpointValues := url.Values{}
	//audienceList := strings.Join([]string{"service_profile", "service_partition", "service_files", "service_notification", "service_payment", "service_stawi"}, " ")
	//endpointValues.Add("audience", audienceList)
	//
	//tokenClient := &clientcredentials.Config{
	//	ClientID:       "1bb24652-3708-4739-ba38-cca42048a1be",
	//	ClientSecret:   "b1QQBzBM04CDvEX9heDx",
	//	TokenURL:       "https://oauth2.chamamobile.com/oauth2/token",
	//	Scopes:         []string{},
	//	EndpointParams: endpointValues,
	//	AuthStyle:      oauth2.AuthStyleAutoDetect,
	//}
	//
	//token, err := tokenClient.Token(ctx)
	//if err != nil {
	//	log.Fatalf("Failed to create token for client: %v", err)
	//}
	//
	//log.Printf("Got token : %s", token.AccessToken)

	partitionCli, err := partitionv1.NewPartitionsClient(ctx,
		apis.WithEndpoint("localhost:50051"),
		apis.WithTokenEndpoint("https://oauth2.chamamobile.com/oauth2/token"),
		apis.WithTokenUsername("ddcc5752-e9bf-4fe0-a73a-669398fe218e"),
		apis.WithTokenPassword("b1QQBzBM04CDvEX9heDx"),
		apis.WithAudiences("service_profile", "service_partition", "service_files", "service_notification", "service_payment", "service_stawi"),
	)

	if err != nil {
		log.Fatalf("Failed to create partition client: %v", err)
	}

	partition, err := partitionCli.GetPartition(ctx, "9bsv0s0hijjg02qks6kg")
	if err != nil {
		log.Fatalf("Failed to get partition: %v", err)
	}

	fmt.Println("Successfully retrieved partition:", partition)
}
