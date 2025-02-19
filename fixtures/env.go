package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/talvor/asyncapi/config"
	"github.com/talvor/asyncapi/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestEnv struct {
	Config *config.Config
	DB     *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {
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

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ProviderType:     testcontainers.ProviderPodman,
		ContainerRequest: req,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, pgContainer)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	conf.SetDatabasePort(mappedPort.Port())

	db, err := store.NewPostgresDB(conf)
	require.NoError(t, err)

	return &TestEnv{
		Config: conf,
		DB:     db,
	}
}

func (te *TestEnv) SetupDB(t *testing.T) func(t *testing.T) {
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", te.Config.ProjectRoot),
		te.Config.DatabaseURL(),
	)
	require.NoError(t, err)

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TeardownDB
}

func (te *TestEnv) TeardownDB(t *testing.T) {
	_, err := te.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ",")))
	require.NoError(t, err)
}
