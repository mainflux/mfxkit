// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/mfxkit/mfxkit"
)

func pingEndpoint(svc mfxkit.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pingReq)
		if err := req.validate(); err != nil {
			return pingRes{}, err
		}

		res, err := svc.Ping(ctx, req.Secret)
		if err != nil {
			return pingRes{}, err
		}

		return pingRes{message: res}, nil
	}
}
