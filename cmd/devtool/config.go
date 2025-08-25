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
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

type Config struct {
	Organization string
	Username     string
	Password     string
}

func loadConfig() (*Config, error) {
	var (
		cfg Config
		err error
	)

	cfg.Organization, err = keyring.Get("sdkdevtool", "organization")
	if err != nil {
		return nil, err
	}
	cfg.Username, err = keyring.Get("sdkdevtool", "username")
	if err != nil {
		return nil, err
	}
	cfg.Password, err = keyring.Get("sdkdevtool", "password")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

var configCmdV2 = &cobra.Command{
	Use:   "config",
	Short: "Store your configuration (organization, username, and password) securely in your system keystore.",
	Run: func(cmd *cobra.Command, args []string) {
		var service = "sdkdevtool"

		if cmdConfigClear {
			err := keyring.DeleteAll(service)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if cmdConfigOrganization == "" || cmdConfigUsername == "" || cmdConfigPassword == "" {
			log.Fatal("Organization, username, and password must be set")
		}

		// set organization in keyring
		err := keyring.Set(service, "organization", cmdConfigOrganization)
		if err != nil {
			log.Fatal(err)
		}

		// set username and password in keyring
		err = keyring.Set(service, "username", cmdConfigUsername)
		if err != nil {
			log.Fatal(err)
		}

		err = keyring.Set(service, "password", cmdConfigPassword)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Your credentials have been saved securely in your system keystore.")
	},
}

var (
	cmdConfigOrganization string
	cmdConfigUsername     string
	cmdConfigPassword     string
	cmdConfigClear        bool
)

func init() {
	configCmdV2.Flags().StringVar(&cmdConfigOrganization, "organization", "", "Organization name")
	configCmdV2.Flags().StringVar(&cmdConfigUsername, "username", "", "Username")
	configCmdV2.Flags().StringVar(&cmdConfigPassword, "password", "", "Password")
	configCmdV2.Flags().BoolVar(&cmdConfigClear, "clear", false, "Clear configuration")

	rootCmd.AddCommand(configCmdV2)
}
