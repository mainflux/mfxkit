// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package env

// NewConfig gets configuration from environment variable.
func NewConfig[T any](opts ...Options) (T, error) {
	var cfg T
	if err := Parse(&cfg, opts...); err != nil {
		return cfg, err
	}
	
	return cfg, nil
}
