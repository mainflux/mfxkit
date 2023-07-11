// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	mflog "github.com/mainflux/mainflux/logger"
	grpcClient "github.com/mainflux/mfxkit/internal/clients/grpc"
)

func Close(log mflog.Logger, clientHandler grpcClient.ClientHandler) {
	if err := clientHandler.Close(); err != nil {
		log.Warn(err.Error())
	}
}
