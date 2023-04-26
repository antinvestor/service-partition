package testsutil

import (
	"context"
	"github.com/antinvestor/service-partition/config"
	"github.com/pitabwire/frame"
)

func GetTestService(name string) (context.Context, *frame.Service, error) {
	dbURL := frame.GetEnv("TEST_DATABASE_URL",
		"postgres://ant:secret@localhost:5423/service_partition?sslmode=disable")
	mainDB := frame.DatastoreCon(dbURL, false)

	var partitionConfig config.PartitionConfig
	err := frame.ConfigProcess("", &partitionConfig)
	if err != nil {
		return nil, nil, err
	}
	partitionConfig.Oauth2ServiceAdminURI = "http://localhost:4445"

	ctx, service := frame.NewService(name, frame.Config(&partitionConfig), mainDB, frame.NoopDriver())
	service.Init()

	return ctx, service, nil
}
