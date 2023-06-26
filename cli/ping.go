// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package cli

import "github.com/spf13/cobra"

// NewPingCmd returns ping command.
func NewPingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ping <secret>",
		Short: "Ping mfxkit",
		Long: "Ping mfxkit\n" +
			"For example:\n" +
			"\tmfxkit-cli ping <secret>\n",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				logUsage(cmd.Use)

				return
			}

			greeting, err := sdk.Ping(args[0])
			if err != nil {
				logError(err)

				return
			}

			logJSON(greeting)

		},
	}
}
