package orderhandler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wathuta/technical_test/orders/internal/repository"
	pproto "github.com/wathuta/technical_test/protos_gen/orders"
)

const (
	defaultPageSize = 10
	maxPageSize     = 1000
	maxUpdateRetry  = 5
)

var (
	errInternal                   = status.Error(codes.Internal, "internal error")
	errNotFound                   = status.Error(codes.InvalidArgument, "resource not found")
	errResourceRequired           = status.Error(codes.InvalidArgument, "resource required")
	errResourceUpdateMaskRequired = status.Error(codes.InvalidArgument, "resource update mask required")
)

type ProfileService struct {
	pproto.UnimplementedOrderServiceServer
	persist repository.Repository
}

func New(persist repository.Repository) *ProfileService {
	return &ProfileService{
		persist: persist,
	}
}
