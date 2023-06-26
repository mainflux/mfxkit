// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package main contains cli main function to run the cli.
package main

import (
	"log"

	"github.com/mainflux/mfxkit/cli"
	sdk "github.com/mainflux/mfxkit/pkg/sdk/go"
	"github.com/spf13/cobra"
)

const defURL string = "http://localhost:9099"

func main() {
	msgContentType := string(sdk.CTJSON)
	sdkConf := sdk.Config{
		MFxkitURL:       defURL,
		MsgContentType:  sdk.ContentType(msgContentType),
		TLSVerification: false,
	}

	// Root
	var rootCmd = &cobra.Command{
		Use: "mfxkit-cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.ParseConfig()

			sdkConf.MsgContentType = sdk.ContentType(msgContentType)
			s := sdk.NewSDK(sdkConf)
			cli.SetSDK(s)
		},
	}

	// API commands
	healthCmd := cli.NewHealthCmd()
	pingCmd := cli.NewPingCmd()

	// Root Commands
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(pingCmd)

	// Root Flags

	rootCmd.PersistentFlags().StringVarP(
		&sdkConf.MFxkitURL,
		"mfxkit-url",
		"m",
		sdkConf.MFxkitURL,
		"MFxkit service URL",
	)

	rootCmd.PersistentFlags().StringVarP(
		&msgContentType,
		"content-type",
		"y",
		msgContentType,
		"Message content type",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&sdkConf.TLSVerification,
		"insecure",
		"i",
		sdkConf.TLSVerification,
		"Do not check for TLS cert",
	)

	rootCmd.PersistentFlags().StringVarP(
		&cli.ConfigPath,
		"config",
		"c",
		cli.ConfigPath,
		"Config path",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&cli.RawOutput,
		"raw",
		"r",
		cli.RawOutput,
		"Enables raw output mode for easier parsing of output",
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
