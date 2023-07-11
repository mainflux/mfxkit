// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package tracing

import (
	"context"

	"github.com/mainflux/mfxkit/mfxkit"
	"go.opentelemetry.io/otel/trace"
)

var _ mfxkit.Service = (*tracingMiddleware)(nil)

type tracingMiddleware struct {
	tracer trace.Tracer
	svc    mfxkit.Service
}

// New returns a new group service with tracing capabilities.
func New(svc mfxkit.Service, tracer trace.Tracer) mfxkit.Service {
	return &tracingMiddleware{tracer, svc}
}

// Ping traces the "Ping" operation of the wrapped policies.Service.
func (tm *tracingMiddleware) Ping(ctx context.Context, secret string) (string, error) {
	ctx, span := tm.tracer.Start(ctx, "svc_ping")
	defer span.End()

	return tm.svc.Ping(ctx, secret)
}
