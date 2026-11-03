package svc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) Ping(ctx context.Context, request *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
