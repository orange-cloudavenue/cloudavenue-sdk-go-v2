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

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// settings hold the values of all client options.
type settings struct {
	// Organization is the name of the organization to which the client belongs.
	Organization string
	// Console contains all properties related to the console of the client.
	Console consoles.ConsoleName
	// SubClients contains the sub-clients for the client.
	SubClients map[subClientName]subClientInterface
	// CachePassphrase is the passphrase used for caching.
	CachePassphrase string
	// CachePath is the path to the cache file.
	CachePath string
}

func newSettings(organization string) *settings {
	return &settings{
		Organization: organization,
		SubClients:   make(map[subClientName]subClientInterface),
	}
}

// ClientOption  is a function which applies options to a settings object.
type ClientOption func(*settings) error

// * Internal options

// withConsole sets the console for the client.
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

// * Exporter options

// WithCustomEndpoints sets custom endpoints for the sub-clients.
func WithCustomEndpoints(endpoints consoles.Services) ClientOption {
	return func(s *settings) error {
		logger := xlogger.WithGroup("client").WithGroup("options").WithGroup("WithCustomEndpoints")
		logger.Debug("Overriding endpoints in the console", "console", s.Console.GetSiteID())
		s.Console.OverrideEndpoint(endpoints)
		return nil
	}
}

// WithCloudAvenueCredential sets the credential for the client.
func WithCloudAvenueCredential(username, password string) ClientOption {
	return func(s *settings) error {
		logger := xlogger.WithGroup("client").WithGroup("options").WithGroup("WithCloudAvenueCredential")

		// Auth client is created before the loop to avoid creating it multiple times.
		// auth cloudavenue is shared between sub-clients vmware and cerberus.
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

// WithLogger sets the logger for the client.
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

// WithCache store the tokens in a cache
func WithCache(passphrase, path string) ClientOption {
	return func(s *settings) error {
		s.CachePassphrase = passphrase
		s.CachePath = path
		return nil
	}
}
