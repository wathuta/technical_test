package internal

import (
	"context"
	"net"
	"os"
	"time"

	handler "github.com/wathuta/technical_test/orders/internal/handler"
	pproto "github.com/wathuta/technical_test/protos_gen/orders"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	"github.com/wathuta/technical_test/orders/internal/repository"
)

type Service struct {
	grpcSrv                 *grpc.Server
	GracefulShutdownTimeout time.Duration

	repo repository.Repository
}
type Options struct {
	ListenAddress           string
	GracefulShutdownTimeout time.Duration
}

func NewService(ctx context.Context, repo repository.Repository, opts Options) (*Service, error) {
	listener, err := net.Listen("tcp", opts.ListenAddress)
	if err != nil {
		return nil, err
	}

	// Set up gRPC server
	grpcSrv := grpc.NewServer()
	if err != nil {
		return nil, err
	}
	profileSvc := handler.New(repo)

	pproto.RegisterOrderServiceServer(grpcSrv, profileSvc)

	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			slog.Error("error", err)
			os.Exit(1)
		}
	}()

	return &Service{repo: repo, grpcSrv: grpcSrv}, nil
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
