// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"time"

	mflog "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mfxkit/mfxkit"
)

var _ mfxkit.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger mflog.Logger
	svc    mfxkit.Service
}

// LoggingMiddleware adds logging facilities to the mfxkit service.
func LoggingMiddleware(svc mfxkit.Service, logger mflog.Logger) mfxkit.Service {
	return &loggingMiddleware{logger, svc}
}

// Ping logs the "Ping" request. It logs the secret and the time it took to complete the request.
// If the request fails, it logs the error.
func (lm *loggingMiddleware) Ping(ctx context.Context, secret string) (response string, err error) {
	defer func(begin time.Time) {
		message := fmt.Sprintf("Method ping for secret %s took %s to complete", secret, time.Since(begin))
		if err != nil {
			lm.logger.Warn(fmt.Sprintf("%s with error: %s.", message, err))

			return
		}
		lm.logger.Info(fmt.Sprintf("%s without errors.", message))
	}(time.Now())

	return lm.svc.Ping(ctx, secret)
}
