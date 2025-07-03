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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRequest_Vmware(t *testing.T) {
	client, err := newMockClient()
	assert.Nil(t, err, "Error creating mock client")

	endpointSessionVmware, err := GetEndpoint("CreateSessionVmware", MethodPOST)
	assert.Nil(t, err, "Error getting endpoint for CreateSessionVmware")
	defer endpointSessionVmware.CleanMockResponse()

	tests := []struct {
		name           string
		expectedErr    bool
		expectedResp   any
		expectedStatus int
	}{
		{
			name:           "success",
			expectedErr:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad request",
			expectedErr:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal server error",
			expectedErr:    true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "not found",
			expectedErr: true,
			expectedResp: vmwareError{
				Message:       "Resource not found",
				StatusCode:    404,
				StatusMessage: "RESOURCE_NOT_FOUND",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "unauthorized",
			expectedErr:    true,
			expectedStatus: http.StatusUnauthorized,
			expectedResp: vmwareError{
				Message:       "Unauthorized access",
				StatusCode:    401,
				StatusMessage: "UNAUTHORIZED",
			},
		},
		{
			name:           "unknown error",
			expectedErr:    true,
			expectedStatus: http.StatusBadGateway,
			expectedResp: struct {
				Foo string `json:"foo"`
			}{
				Foo: "Unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer endpointSessionVmware.CleanMockResponse()

			endpointSessionVmware.SetMockResponse(tt.expectedResp, &tt.expectedStatus)

			req, err := client.NewRequest(t.Context(), ClientVmware)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error creating request")
			assert.NotNil(t, req, "Expected request to be created")

			endpointSessionVmware.CleanMockResponse()
		})
	}
}
