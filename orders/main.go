package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/wathuta/technical_test/orders/internal"
	"github.com/wathuta/technical_test/orders/internal/config"
	database "github.com/wathuta/technical_test/orders/internal/platform/postgres"
	"github.com/wathuta/technical_test/orders/internal/repository"
	"golang.org/x/exp/slog"
)

func main() {
	var err error
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	slog.SetDefault(l)
	if !config.HasAllEnvVariables() {
		envFileName := ".env.dev"
		slog.Info("loading env file", "fileName", envFileName)
		err = godotenv.Load(envFileName)
		if err != nil {
			slog.Error("unable to load env vars", "error", err)
			os.Exit(1)
		}
	}

	if os.Getenv(config.RunMigrationsEnvVar) == "true" {
		err = database.RunMigrations()
		if err != nil {
			slog.Error("unable to db migrations", "error", err)
			os.Exit(1)
		}
	}

	r, err := repository.NewPostgresRepository()
	if err != nil {
		slog.Error("unable to connect to db", "error", err)
		os.Exit(1)
	}
	service, err := internal.NewService(context.Background(), r, internal.Options{ListenAddress: os.Getenv(config.ListenAddressEnvVar)})
	if err != nil {
		slog.Error("failed to start service", "error", err)
		os.Exit(1)
	}
	shutdownOnSignal(service)
}
func waitForShutdownSignal() string {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until signaled
	sig := <-c

	return sig.String()
}

func shutdownOnSignal(svc *internal.Service) {
	signalName := waitForShutdownSignal()
	slog.Info("Received signal, starting shutdown", "signal", signalName)

	if svc.Shutdown() {
		slog.Info("Shutdown complete")
	} else {
		slog.Info("Shutdown timed out")
	}
}
