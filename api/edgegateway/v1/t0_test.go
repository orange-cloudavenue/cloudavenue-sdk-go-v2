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
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func Test_ListT0(t *testing.T) {
	tests := []struct {
		name string

		mockResponse       any
		mockResponseStatus int

		expectedErr   bool
		errorContains string
	}{
		{
			name:        "List T0",
			expectedErr: false,
		},
		{
			name:               "Error 500",
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        true, // Error HTTP 500 returns an error after retries are exhausted.
		},
		{
			name:               "Error 404",
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true, // Error HTTP 404 should return an error.
		},
		{
			name: "Simulate unknown class of service",
			mockResponse: &itypes.APIResponseT0Names{
				generator.MustGenerate("{resource_name:t0}"),
			},
			expectedErr:        false,
			mockResponseStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)
			ep := endpoints.ListT0()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(ep)
				ms.SetResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			t0s, err := eC.ListT0(t.Context())
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error while listing T0s")
				assert.NotNil(t, t0s, "Expected non-nil T0s response")
			}
		})
	}
}

func Test_GetT0(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsGetT0

		mockResponse       any
		mockResponseStatus int

		expectedErr   bool
		errorContains string
	}{
		{
			name: "Valid T0",
			params: types.ParamsGetT0{
				T0Name: "prvrf01eocb0001234allsp01",
			},
			mockResponse: &itypes.APIResponseT0{
				Type:       "tier-0-vrf",
				Name:       "prvrf01eocb0001234allsp01",
				Properties: itypes.APIResponseT0Properties{ClassOfService: "SHARED_STANDARD"},
			},
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name: "Invalid TO name",
			params: types.ParamsGetT0{
				T0Name: "invalid_t0_name",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        true, // Error HTTP 500 returns an error after retries are exhausted.
		},
		{
			name: "Simulate empty response",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponse:       &itypes.APIResponseT0{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "Get by EdgeGateway Name",
			params: types.ParamsGetT0{
				EdgegatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Children: []itypes.APIResponseNetworkServicesChildren{{
						Type: "edge-gateway",
						Name: "test-edgegateway-name",
					}},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
		},
		{
			name: "Get by EdgeGateway ID",
			params: types.ParamsGetT0{
				EdgegatewayID: "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Children: []itypes.APIResponseNetworkServicesChildren{{
						Type: "edge-gateway",
						Name: generator.MustGenerate("{resource_name:edgegateway}"),
						Properties: struct {
							RateLimit int    `json:"rateLimit,omitempty"`
							EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
						}{
							EdgeUUID: "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						},
					}},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
		},
		{
			name: "EdgeGateway Name not found",
			params: types.ParamsGetT0{
				EdgegatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponse:       &itypes.APIResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "EdgeGateway ID not found",
			params: types.ParamsGetT0{
				EdgegatewayID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       &itypes.APIResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "Error 404",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true, // Error HTTP 404 should return an error.
		},
		{
			name:          "Missing edge gateway reference",
			params:        types.ParamsGetT0{},
			expectedErr:   true,
			errorContains: "edge gateway id or name is required",
		},
		{
			name:          "Invalid edge gateway ID",
			params:        types.ParamsGetT0{EdgegatewayID: "not-an-edge-gateway-urn"},
			expectedErr:   true,
			errorContains: "invalid edge gateway ID",
		},
		{
			name:          "Invalid edge gateway name",
			params:        types.ParamsGetT0{EdgegatewayName: "invalid edge gateway name"},
			expectedErr:   true,
			errorContains: "invalid edge gateway name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)
			epList := endpoints.ListT0()
			epGet := cav.MustGetEndpoint("GetT0")
			epServices := endpoints.GetEdgeGatewayServices()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				switch tt.mockResponse.(type) {
				case *itypes.APIResponseT0, *itypes.APIResponseT0s:
					ms.CleanResponse(epGet)
					ms.SetResponse(epGet, tt.mockResponse, &tt.mockResponseStatus)
				case *itypes.APIResponseNetworkServices:
					ms.CleanResponse(epServices)
					ms.SetResponse(epServices, tt.mockResponse, &tt.mockResponseStatus)
				default:
					ms.CleanResponse(epGet)
					ms.SetResponse(epGet, tt.mockResponse, &tt.mockResponseStatus)
				}
				ms.CleanResponse(epList)
			}

			if tt.params.EdgegatewayID != "" || tt.params.EdgegatewayName != "" {
				detail := &itypes.APIResponseT0{
					Type:       "tier-0-vrf",
					Name:       "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{ClassOfService: "SHARED_STANDARD"},
					Children: []itypes.APIResponseT0Children{{
						Type: "edge-gateway",
						Name: func() string {
							if tt.params.EdgegatewayName != "" {
								return tt.params.EdgegatewayName
							}
							return generator.MustGenerate("{resource_name:edgegateway}")
						}(),
						Properties: struct {
							RateLimit int    `json:"rateLimit,omitempty" fake:"5"`
							EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
						}{
							RateLimit: 5,
							EdgeUUID: func() string {
								if tt.params.EdgegatewayID != "" {
									return tt.params.EdgegatewayID
								}
								return generator.MustGenerate("{urn:edgegateway}")
							}(),
						},
					}},
				}
				status := http.StatusOK
				ms.CleanResponse(epGet)
				ms.SetResponse(epGet, detail, &status)
			}

			if tt.name == "Simulate empty response" {
				status := http.StatusOK
				ms.CleanResponse(epGet)
				ms.SetResponse(epGet, &itypes.APIResponseT0{}, &status)
			}

			t0, err := eC.GetT0(t.Context(), tt.params)

			if tt.expectedErr {
				require.Error(t, err)
				assert.Nil(t, t0, "Expected nil T0 response")
				if tt.errorContains != "" {
					assert.ErrorContains(t, err, tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, t0)
				if tt.params.T0Name != "" {
					assert.Equal(t, tt.params.T0Name, t0.Name, "Expected T0 name to match the requested name")
				}
				if tt.params.EdgegatewayID != "" || tt.params.EdgegatewayName != "" {
					assert.NotEmpty(t, t0.EdgeGateways, "Expected T0 to have edge gateways")
				}
			}
		})
	}
}

func Test_GetT0NotFoundPreservesResponseMetadata(t *testing.T) {
	eC, ms := newClient(t)
	ep := endpoints.GetEdgeGatewayServices()
	ms.CleanResponse(ep)
	status := http.StatusOK
	ms.SetResponse(ep, &itypes.APIResponseNetworkServices{}, &status)

	_, err := eC.GetT0(t.Context(), types.ParamsGetT0{
		EdgegatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
	})
	require.Error(t, err)

	var apiErr *pkgerrors.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, http.MethodGet, apiErr.Method)
	assert.NotEmpty(t, apiErr.Endpoint)
}
