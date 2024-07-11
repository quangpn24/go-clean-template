package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePostgresConnect(t *testing.T) {
	testDBName := "db-name"
	testUser := "user"
	testPassword := "123456"

	db := CreateConnection(t, testDBName, testUser, testPassword)

	assert.NotNil(t, db)
}

func TestMigrateTestDatabase(t *testing.T) {
	// Arrange
	testDBName := "db-name"
	testUser := "user"
	testPassword := "123456"
	db := CreateConnection(t, testDBName, testUser, testPassword)
	migrationPath := "../../migrations"

	// Act
	MigrateTestDatabase(t, db, migrationPath)

	// Assert
	sqlDB, err := db.DB()
	assert.NoError(t, err)

	// Check that the migration table exists
	var count int
	err = sqlDB.QueryRow("SELECT count(*) FROM information_schema.tables WHERE table_name = 'gorp_migrations'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
