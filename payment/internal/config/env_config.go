package config

import (
	"os"

	"golang.org/x/exp/slog"
)

const (
	RuntimeStageEnvVar              = "RUNTIME_STAGE"
	DBPortEnvVar                    = "DB_PORT"
	DBHostEnvVar                    = "DB_HOST"
	DBUserEnvVar                    = "DB_USER"
	DBPassWordEnvVar                = "DB_PASSWORD"
	DBNameEnvVar                    = "DB_NAME"
	DBSSLModeEnvVar                 = "DB_SSL_MODE"
	DBMaxConnectionsEnvVar          = "DB_MAX_CONNECTIONS"
	DBMaxIdleConnectionsEnvVar      = "DB_MAX_IDLE_CONNECTIONS"
	DBMaxLifetimeConnectionsEnvVar  = "DB_MAX_LIFETIME_CONNECTIONS"
	RunMigrationsEnvVar             = "RUN_MIGRATIONS"
	GRPCListenAddressEnvVar         = "GRPC_LISTEN_ADDRESS"
	HTTPListenAddressEnvVar         = "HTTP_LISTEN_ADDRESS"
	MpesaConsumerKeyEnvVar          = "MPESA_CONSUMER_KEY"
	MpesaConsumerSecreteEnvVar      = "MPESA_CONSUMER_SECRETE"
	MpesaPassKeyEnv                 = "MPESA_PASSKEY"
	CallBackBaseURL                 = "CALLBACK_BASEURL"
	OrderServiceListenAddressEnvVar = "ORDER_SERVICE_LISTEN_ADDRESS"
)

func HasAllEnvVariables() bool {
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
		value, ok := os.LookupEnv(v)
		if !ok || len(value) < 1 {
			slog.Error("mandatory env variable not set", slog.String("key", v))
			return false
		}
	}
	return true

}
