package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/talvor/asyncapi/config"
	"github.com/talvor/asyncapi/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Config           *config.Config
	DB               *sql.DB
	ContainerCleanup func()
}

type ContainerCleanup func()

func NewTestEnv() (*TestEnv, error) {
	os.Setenv("ENV", string(config.EnvTest))

	conf := config.GetConfig()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":     conf.DatabaseUser,
			"POSTGRES_PASSWORD": conf.DatabasePassword,
			"POSTGRES_DB":       conf.DatabaseName,
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     testcontainers.ProviderPodman,
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		container.Terminate(context.Background())
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	conf.SetDatabasePort(mappedPort.Port())

	db, err := store.NewPostgresDB(conf)
	if err != nil {
		return nil, err
	}

	return &TestEnv{
		Config:           conf,
		DB:               db,
		ContainerCleanup: cleanup,
	}, nil
}

func (te *TestEnv) SetupDB() error {
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", te.Config.ProjectRoot),
		te.Config.DatabaseURL(),
	)
	if err != nil {
		return fmt.Errorf("failed to migrate db %w", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate db %w", err)
	}
	return nil
}

func (te *TestEnv) TeardownDB() error {
	_, err := te.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ",")))
	if err != nil {
		return fmt.Errorf("failed to cleanup db %w", err)
	}
	return nil
}
