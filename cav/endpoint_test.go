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

	"github.com/stretchr/testify/assert"
)

func Test_Endpoint_Register(t *testing.T) {
	tests := []struct {
		Endpoint
		expectedError bool
	}{
		{
			Endpoint: Endpoint{
				Category:         Category("fake"),
				Version:          VersionV1,
				Name:             "fake",
				Method:           MethodPOST,
				SubClient:        ClientVmware,
				PathTemplate:     "/1.2.3/sessions",
				PathParams:       []PathParam{},
				QueryParams:      []QueryParam{},
				DocumentationURL: "https://foo.bar",
				RequestFunc:      nil,
			},
			expectedError: false,
		},
		{
			Endpoint: Endpoint{
				Category: "", // Invalid category
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Endpoint.Register()
			if tt.expectedError {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}
