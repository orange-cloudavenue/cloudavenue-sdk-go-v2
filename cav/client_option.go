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
	"errors"
	"log"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// settings hold the values of all client options.
type settings struct {
	// Organization is the name of the organization to which the client belongs.
	Organization string
	// Console contains all properties related to the console of the client.
	Console consoles.Console
	// SubClients contains the sub-clients for the client.
	SubClients map[SubClientName]SubClient

	httpClient *resty.Client
}

func newSettings(organization string) *settings {
	return &settings{
		Organization: organization,
		SubClients:   make(map[SubClientName]SubClient),
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
		s.Console.OverrideEndpoint(endpoints)
		return nil
	}
}

// WithCloudAvenueCredential sets the credential for the client.
func WithCloudAvenueCredential(username, password string) ClientOption {
	return func(s *settings) error {
		for _, client := range []SubClientName{ClientCerberus, ClientVmware} {
			if _, ok := s.SubClients[client]; !ok {
				s.SubClients[client] = subClients[client]
			}

			log.Default().Printf("Setting console %s for client %s", s.Console, client)
			s.SubClients[client].SetConsole(s.Console)

			cred, err := NewCloudavenueCredential(s.Console, s.Organization, username, password)
			if err != nil {
				return err
			}
			s.SubClients[client].SetCredential(cred)
		}

		return nil
	}
}
