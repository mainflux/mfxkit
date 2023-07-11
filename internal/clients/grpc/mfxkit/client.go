// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"github.com/mainflux/mainflux/pkg/errors"
	grpcClient "github.com/mainflux/mfxkit/internal/clients/grpc"
	"github.com/mainflux/mfxkit/internal/env"
	"github.com/mainflux/mfxkit/mfxkit"
	api "github.com/mainflux/mfxkit/mfxkit/api/grpc"
)

const envAuthGrpcPrefix = "MF_MFXKIT_GRPC_"

var errGrpcConfig = errors.New("failed to load grpc configuration")

// Setup loads Auth gRPC configuration from environment variable and creates new Auth gRPC API.
func Setup(envPrefix, jaegerURL, svcName string) (mfxkit.ServiceClient, grpcClient.ClientHandler, error) {
	config := grpcClient.Config{}
	if err := env.Parse(&config, env.Options{Prefix: envAuthGrpcPrefix, AltPrefix: envPrefix}); err != nil {
		return nil, nil, errors.Wrap(errGrpcConfig, err)
	}
	c, ch, err := grpcClient.Setup(config, svcName, jaegerURL)
	if err != nil {
		return nil, nil, err
	}

	return api.NewClient(c.ClientConn, config.Timeout), ch, nil
}
