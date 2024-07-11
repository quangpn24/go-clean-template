package testutil

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MigrateTestDatabase(t testing.TB, db *gorm.DB, migrationPath string) {
	t.Helper()

	migrations := &migrate.FileMigrationSource{
		Dir: migrationPath,
	}

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	_, err = migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	assert.NoError(t, err)
}

func CreateConnection(t testing.TB, dbName string, dbUser string, dbPass string) *gorm.DB {
	cont := SetupPostgresContainer(t, dbName, dbUser, dbPass)
	host, _ := cont.Host(context.Background())
	port, _ := cont.MappedPort(context.Background(), "5432")

	dsnStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host,
		dbUser,
		dbPass,
		dbName,
		port.Port(),
	)

	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{TranslateError: true})

	assert.NoError(t, err)
	assert.NotNil(t, db)

	return db
}
