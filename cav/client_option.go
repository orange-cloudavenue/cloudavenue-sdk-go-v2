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

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	subclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/subClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// settings hold the values of all client options.
type settings struct {
	// Organization is the name of the organization to which the client belongs.
	Organization string
	// Console contains all properties related to the console of the client.
	Console consoles.Console
	// SubClients contains the sub-clients for the client.
	SubClients map[subclient.Name]subclient.Client

	httpClient *resty.Client
}

func newSettings(organization string) *settings {
	return &settings{
		Organization: organization,
		SubClients:   make(map[subclient.Name]subclient.Client),
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

// WithCloudAvenueCredential sets the credential for the client.
func WithCloudAvenueCredential(username, password string) ClientOption {
	return func(s *settings) error {
		for _, client := range []subclient.Name{subclient.Cerberus, subclient.Vmware} {
			if s.SubClients[client] == nil {
				s.SubClients[client] = subclient.Clients[client]
			}

			s.SubClients[client].SetConsole(s.Console)

			cred, err := auth.NewCloudavenueCredential(s.httpClient, s.Console, s.Organization, username, password)
			if err != nil {
				return err
			}
			s.SubClients[client].SetCredential(cred)
		}

		return nil
	}
}
