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
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Endpoint_Register(t *testing.T) {
	endpointName := "FakeTestEndpointRegister"
	endpoints.mu.Lock()
	delete(endpoints.Map, endpointName)
	endpoints.mu.Unlock()
	t.Cleanup(func() {
		endpoints.mu.Lock()
		delete(endpoints.Map, endpointName)
		endpoints.mu.Unlock()
	})

	tests := []struct {
		Endpoint
		expectedError bool
	}{
		{
			Endpoint: Endpoint{
				Name:             endpointName,
				Method:           MethodPOST,
				Backend:          BackendVMware,
				PathTemplate:     "/1.2.3/sessions",
				PathParams:       []PathParam{},
				QueryParams:      []QueryParam{},
				DocumentationURL: "https://foo.bar",
				Description:      "This is a fake endpoint",
			},
			expectedError: false,
		},
		{
			Endpoint: Endpoint{
				Name: "invalid space name", // Invalid name
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.expectedError {
				panicked := false
				defer func() {
					if r := recover(); r != nil {
						panicked = true
					}
					require.True(t, panicked, "expected Register to panic")
				}()
			}
			tt.Endpoint.Register()
			if tt.expectedError {
				return
			}

			registered, err := GetEndpoint(tt.Name)
			require.NoError(t, err)
			require.NotNil(t, registered)
		})
	}
}
