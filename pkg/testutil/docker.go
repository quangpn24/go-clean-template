package testutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgresContainer(t testing.TB, dbname, user, password string) testcontainers.Container {
	ctx := context.Background()
	postgresql, err := postgres.Run(ctx, "docker.io/postgres:15.2-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)

	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, postgresql.Terminate(ctx))
	})

	return postgresql
}
