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
	"context"
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetEdgeGatewayServices(t *testing.T) {
	tests := []struct {
		name               string
		params             *types.ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid Edge Gateway services",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway services with name",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockQueryResponse:       nil,
			mockQueryResponseStatus: http.StatusNotFound,
			expectedErr:             true,
		},
		{
			name:        "Simulate empty params",
			params:      &types.ParamsEdgeGateway{},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        true,
		},
		{
			name: "Simulate empty response",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       &itypes.ApiResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)
			ep := endpoints.GetEdgeGatewayServices()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				ms.CleanResponse(ep)
				// Set the mock response
				ms.SetResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			// Set up mock query response
			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				// Clean all default mock responses
				ms.CleanResponse(epQuery)
				// Set the mock query response
				ms.SetResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			// Call the GetNetworkServices method
			result, err := eC.GetServices(t.Context(), *tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil: %v", result)
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
				assert.NotNil(t, result, "Result should not be nil: %v", tt.params)
			}
		})
	}
}

func TestGetEdgeGatewayServices_ContextDeadlineExceeded(t *testing.T) {
	mC, _, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err = eC.GetServices(ctx, types.ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestEnableCloudavenueServices(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		expectedErr bool
	}{
		{
			name: "Enable network services with valid ID",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable network services with valid Name",
			params: types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable network services with empty params",
			params: types.ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        true,
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockQueryResponse:       nil,
			mockQueryResponseStatus: http.StatusNotFound,
			expectedErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)
			ep := endpoints.EnableCloudavenueServices()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ms.CleanResponse(ep)
				ms.SetResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			// Set up mock query response
			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				t.Log("Setting up mock query response for:", tt.name)
				ms.CleanResponse(epQuery)
				ms.SetResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			err := eC.EnableCloudavenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
			}
		})
	}
}

func TestEnableCloudavenueServices_ContextDeadlineExceeded(t *testing.T) {
	mC, _, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	err = eC.EnableCloudavenueServices(ctx, types.ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestDisableCloudavenueServices(t *testing.T) {
	validEdgeGWName := generator.MustGenerate("{resource_name:edgegateway}")
	validEdgeGWID := "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c"

	tests := []struct {
		name   string
		params types.ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		mockGetNetworkServicesResponse       any
		mockGetNetworkServicesResponseStatus int

		expectedErr bool
	}{
		{
			name: "Disable network services with valid ID",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Name: "test-t0",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			expectedErr:                          false,
		},
		{
			name: "Disable network services with valid Name",
			params: types.ParamsEdgeGateway{
				Name: validEdgeGWName,
			},
			mockQueryResponse: &itypes.ApiResponseQueryEdgeGateway{
				Record: []itypes.ApiResponseQueryEdgeGatewayRecord{
					{
						ID:   validEdgeGWID,
						HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						Name: validEdgeGWName,
					},
				},
			},
			mockQueryResponseStatus: http.StatusOK,
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			expectedErr:                          false,
		},
		{
			name: "Disable network services with empty params",
			params: types.ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Failed to get network services",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockGetNetworkServicesResponse:       nil,
			mockGetNetworkServicesResponseStatus: http.StatusNotFound,
			expectedErr:                          true,
		},

		{
			name: "Error 500",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			mockResponseStatus:                   http.StatusInternalServerError,
			expectedErr:                          true,
		},
		{
			name: "Error 401",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			mockResponseStatus:                   http.StatusUnauthorized,
			expectedErr:                          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)
			ep := endpoints.DisableCloudavenueServices()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ms.CleanResponse(ep)
				ms.SetResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				epQuery := endpoints.QueryEdgeGateway()
				ms.CleanResponse(epQuery)
				ms.SetResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
				ms.CleanResponse(endpoints.ListVdc())
				ms.SetResponse(endpoints.ListVdc(), tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			epGetNetworkServices := endpoints.GetEdgeGatewayServices()
			// Set up mock response for GetNetworkServices
			if tt.mockGetNetworkServicesResponse != nil || tt.mockGetNetworkServicesResponseStatus != 0 {
				t.Log("Setting up mock GetNetworkServices response for:", tt.name)
				ms.CleanResponse(epGetNetworkServices)
				ms.SetResponse(epGetNetworkServices, tt.mockGetNetworkServicesResponse, &tt.mockGetNetworkServicesResponseStatus)
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockGetNetworkServicesResponse, &tt.mockGetNetworkServicesResponseStatus)
			}

			err := eC.DisableCloudavenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
			}

			ms.CleanResponse(ep)
			ms.CleanResponse(endpoints.QueryEdgeGateway())
			ms.CleanResponse(endpoints.ListVdc())
			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
		})
	}
}

func TestGetCloudavenueServices(t *testing.T) {
	validEdgeGWName := generator.MustGenerate("{resource_name:edgegateway}")
	validEdgeGWID := "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c"

	tests := []struct {
		name   string
		params types.ParamsEdgeGateway

		mockQueryResponse       any
		mockQueryResponseStatus int

		mockGetNetworkServicesResponse       any
		mockGetNetworkServicesResponseStatus int

		expectedErr bool
	}{
		{
			name: "Get Cloud Avenue services with valid ID",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			expectedErr:                          false,
		},
		{
			name: "Get Cloud Avenue services with valid Name",
			params: types.ParamsEdgeGateway{
				Name: validEdgeGWName,
			},
			mockQueryResponse: &itypes.ApiResponseQueryEdgeGateway{
				Record: []itypes.ApiResponseQueryEdgeGatewayRecord{
					{
						ID:   validEdgeGWID,
						HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						Name: validEdgeGWName,
					},
				},
			},
			mockQueryResponseStatus: http.StatusOK,
			mockGetNetworkServicesResponse: &itypes.ApiResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.ApiResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.ApiResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      "cav-services",
									ServiceID: "test-service-id",
								},
							},
						},
					},
				},
			},
			mockGetNetworkServicesResponseStatus: http.StatusOK,
			expectedErr:                          false,
		},
		{
			name: "Get Cloud Avenue services with empty params",
			params: types.ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponseStatus: http.StatusInternalServerError,
			expectedErr:                          true,
		},
		{
			name: "Error 401",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockGetNetworkServicesResponseStatus: http.StatusUnauthorized,
			expectedErr:                          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				epQuery := endpoints.QueryEdgeGateway()
				ms.CleanResponse(epQuery)
				ms.SetResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
				ms.CleanResponse(endpoints.ListVdc())
				ms.SetResponse(endpoints.ListVdc(), tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			epGetNetworkServices := endpoints.GetEdgeGatewayServices()
			if tt.mockGetNetworkServicesResponse != nil || tt.mockGetNetworkServicesResponseStatus != 0 {
				t.Log("Setting up mock GetNetworkServices response for:", tt.name)
				ms.CleanResponse(epGetNetworkServices)
				ms.SetResponse(epGetNetworkServices, tt.mockGetNetworkServicesResponse, &tt.mockGetNetworkServicesResponseStatus)
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockGetNetworkServicesResponse, &tt.mockGetNetworkServicesResponseStatus)
			}

			result, err := eC.GetCloudavenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil")
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, result, "Result should not be nil")
			}

			ms.CleanResponse(endpoints.QueryEdgeGateway())
			ms.CleanResponse(endpoints.ListVdc())
			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
		})
	}
}
