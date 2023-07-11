// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-zoo/bone"
	chclient "github.com/mainflux/callhome/pkg/client"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mfxkit/internal"
	jaegerClient "github.com/mainflux/mfxkit/internal/clients/jaeger"
	"github.com/mainflux/mfxkit/internal/env"
	"github.com/mainflux/mfxkit/internal/server"
	grpcserver "github.com/mainflux/mfxkit/internal/server/grpc"
	httpserver "github.com/mainflux/mfxkit/internal/server/http"
	"github.com/mainflux/mfxkit/mfxkit"
	"github.com/mainflux/mfxkit/mfxkit/api"
	grpcapi "github.com/mainflux/mfxkit/mfxkit/api/grpc"
	httpapi "github.com/mainflux/mfxkit/mfxkit/api/http"
	"github.com/mainflux/mfxkit/mfxkit/tracing"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	svcName        = "ping"
	envPrefix      = "MF_MFXKIT_"
	envPrefixHTTP  = "MF_MFXKIT_HTTP_"
	envPrefixGRPC  = "MF_MFXKIT_GRPC_"
	defSvcHTTPPort = "9099"
	defSvcGRPCPort = "9199"
)

type config struct {
	LogLevel      string `env:"MF_MFXKIT_LOG_LEVEL"   envDefault:"error"`
	Secret        string `env:"MF_MFXKIT_SECRET"      envDefault:"secret"`
	JaegerURL     string `env:"MF_JAEGER_URL"         envDefault:"http://jaeger:14268/api/traces"`
	SendTelemetry bool   `env:"MF_SEND_TELEMETRY"     envDefault:"true"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err.Error())
	}

	logger, err := mflog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to init logger: %s", err.Error())
	}

	tp, err := jaegerClient.NewProvider(svcName, cfg.JaegerURL)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to init Jaeger: %s", err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)

	svc := newService(cfg.Secret, logger, tracer)

	httpServerConfig := server.Config{Port: defSvcHTTPPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHTTP, AltPrefix: envPrefix}); err != nil {
		logger.Fatal(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err.Error()))
	}
	mux := bone.New()
	hs := httpserver.New(ctx, cancel, svcName, httpServerConfig, httpapi.MakeHandler(svc, mux, logger), logger)

	registerAuthServiceServer := func(srv *grpc.Server) {
		reflection.Register(srv)
		mfxkit.RegisterServiceServer(srv, grpcapi.NewServer(svc))
	}
	grpcServerConfig := server.Config{Port: defSvcGRPCPort}
	if err := env.Parse(&grpcServerConfig, env.Options{Prefix: envPrefixGRPC, AltPrefix: envPrefix}); err != nil {
		log.Fatalf("failed to load %s gRPC server configuration : %s", svcName, err.Error())
	}
	gs := grpcserver.New(ctx, cancel, svcName, grpcServerConfig, registerAuthServiceServer, logger)

	if cfg.SendTelemetry {
		chc := chclient.New(svcName, mainflux.Version, logger, cancel)
		go chc.CallHome(ctx)
	}

	g.Go(func() error {
		return hs.Start()
	})
	g.Go(func() error {
		return gs.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, hs, gs)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("%s service terminated: %s", svcName, err))
	}
}

func newService(secret string, logger mflog.Logger, tracer trace.Tracer) mfxkit.Service {
	svc := mfxkit.New(secret)
	svc = tracing.New(svc, tracer)
	svc = api.LoggingMiddleware(svc, logger)
	counter, latency := internal.MakeMetrics(svcName, "api")
	svc = api.MetricsMiddleware(svc, counter, latency)

	return svc
}
