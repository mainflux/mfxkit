// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HealthInfo contains version endpoint response.
type HealthInfo struct {
	// Status contains service status.
	Status string `json:"status"`

	// Version contains current service version.
	Version string `json:"version"`

	// Commit represents the git hash commit.
	Commit string `json:"commit"`

	// Description contains service description.
	Description string `json:"description"`

	// BuildTime contains service build time.
	BuildTime string `json:"build_time"`
}

func (sdk mfSDK) Health() (HealthInfo, error) {
	url := fmt.Sprintf("%s/health", sdk.mfxkitURL)

	resp, err := sdk.client.Get(url)
	if err != nil {
		return HealthInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthInfo{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var h HealthInfo
	if err := json.NewDecoder(resp.Body).Decode(&h); err != nil {
		return HealthInfo{}, err
	}

	return h, nil
}
