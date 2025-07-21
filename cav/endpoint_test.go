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
)

func Test_Endpoint_Register(t *testing.T) {
	tests := []struct {
		Endpoint
		expectedError bool
	}{
		{
			Endpoint: Endpoint{
				api:              API("fake"),
				version:          VersionV1,
				Name:             "fake",
				Method:           MethodPOST,
				SubClient:        ClientVmware,
				PathTemplate:     "/1.2.3/sessions",
				PathParams:       []PathParam{},
				QueryParams:      []QueryParam{},
				DocumentationURL: "https://foo.bar",
				RequestFunc:      nil,
				Description:      "This is a fake endpoint",
			},
			expectedError: false,
		},
		{
			Endpoint: Endpoint{
				api: "", // Invalid api
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.expectedError {
				defer func() {
					if r := recover(); r != nil {
						// We successfully recovered from panic
						t.Log("Test passed, panic was caught!")
					}
				}()
			}
			tt.Endpoint.Register()
		})
	}
}
