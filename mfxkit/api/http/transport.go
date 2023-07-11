// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux"
	mflog "github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mfxkit/mfxkit"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/go-kit/kit/otelkit"
)

const contentType = "application/json"

// errUnsupportedContentType indicates unacceptable or lack of Content-Type.
var errUnsupportedContentType = errors.New("unsupported content type")

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc mfxkit.Service, mux *bone.Mux, logger mflog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(loggingErrorEncoder(logger, encodeError)),
	}

	mux.Post("/ping", kithttp.NewServer(
		otelkit.EndpointMiddleware(otelkit.WithOperation("ping"))(pingEndpoint(svc)),
		decodePing,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/health", mainflux.Health("mfxkit"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodePing(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errUnsupportedContentType
	}

	req := pingReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(mainflux.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case mfxkit.ErrMalformedEntity,
		io.ErrUnexpectedEOF,
		io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	case mfxkit.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// loggingErrorEncoder is a go-kit error encoder logging decorator.
func loggingErrorEncoder(logger mflog.Logger, enc kithttp.ErrorEncoder) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		switch {
		case errors.Contains(err, mfxkit.ErrMalformedEntity),
			errors.Contains(err, io.ErrUnexpectedEOF),
			errors.Contains(err, io.EOF),
			errors.Contains(err, mfxkit.ErrUnauthorizedAccess),
			errors.Contains(err, errUnsupportedContentType):
			logger.Error(err.Error())
		}

		enc(ctx, err, w)
	}
}
