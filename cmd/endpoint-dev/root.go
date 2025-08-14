/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"github.com/spf13/cobra"
)

var (
	loggerLevel string
	mockFlag    bool
	configFile  string

	rootCmd = &cobra.Command{
		Use:   "endpoint-dev",
		Short: "CloudAvenue Endpoint development CLI",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&loggerLevel, "logger", "debug", "Set the logger level (e.g., debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVar(&mockFlag, "mock", false, "Use the mock client (default: false)")
	rootCmd.PersistentFlags().StringVar(&configFile, "file", "", "Path to config file (default: $HOME/.sdkv2/devtools/config.yaml)")
}
