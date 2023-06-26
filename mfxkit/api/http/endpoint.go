// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/mfxkit/mfxkit"
)

func pingEndpoint(svc mfxkit.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pingReq)

		if err := req.validate(); err != nil {
			return nil, err
		}

		greeting, err := svc.Ping(ctx, req.Secret)
		if err != nil {
			return nil, err
		}

		res := pingRes{
			Greeting: greeting,
		}
		
		return res, nil
	}
}
