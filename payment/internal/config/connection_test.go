// config_test.go

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionURLBuilder_Postgres(t *testing.T) {
	// Set the required environment variables for PostgreSQL
	os.Setenv(DBHostEnvVar, "localhost")
	os.Setenv(DBPortEnvVar, "5432")
	os.Setenv(DBUserEnvVar, "user")
	os.Setenv(DBPassWordEnvVar, "password")
	os.Setenv(DBNameEnvVar, "dbname")
	os.Setenv(DBSSLModeEnvVar, "disable")

	expectedURL := "host=localhost port=5432 user=user password=password dbname=dbname sslmode=disable"

	url, err := ConnectionURLBuilder("postgres")

	assert.Nil(t, err)
	assert.Equal(t, expectedURL, *url)
}

func TestConnectionURLBuilder_PostgresMigrate(t *testing.T) {
	// Set the required environment variables for PostgreSQL Migrate
	os.Setenv(DBUserEnvVar, "user")
	os.Setenv(DBPassWordEnvVar, "password")
	os.Setenv(DBHostEnvVar, "localhost")
	os.Setenv(DBPortEnvVar, "5432")
	os.Setenv(DBNameEnvVar, "dbname")
	os.Setenv(DBSSLModeEnvVar, "disable")

	expectedURL := "postgres://user:password@localhost:5432/dbname?sslmode=disable"

	url, err := ConnectionURLBuilder("postgres-migrate")

	assert.Nil(t, err)
	assert.Equal(t, expectedURL, *url)
}

func TestConnectionURLBuilder_UnsupportedService(t *testing.T) {
	// Test unsupported service name
	url, err := ConnectionURLBuilder("unknown-service")

	assert.NotNil(t, err)
	assert.Nil(t, url)
	assert.Equal(t, "connection name 'unknown-service' is not supported", err.Error())
}
