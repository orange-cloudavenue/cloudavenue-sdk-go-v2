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

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func Test_GetEdgeGatewayBandwidth(t *testing.T) {
	validEdgeGatewayName := generator.MustGenerate("{resource_name:edgegateway}")

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
				ID: "urn:vcloud:gateway:test-edge-gw-id",
			},
			mockResponse: func() *itypes.APIResponseT0s {
				child := itypes.APIResponseT0Children{
					Type: "edge-gateway",
					Name: "test-edge-gw",
				}
				child.Properties.RateLimit = 5
				child.Properties.EdgeUUID = "urn:vcloud:gateway:test-edge-gw-id"
				return &itypes.APIResponseT0s{
					{
						Type:     "tier-0-vrf",
						Name:     "test-t0",
						Children: []itypes.APIResponseT0Children{child},
					},
				}
			}(),
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: types.ParamsEdgeGateway{
				Name: validEdgeGatewayName,
			},
			mockResponse: func() *itypes.APIResponseT0s {
				child := itypes.APIResponseT0Children{
					Type: "edge-gateway",
					Name: validEdgeGatewayName,
				}
				child.Properties.RateLimit = 5
				child.Properties.EdgeUUID = "urn:vcloud:gateway:test-edge-gw-id"
				return &itypes.APIResponseT0s{
					{
						Type:     "tier-0-vrf",
						Name:     "test-t0",
						Children: []itypes.APIResponseT0Children{child},
					},
				}
			}(),
			mockResponseStatus: 200,
			expectedErr:        false,
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
			expectedErr:        true,
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
			eC, ms := newClient(t)

			// Set up mock response on both endpoint identities sharing the same path.
			ep := endpoints.ListT0()
			epSharedPath := endpoints.GetEdgeGatewayServices()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ms.CleanResponse(ep)
				ms.CleanResponse(epSharedPath)
				ms.SetResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
				ms.SetResponse(epSharedPath, tt.mockResponse, &tt.mockResponseStatus)
			}

			result, err := eC.GetBandwidth(t.Context(), tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil when error is expected")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", err)
				if assert.NotNil(t, result, "Result should not be nil") {
					assert.NotEmpty(t, result.ID, "Expected edge gateway ID to match")
					assert.NotEmpty(t, result.Name, "Expected edge gateway name to match")
				}
			}
		})
	}
}
