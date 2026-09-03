/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"context"
	"errors"
	"log/slog"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/http-client"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// settings stores client option values during construction.
type settings struct {
	Organization    string
	Console         consoles.ConsoleName
	SubClients      map[subClientName]subClientInterface
	CachePassphrase string
	CachePath       string
}

// ClientOption applies configuration to settings.
type ClientOption func(*settings) error

// WithCustomEndpoints overrides service endpoints on detected console.
func WithCustomEndpoints(endpoints consoles.Services) ClientOption {
	return func(s *settings) error {
		logger := xlogger.WithGroup("client").WithGroup("options").WithGroup("WithCustomEndpoints")
		logger.Debug("Overriding endpoints in the console", "console", s.Console.GetSiteID())
		s.Console.OverrideEndpoint(endpoints)
		return nil
	}
}

// WithCloudAvenueCredential configures shared CloudAvenue credentials.
func WithCloudAvenueCredential(username, password string) ClientOption {
	return func(s *settings) error {
		logger := xlogger.WithGroup("client").WithGroup("options").WithGroup("WithCloudAvenueCredential")

		cred, err := newCloudavenueCredential(s.Console, s.Organization, username, password)
		if err != nil {
			logger.Error("Failed to create Cloudavenue credential", "error", err)
			return err
		}

		for _, client := range []subClientName{ClientCerberus, ClientVmware} {
			if _, ok := s.SubClients[client]; !ok {
				s.SubClients[client] = subClients[client]
			}

			s.SubClients[client].setConsole(s.Console)
			s.SubClients[client].setCredential(cred)
		}

		return nil
	}
}

// WithLogger sets package and client logger.
func WithLogger(customLogger *slog.Logger) ClientOption {
	return func(_ *settings) error {
		xlog.SetGlobalLogger(customLogger)
		xlogger = customLogger
		if xlogger.Enabled(context.Background(), slog.LevelDebug) {
			httpclient.DebugMode = true
		}
		return nil
	}
}

// WithCache enables session cache persistence.
func WithCache(passphrase, path string) ClientOption {
	return func(s *settings) error {
		s.CachePassphrase = passphrase
		s.CachePath = path
		return nil
	}
}

func newSettings(organization string) *settings {
	return &settings{
		Organization: organization,
		SubClients:   make(map[subClientName]subClientInterface),
	}
}

func withConsole() ClientOption {
	return func(s *settings) error {
		c, ok := consoles.FindByOrganizationName(s.Organization)
		if !ok {
			return errors.New("console not found")
		}
		s.Console = c
		return nil
	}
}
