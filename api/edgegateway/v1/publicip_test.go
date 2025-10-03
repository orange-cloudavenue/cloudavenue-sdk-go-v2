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

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
)

func TestListEdgegatewayPublicIP(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		params types.ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid request",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid request",
			params: types.ParamsEdgeGateway{
				ID:   "invalid-id",
				Name: "invalid-name",
			},
			expectedErr: true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.GetEdgeGatewayServices().CleanMockResponse()
				endpoints.GetEdgeGatewayServices().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.PublicIps, "Public IPs should not be empty")
			for _, ip := range resp.PublicIps {
				assert.NotEmpty(t, ip.ID, "Public IP ID should not be empty")
				assert.NotEmpty(t, ip.IP, "Public IP Address should not be empty")
			}
		})
	}
}

func TestGetEdgegatewayPublicIP(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		params types.ParamsGetEdgeGatewayPublicIP

		mockResponse       any
		mockResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid request",
			params: types.ParamsGetEdgeGatewayPublicIP{
				EdgeGatewayID: generator.MustGenerate("{urn:edgegateway}"),
				IP:            generator.MustGenerate("{ipv4address}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid request by name",
			params: types.ParamsGetEdgeGatewayPublicIP{
				IP:              generator.MustGenerate("{ipv4address}"),
				EdgeGatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Failed request by name",
			params: types.ParamsGetEdgeGatewayPublicIP{
				IP:              generator.MustGenerate("{ipv4address}"),
				EdgeGatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockListResponseStatus: 404,
			expectedErr:            true,
		},
		{
			name: "Invalid request",
			params: types.ParamsGetEdgeGatewayPublicIP{
				EdgeGatewayID:   "invalid-id",
				EdgeGatewayName: "invalid-name",
			},
			expectedErr: true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsGetEdgeGatewayPublicIP{
				EdgeGatewayID: generator.MustGenerate("{urn:edgegateway}"),
				IP:            generator.MustGenerate("{ipv4address}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Simulate empty response",
			params: types.ParamsGetEdgeGatewayPublicIP{
				EdgeGatewayID: generator.MustGenerate("{urn:edgegateway}"),
				IP:            generator.MustGenerate("{ipv4address}"),
			},
			mockResponse:       &itypes.ApiResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.GetEdgeGatewayServices().CleanMockResponse()
				endpoints.GetEdgeGatewayServices().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListResponseStatus != 0 {
				endpoints.QueryEdgeGateway().CleanMockResponse()
				endpoints.QueryEdgeGateway().SetMockResponse(tt.mockListResponse, &tt.mockListResponseStatus)
			}

			client := newClient(t)

			resp, err := client.GetPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.Equal(t, tt.params.IP, resp.IP, "Public IP Address should match")
		})
	}
}

func TestCreateEdgegatewayPublicIP(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		params types.ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockJobResponse       any
		mockJobResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid request",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						//  {
						//     "details": "195.25.101.7",
						//     "name": "reserve_ip for Org cav01ev01ocb0006205 for public ip",
						//     "status": "DONE"
						//  },
						{
							Details: generator.MustGenerate("{ipv4address}"),
							Name:    "reserve_ip for Org cav01ev01ocb0001234 for public ip",
							Status:  "DONE",
						},
					},
					Name:        "Create PublicIP Job",
					Status:      "DONE",
					Description: "PublicIP created successfully",
				},
			},
			mockJobResponseStatus: 200,
			expectedErr:           false,
		},
		{
			name: "Valid request by name",
			params: types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						//  {
						//     "details": "195.25.101.7",
						//     "name": "reserve_ip for Org cav01ev01ocb0006205 for public ip",
						//     "status": "DONE"
						//  },
						{
							Details: generator.MustGenerate("{ipv4address}"),
							Name:    "reserve_ip for Org cav01ev01ocb0001234 for public ip",
							Status:  "DONE",
						},
					},
					Name:        "Create PublicIP Job",
					Status:      "DONE",
					Description: "PublicIP created successfully",
				},
			},
			mockJobResponseStatus: 200,
			expectedErr:           false,
		},
		{
			name: "Failed request by name",
			params: types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockListResponseStatus: 404,
			expectedErr:            true,
		},
		{
			name: "Job failed",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockJobResponseStatus: 400,
			expectedErr:           true,
		},
		{
			name: "Invalid request",
			params: types.ParamsEdgeGateway{
				ID: "invalid-id",
			},
			expectedErr: true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.CreatePublicIp().CleanMockResponse()
				endpoints.CreatePublicIp().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListResponseStatus != 0 {
				endpoints.QueryEdgeGateway().CleanMockResponse()
				endpoints.QueryEdgeGateway().SetMockResponse(tt.mockListResponse, &tt.mockListResponseStatus)
			}

			if tt.mockJobResponseStatus != 0 {
				endpoints.GetJobCerberus().CleanMockResponse()
				endpoints.GetJobCerberus().SetMockResponse(tt.mockJobResponse, &tt.mockJobResponseStatus)
			}

			client := newClient(t)

			resp, err := client.CreatePublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
		})
	}
}

func TestDeleteEdgegatewayPublicIP(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		params types.ParamsDeleteEdgeGatewayPublicIP

		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid request",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid request",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: "invalid-ip",
			},
			expectedErr: true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.DisableCloudavenueServices().CleanMockResponse()
				endpoints.DisableCloudavenueServices().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			err := client.DeletePublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}
