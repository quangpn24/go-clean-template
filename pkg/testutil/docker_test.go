package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupPostgresContainer(t *testing.T) {
	//Arrange
	dbName, dbUser, dbPass := "test1", "user", "123456"

	//Act
	cont := SetupPostgresContainer(t, dbName, dbUser, dbPass)

	//Assert
	assert.NotNil(t, cont)
}
