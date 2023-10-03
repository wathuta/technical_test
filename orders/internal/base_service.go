package internal

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	handler "github.com/wathuta/technical_test/orders/internal/handler"

	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	ordersPb "github.com/wathuta/technical_test/protos_gen/orders"
	prductsPb "github.com/wathuta/technical_test/protos_gen/products"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	"github.com/wathuta/technical_test/orders/internal/repository"
)

type Service struct {
	grpcSrv                 *grpc.Server
	GracefulShutdownTimeout time.Duration

	db *sqlx.DB
}
type Options struct {
	ListenAddress           string
	GracefulShutdownTimeout time.Duration
}

func NewService(ctx context.Context, db *sqlx.DB, opts Options) (*Service, error) {
	listener, err := net.Listen("tcp", opts.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Set up gRPC server
	grpcSrv := grpc.NewServer(grpc.EmptyServerOption{})
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(db)

	handler := handler.New(repo)

	ordersPb.RegisterOrderServiceServer(grpcSrv, handler)
	customersPb.RegisterCustomerServiceServer(grpcSrv, handler)
	prductsPb.RegisterProductServiceServer(grpcSrv, handler)

	go func() {
		slog.Info("starting the server", "listening address:", listener.Addr().String())
		if err := grpcSrv.Serve(listener); err != nil {
			slog.Error("error", err)
			os.Exit(1)
		}
	}()

	return &Service{db: db, grpcSrv: grpcSrv}, nil
}

func (s *Service) Shutdown() bool {
	c := make(chan struct{})

	go func() {
		defer close(c)

		// Block until all pending RPCs are finished
		s.grpcSrv.GracefulStop()
	}()

	select {
	case <-time.After(s.GracefulShutdownTimeout):
		// Timeout
		s.grpcSrv.Stop()
		<-c
		return false

	case <-c:
		// Shutdown completed within the timeout
		return true
	}
}
func (s *Service) Close() {
	s.db.Close()
}
