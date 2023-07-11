// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mfxkit

import (
	"context"
	"errors"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")
)

type service struct {
	secret string
}

var _ Service = (*service)(nil)

// New instantiates the mfxkit service implementation.
func New(secret string) Service {
	return &service{
		secret: secret,
	}
}

func (ks *service) Ping(ctx context.Context, secret string) (string, error) {
	if ks.secret != secret {
		return "", ErrUnauthorizedAccess
	}

	return "Hello World :)", nil
}
