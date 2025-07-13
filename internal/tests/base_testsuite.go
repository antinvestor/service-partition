package tests

import (
	"context"
	"fmt"
	"github.com/antinvestor/service-partition/config"
	"github.com/antinvestor/service-partition/service/repository"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net"
	"path/filepath"
	"testing"

	"github.com/pitabwire/frame"
	"github.com/pitabwire/frame/tests"
	"github.com/pitabwire/frame/tests/deps/testpostgres"
	"github.com/pitabwire/frame/tests/testdef"
	"github.com/pitabwire/util"
	"github.com/stretchr/testify/require"
)

const (
	PostgresqlDBImage = "paradedb/paradedb:latest"
	OryHydraImage     = "oryd/hydra:latest"

	DefaultRandomStringLength = 8
)

type BaseTestSuite struct {
	tests.FrameBaseTestSuite
}

func initResources(_ context.Context) []testdef.TestResource {
	pg := testpostgres.NewPGDepWithCred(PostgresqlDBImage, "ant", "s3cr3t", "service_profile")
	resources := []testdef.TestResource{pg}
	return resources
}

func (bs *BaseTestSuite) migrateHydraContainer(ctx context.Context, postgresqlUri, configFilePath string) error {

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:    OryHydraImage,
			Networks: []string{bs.Network.Name},
			Cmd:      []string{"migrate", "sql", "up", "--read-from-env", "--yes"},
			Env: map[string]string{
				"LOG_LEVEL": "debug",
				"DSN":       postgresqlUri,
			},

			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      configFilePath,
					ContainerFilePath: "/etc/config/hydra.yml",
					FileMode:          0o755,
				},
			},
			WaitingFor: wait.ForExit(),
		},

		Started: true,
	})
	if err != nil {
		return err
	}

	err = container.Terminate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BaseTestSuite) CreateHydraContainer(t *testing.T) (string, func(), error) {

	ctx := t.Context()

	var datastoreUri frame.DataSource
	for _, res := range bs.Resources() {
		if res.GetInternalDS().IsPostgres() {
			datastoreUri = res.GetInternalDS()
		}
	}

	configFilePath, err := filepath.Abs("../../internal/tests/hydra.yaml")
	require.NoError(t, err)

	err = bs.migrateHydraContainer(ctx, datastoreUri.String(), configFilePath)
	require.NoError(t, err)

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        OryHydraImage,
			Networks: []string{bs.Network.Name},
			ExposedPorts: []string{"4444/tcp", "4445/tcp"},
			Cmd:          []string{"serve", "all", "--config", "/etc/config/hydra.yml", "--dev"},
			Env: map[string]string{
				"LOG_LEVEL": "debug",
				"DSN":       datastoreUri.String(),
			},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      configFilePath,
					ContainerFilePath: "/etc/config/hydra.yml",
					FileMode:          0o755,
				},
			},
			WaitingFor: wait.ForHTTP("/health/ready").WithPort("4445/tcp"),
		},
		Started: true,
	})


	cleanupFunc := func() {
		if c != nil {
			_ = c.Terminate(ctx)
		}
	}
	if err != nil {
		return "", cleanupFunc, err
	}

	containerPort, err := c.MappedPort(ctx, "4445/tcp")
	if err != nil {
		return "", cleanupFunc, err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return "", cleanupFunc, err
	}

	connStr := fmt.Sprintf("http://%s", net.JoinHostPort(host, containerPort.Port()))
	return connStr,cleanupFunc, nil
}

func (bs *BaseTestSuite) SetupSuite() {
	bs.InitResourceFunc = initResources
	bs.FrameBaseTestSuite.SetupSuite()
}

func (bs *BaseTestSuite) CreateService(
	t *testing.T,
	depOpts *testdef.DependancyOption,
) (*frame.Service, context.Context) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	cfg, err := frame.ConfigFromEnv[config.PartitionConfig]()
	require.NoError(t, err)

	cfg.LogLevel = "debug"
	cfg.RunServiceSecurely = false
	cfg.ServerPort = ""

	for _, res := range depOpts.Database() {
		testDS, cleanup, err0 := res.GetRandomisedDS(t.Context(), depOpts.Prefix())
		require.NoError(t, err0)

		t.Cleanup(func() {
			cleanup(t.Context())
		})

		cfg.DatabasePrimaryURL = []string{testDS.String()}
		cfg.DatabaseReplicaURL = []string{testDS.String()}
	}

	ctx, svc := frame.NewServiceWithContext(t.Context(), "profile tests",
		frame.WithConfig(&cfg),
		frame.WithDatastore(),
		frame.WithNoopDriver())

	svc.Init(ctx)

	err = repository.Migrate(ctx, svc, "../../migrations/0001")
	require.NoError(t, err)

	err = svc.Run(ctx, "")
	require.NoError(t, err)

	return svc, ctx
}

func (bs *BaseTestSuite) TearDownSuite() {
	bs.FrameBaseTestSuite.TearDownSuite()
}

// WithTestDependancies Creates subtests with each known DependancyOption.
func (bs *BaseTestSuite) WithTestDependancies(t *testing.T, testFn func(t *testing.T, dep *testdef.DependancyOption)) {
	options := []*testdef.DependancyOption{
		testdef.NewDependancyOption("default", util.RandomString(DefaultRandomStringLength), bs.Resources()),
	}

	tests.WithTestDependancies(t, options, testFn)
}
