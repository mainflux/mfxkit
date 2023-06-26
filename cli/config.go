// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Offset    uint   `toml:"offset"`
	Limit     uint   `toml:"limit"`
	Name      string `toml:"name"`
	RawOutput bool   `toml:"raw_output"`
}

// read - retrieve config from a file.
func read(file string) (Config, error) {
	data, err := os.ReadFile(file)
	c := Config{}
	if err != nil {
		return c, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := toml.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config TOML: %w", err)
	}

	return c, nil
}

func ParseConfig() {
	if ConfigPath == "" {
		// No config file
		return
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		errConfigNotFound := errors.Wrap(errors.New("config file was not found"), err)
		logError(errConfigNotFound)

		return
	}

	config, err := read(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	if config.RawOutput {
		RawOutput = config.RawOutput
	}
}
