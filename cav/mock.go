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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	subclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/subClient"
)

// TODO move to internal/subClient/mock.go
// WithMock sets the credential for the client.
func WithMock() ClientOption {
	return func(s *settings) error {
		if s.SubClients[mock] == nil {
			s.SubClients[mock] = subclient.Clients[mock]
		}

		s.SubClients[mock].SetConsole(s.Console)
		s.SubClients[mock].SetCredential(auth.NewMockAuth(map[string]string{
			"X-Mock": "mock",
		}))

		return nil
	}
}

// WithMockJob sets the mock job client for testing purposes.
func WithMockJob() ClientOption {
	return func(s *settings) error {
		if s.SubClients[mockJob] == nil {
			s.SubClients[mockJob] = subclient.Clients[mockJob]
		}

		s.SubClients[mockJob].SetConsole(s.Console)
		s.SubClients[mockJob].SetCredential(auth.NewMockAuth(map[string]string{
			"X-Mock-Job": "mock-job",
		}))

		return nil
	}
}
