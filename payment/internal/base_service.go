package internal

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/payment/internal/config"
	"github.com/wathuta/technical_test/payment/internal/handler"
	"github.com/wathuta/technical_test/payment/internal/platform/mpesa"
	"github.com/wathuta/technical_test/payment/internal/repository"
	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
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
	mpesaService := mpesa.NewMpesa(&mpesa.MpesaOpts{
		ConsumerKey:    os.Getenv(config.MpesaConsumerKeyEnvVar),
		ConsumerSecret: os.Getenv(config.MpesaConsumerSecreteEnvVar),
	})

	handler := handler.New(repo, mpesaService)

	paymentpb.RegisterPaymentServiceServer(grpcSrv, handler)

	go func() {
		fmt.Println("GRPC Server is running on:", listener.Addr())
		if err := grpcSrv.Serve(listener); err != nil {
			slog.Error("error", err)
			os.Exit(1)
		}
	}()
	go func() {
		serveHTTP(handler)
	}()

	return &Service{db: db, grpcSrv: grpcSrv}, nil
}

func serveHTTP(h *handler.Handler) {
	mux := gin.Default()
	mux.POST("/callback", h.CallbackHandler)
	mux.Run("localhost:5002")

	fmt.Println("REST server is running on localhost:5002")
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
