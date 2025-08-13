/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
)

func Test_GetEdgeGatewayBandwidth(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: types.ParamsEdgeGateway{
				ID: "invalid-id",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 500,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Error 404",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC := newClient(t)

			// Set up mock response
			ep := endpoints.ListT0()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			result, err := eC.GetBandwidth(t.Context(), tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil when error is expected")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", err)
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotEmpty(t, result.ID, "Expected edge gateway ID to match")
				assert.NotEmpty(t, result.Name, "Expected edge gateway name to match")
			}
		})
	}
}
