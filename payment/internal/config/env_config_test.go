// config_test.go

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAllEnvVariables_AllSet(t *testing.T) {
	// Set all required environment variables
	requiredEnvVars := []string{
		RuntimeStageEnvVar,
		DBPortEnvVar,
		DBHostEnvVar,
		DBUserEnvVar,
		DBPassWordEnvVar,
		DBNameEnvVar,
		DBSSLModeEnvVar,
		DBMaxConnectionsEnvVar,
		DBMaxIdleConnectionsEnvVar,
		DBMaxLifetimeConnectionsEnvVar,
		RunMigrationsEnvVar,
		GRPCListenAddressEnvVar,
		HTTPListenAddressEnvVar,
		MpesaConsumerKeyEnvVar,
		MpesaConsumerSecreteEnvVar,
		MpesaPassKeyEnv,
		OrderServiceListenAddressEnvVar,
	}

	for _, v := range requiredEnvVars {
		defer os.Unsetenv(v)  // Clear the environment variable after the test
		os.Setenv(v, "value") // Set a dummy value for the required env variable
	}

	result := HasAllEnvVariables()
	assert.True(t, result, "Expected HasAllEnvVariables to return true")
}

func TestHasAllEnvVariables_MissingVariable(t *testing.T) {
	// Set some required environment variables but leave one missing
	requiredEnvVars := []string{
		RuntimeStageEnvVar,
		DBPortEnvVar,
		DBHostEnvVar,
		DBUserEnvVar,
		DBPassWordEnvVar,
		DBNameEnvVar,
		DBSSLModeEnvVar,
		DBMaxConnectionsEnvVar,
		DBMaxIdleConnectionsEnvVar,
		DBMaxLifetimeConnectionsEnvVar,
		RunMigrationsEnvVar,
		GRPCListenAddressEnvVar,
		HTTPListenAddressEnvVar,
		MpesaConsumerKeyEnvVar,
		MpesaConsumerSecreteEnvVar,
		// Missing: MpesaPassKeyEnv
		OrderServiceListenAddressEnvVar,
	}

	for _, v := range requiredEnvVars {
		defer os.Unsetenv(v)  // Clear the environment variable after the test
		os.Setenv(v, "value") // Set a dummy value for the required env variables
	}

	result := HasAllEnvVariables()
	assert.False(t, result, "Expected HasAllEnvVariables to return false")
}
