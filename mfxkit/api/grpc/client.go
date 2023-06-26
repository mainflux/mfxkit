// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/mainflux/mfxkit/mfxkit"
	"go.opentelemetry.io/contrib/instrumentation/github.com/go-kit/kit/otelkit"
	"google.golang.org/grpc"
)

const svcName = "mainflux.mfxkit.Service"

var _ mfxkit.ServiceClient = (*grpcClient)(nil)

type grpcClient struct {
	ping    endpoint.Endpoint
	timeout time.Duration
}

// NewClient returns new gRPC client instance.
func NewClient(conn *grpc.ClientConn, timeout time.Duration) mfxkit.ServiceClient {
	return &grpcClient{
		ping: otelkit.EndpointMiddleware(otelkit.WithOperation("ping"))(kitgrpc.NewClient(
			conn,
			svcName,
			"Ping",
			encodePingRequest,
			decodePingResponse,
			mfxkit.PingRes{},
		).Endpoint()),

		timeout: timeout,
	}
}

func (client *grpcClient) Ping(ctx context.Context, req *mfxkit.PingReq, _ ...grpc.CallOption) (r *mfxkit.PingRes, err error) {
	ctx, close := context.WithTimeout(ctx, client.timeout)
	defer close()
	preq := pingReq{Secret: req.GetSecret()}
	res, err := client.ping(ctx, preq)
	if err != nil {
		return &mfxkit.PingRes{}, err
	}

	pres := res.(pingRes)

	return &mfxkit.PingRes{Message: pres.message}, err
}

func encodePingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(pingReq)

	return &mfxkit.PingReq{
		Secret: req.Secret,
	}, nil
}

func decodePingResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(*mfxkit.PingRes)

	return pingRes{message: res.GetMessage()}, nil
}
