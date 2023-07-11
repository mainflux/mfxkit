// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mfxkit

import "context"

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Ping compares a given string with secret.
	Ping(ctx context.Context, secret string) (res string, err error)
}
