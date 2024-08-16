package postgrestore

import (
	"context"
	"fmt"
	"testing"

	"go-clean-template/pkg/config"
	"go-clean-template/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	// Arrange
	testDBName := "test-db"
	testDBHost := "localhost"
	testDbPort := 5432
	testDBUsername := "test-user"
	testDBPassword := "P@ssw0rd"
	cfg := &config.Config{
		DB: struct {
			Name      string `envconfig:"DB_NAME"`
			Host      string `envconfig:"DB_HOST"`
			Port      int    `envconfig:"DB_PORT"`
			User      string `envconfig:"DB_USER"`
			Pass      string `envconfig:"DB_PASS"`
			EnableSSL bool   `envconfig:"ENABLE_SSL"`
		}{
			Name:      testDBName,
			Host:      testDBHost,
			Port:      testDbPort,
			User:      testDBUsername,
			Pass:      testDBPassword,
			EnableSSL: false,
		},
	}

	// Act
	actual := ParseFromConfig(cfg)

	// Assert
	expected := Options{
		DBName:   testDBName,
		DBUser:   testDBUsername,
		Password: testDBPassword,
		Host:     testDBHost,
		Port:     fmt.Sprintf("%d", testDbPort),
		SSLMode:  false,
	}

	assert.Equal(t, expected, actual)
}

func TestNewDB(t *testing.T) {
	t.Run("good case", func(t *testing.T) {
		// Arrange
		dbName, dbUser, dbPass := "test1", "test1", "123456"
		cont := testutil.SetupPostgresContainer(t, dbName, dbUser, dbPass)
		host, _ := cont.Host(context.Background())
		port, _ := cont.MappedPort(context.Background(), "5432")

		// Act
		db, err := NewDB(Options{
			DBName:   dbName,
			DBUser:   dbUser,
			Password: dbPass,
			Host:     host,
			Port:     port.Port(),
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("bad case", func(t *testing.T) {
		// Arrange
		dbName, dbUser, dbPass := "test1", "test1", "12345"
		cont := testutil.SetupPostgresContainer(t, dbName, dbUser, dbPass)
		host, _ := cont.Host(context.Background())
		port, _ := cont.MappedPort(context.Background(), "5432")
		opt := Options{
			DBName:   dbName,
			DBUser:   dbUser,
			Password: "wrong password",
			Host:     host,
			Port:     port.Port(),
		}

		// Act
		db, err := NewDB(opt)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
