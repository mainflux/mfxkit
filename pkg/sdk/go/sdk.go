// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

const (
	// CTJSON represents JSON content type.
	CTJSON ContentType = "application/json"

	// CTBinary represents binary content type.
	CTBinary ContentType = "application/octet-stream"

	// EnabledStatus represents enable status for a client.
	EnabledStatus = "enabled"

	// DisabledStatus represents disabled status for a client.
	DisabledStatus = "disabled"

	BearerPrefix = "Bearer "

	ThingPrefix = "Thing "
)

// ContentType represents all possible content types.
type ContentType string

var _ SDK = (*mfSDK)(nil)

// SDK contains Mainflux API.
type SDK interface {
	// Ping sends a ping request to Mainflux.
	//
	// Mainflux responds with a greeting message.
	//
	// example:
	//
	//  greeting, err := sdk.Ping("my-secret")
	//  if err != nil {
	//      fmt.Println(err)
	//  }
	//  fmt.Println(greeting)
	Ping(secret string) (string, error)

	// Health sends a health request to Mainflux.
	//
	// Mainflux responds with service status.
	//
	// example:
	//
	//  health, err := sdk.Health()
	//  if err != nil {
	//      fmt.Println(err)
	//  }
	//  fmt.Println(health)
	Health() (HealthInfo, error)
}

type mfSDK struct {
	mfxkitURL string

	msgContentType ContentType
	client         *http.Client
}

// Config contains sdk configuration parameters.
//
// example:
//
//	conf := sdk.Config{
//	    MFxkitURL: "http://localhost:9099",
//	    MsgContentType: sdk.CTJSON,
//	    TLSVerification: false,
//	}
type Config struct {
	MFxkitURL string

	MsgContentType  ContentType
	TLSVerification bool
}

// NewSDK returns new mainflux SDK instance.
//
// example:
//
//	conf := sdk.Config{
//	    MFxkitURL: "http://localhost:9099",
//	    MsgContentType: sdk.CTJSON,
//	    TLSVerification: false,
//	}
//
//	sdk := sdk.NewSDK(conf)
func NewSDK(conf Config) SDK {
	return &mfSDK{
		mfxkitURL:      conf.MFxkitURL,
		msgContentType: conf.MsgContentType,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !conf.TLSVerification,
				},
			},
		},
	}
}

// processRequest creates and send a new HTTP request, and checks for errors in the HTTP response.
// It then returns the response headers, the response body, and the associated error(s) (if any).
func (sdk mfSDK) processRequest(method, url, token, contentType string, data []byte, expectedRespCodes ...int) (http.Header, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return make(http.Header), []byte{}, err
	}

	if token != "" {
		if !strings.Contains(token, ThingPrefix) {
			token = BearerPrefix + token
		}
		req.Header.Set("Authorization", token)
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := sdk.client.Do(req)
	if err != nil {
		return make(http.Header), []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return make(http.Header), []byte{}, err
	}

	for _, code := range expectedRespCodes {
		if resp.StatusCode == code {
			return resp.Header, body, nil
		}
	}

	return nil, nil, nil
}
