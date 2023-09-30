package config

import (
	"fmt"
	"os"
)

// ConnectionURLBuilder returns connection urls for the given service
func ConnectionURLBuilder(serviceName string) (*string, error) {
	var url string
	switch serviceName {
	case "postgres":
		// URL for PostgreSQL connection.
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv(DBHostEnvVar),
			os.Getenv(DBPortEnvVar),
			os.Getenv(DBUserEnvVar),
			os.Getenv(DBPassWordEnvVar),
			os.Getenv(DBNameEnvVar),
			os.Getenv(DBSSLModeEnvVar),
		)

	case "postgres-migrate":
		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv(DBUserEnvVar),
			os.Getenv(DBPassWordEnvVar),
			os.Getenv(DBHostEnvVar),
			os.Getenv(DBPortEnvVar),
			os.Getenv(DBNameEnvVar),
			os.Getenv(DBSSLModeEnvVar),
		)
	default:
		// Return error message.
		return nil, fmt.Errorf("connection name '%v' is not supported", serviceName)
	}

	// Return connection URL.
	return &url, nil
}
