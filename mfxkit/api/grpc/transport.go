// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mfxkit/mfxkit"
	"go.opentelemetry.io/contrib/instrumentation/github.com/go-kit/kit/otelkit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ mfxkit.ServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	ping kitgrpc.Handler
	mfxkit.UnimplementedServiceServer
}

// NewServer returns new ServiceServer instance.
func NewServer(svc mfxkit.Service) mfxkit.ServiceServer {
	return &grpcServer{
		ping: kitgrpc.NewServer(
			otelkit.EndpointMiddleware(otelkit.WithOperation("ping"))(pingEndpoint(svc)),
			decodePingRequest,
			encodePingResponse,
		),
	}
}

func (s *grpcServer) Ping(ctx context.Context, req *mfxkit.PingReq) (*mfxkit.PingRes, error) {
	_, res, err := s.ping.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*mfxkit.PingRes), nil
}

func decodePingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*mfxkit.PingReq)

	return pingReq{Secret: req.GetSecret()}, nil
}

func encodePingResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(pingRes)

	return &mfxkit.PingRes{Message: res.message}, nil
}

func encodeError(err error) error {
	switch {
	case errors.Contains(err, nil):
		return nil
	case errors.Contains(err, mfxkit.ErrMalformedEntity):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Contains(err, mfxkit.ErrUnauthorizedAccess):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
