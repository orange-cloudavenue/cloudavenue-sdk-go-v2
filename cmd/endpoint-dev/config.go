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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Organization string `yaml:"organization"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

var (
	org      string
	username string
	password string
)

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "cloudavenue-config.yaml"
	}
	return filepath.Join(home, ".sdkv2", "devtools", "config.yaml")
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup configuration file with organization, username, and password",
	Run: func(cmd *cobra.Command, args []string) {
		path := configFile
		if path == "" {
			path = defaultConfigPath()
		}

		cfg := &Config{}

		// Try to load existing config
		if existing, err := loadConfig(path); err == nil {
			cfg = existing
		}

		// Use flags if provided, otherwise keep existing values
		if org != "" {
			cfg.Organization = org
		}
		if username != "" {
			cfg.Username = username
		}
		if password != "" {
			cfg.Password = password
		}

		// Check required fields
		missing := []string{}
		if cfg.Organization == "" {
			missing = append(missing, "organization")
		}
		if cfg.Username == "" {
			missing = append(missing, "username")
		}
		if cfg.Password == "" {
			missing = append(missing, "password")
		}
		if len(missing) > 0 {
			fmt.Printf("Missing required fields: %s\n", strings.Join(missing, ", "))
			os.Exit(1)
		}

		if err := saveConfig(path, cfg); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration saved to %s\n", path)
	},
}

func init() {
	configCmd.Flags().StringVar(&org, "organization", "", "Organization name")
	configCmd.Flags().StringVar(&username, "username", "", "Username")
	configCmd.Flags().StringVar(&password, "password", "", "Password")

	rootCmd.AddCommand(configCmd)
}
