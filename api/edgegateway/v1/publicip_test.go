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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListEdgegatewayPublicIP(t *testing.T) {
	validEdgeGWName := generator.MustGenerate("{resource_name:edgegateway}")
	validEdgeGWID := "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c"
	validIP := generator.MustGenerate("{ipv4address}")

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
				ID: validEdgeGWID,
			},
			mockResponse: &itypes.ApiResponseNetworkServices{
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
									Name:      "internet",
									ServiceID: "test-publicip-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{
										IP:        validIP,
										Announced: true,
									},
								},
							},
						},
					},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
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
				ID: validEdgeGWID,
			},
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			epServices := endpoints.GetEdgeGatewayServices()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				statusCode := tt.mockResponseStatus
				ms.CleanResponse(epServices)
				ms.SetResponse(epServices, tt.mockResponse, &statusCode)
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockResponse, &statusCode)
			}

			resp, err := client.ListPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.PublicIPs, "Public IPs should not be empty")
			for _, ip := range resp.PublicIPs {
				assert.NotEmpty(t, ip.ID, "Public IP ID should not be empty")
				assert.NotEmpty(t, ip.IP, "Public IP Address should not be empty")
			}

			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
		})
	}
}

func TestGetEdgegatewayPublicIP(t *testing.T) {
	validEdgeGWName := generator.MustGenerate("{resource_name:edgegateway}")
	validEdgeGWID := "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c"
	validIP := generator.MustGenerate("{ipv4address}")

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
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponse: &itypes.ApiResponseNetworkServices{
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
									Name:      "internet",
									ServiceID: "test-publicip-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{
										IP:        validIP,
										Announced: true,
									},
								},
							},
						},
					},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
		},
		{
			name: "Valid request by name",
			params: types.ParamsGetEdgeGatewayPublicIP{
				IP:   validIP,
				Name: validEdgeGWName,
			},
			mockListResponse: &itypes.ApiResponseQueryEdgeGateway{
				Record: []itypes.ApiResponseQueryEdgeGatewayRecord{
					{
						ID:   validEdgeGWID,
						HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						Name: validEdgeGWName,
					},
				},
			},
			mockListResponseStatus: http.StatusOK,
			mockResponse: &itypes.ApiResponseNetworkServices{
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
									Name:      "internet",
									ServiceID: "test-publicip-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{
										IP:        validIP,
										Announced: true,
									},
								},
							},
						},
					},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
		},
		{
			name: "Failed request by name",
			params: types.ParamsGetEdgeGatewayPublicIP{
				IP:   validIP,
				Name: validEdgeGWName,
			},
			mockListResponseStatus: http.StatusNotFound,
			expectedErr:            true,
		},
		{
			name: "Invalid request",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID:   "invalid-id",
				Name: "invalid-name",
			},
			expectedErr: true,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true,
		},
		{
			name: "Simulate empty response",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponse:       &itypes.ApiResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			epServices := endpoints.GetEdgeGatewayServices()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				statusCode := tt.mockResponseStatus
				ms.CleanResponse(epServices)
				ms.SetResponse(epServices, tt.mockResponse, &statusCode)
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockResponse, &statusCode)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				listStatusCode := tt.mockListResponseStatus
				ms.CleanResponse(epQuery)
				ms.SetResponse(epQuery, tt.mockListResponse, &listStatusCode)
				ms.CleanResponse(endpoints.ListVdc())
				ms.SetResponse(endpoints.ListVdc(), tt.mockListResponse, &listStatusCode)
			}

			resp, err := client.GetPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.Equal(t, tt.params.IP, resp.IP, "Public IP Address should match")

			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
			ms.CleanResponse(endpoints.QueryEdgeGateway())
			ms.CleanResponse(endpoints.ListVdc())
		})
	}
}

func TestCreateEdgegatewayPublicIP(t *testing.T) {
	validEdgeGWName := generator.MustGenerate("{resource_name:edgegateway}")
	validEdgeGWID := "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c"
	validIP := "195.25.101.7"

	tests := []struct {
		name   string
		params types.ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockJobResponse       any
		mockJobResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		mockGetNetworkServicesResponse       any
		mockGetNetworkServicesResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid request",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: validIP,
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
									Name:      "internet",
									ServiceID: "test-publicip-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{
										IP:        validIP,
										Announced: true,
									},
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
			name: "Valid request by name",
			params: types.ParamsEdgeGateway{
				Name: validEdgeGWName,
			},
			mockListResponse: &itypes.ApiResponseQueryEdgeGateway{
				Record: []itypes.ApiResponseQueryEdgeGatewayRecord{
					{
						ID:   validEdgeGWID,
						HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						Name: validEdgeGWName,
					},
				},
			},
			mockListResponseStatus: http.StatusOK,
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: validIP,
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
									Name:      "internet",
									ServiceID: "test-publicip-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{
										IP:        validIP,
										Announced: true,
									},
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
			name: "Failed request by name",
			params: types.ParamsEdgeGateway{
				Name: validEdgeGWName,
			},
			mockListResponse:       nil,
			mockListResponseStatus: http.StatusNotFound,
			expectedErr:            true,
		},
		{
			name: "Job failed",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
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
				ID: validEdgeGWID,
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			epCreatePublicIp := endpoints.CreatePublicIp()
			if tt.mockResponseStatus != 0 {
				statusCode := tt.mockResponseStatus
				ms.SetResponse(epCreatePublicIp, tt.mockResponse, &statusCode)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				listStatusCode := tt.mockListResponseStatus
				ms.SetResponse(epQuery, tt.mockListResponse, &listStatusCode)
				ms.CleanResponse(endpoints.ListVdc())
				ms.SetResponse(endpoints.ListVdc(), tt.mockListResponse, &listStatusCode)
			}

			epGetJob := endpoints.GetJobCerberus()
			if tt.mockJobResponseStatus != 0 {
				jobStatusCode := tt.mockJobResponseStatus
				ms.SetResponse(epGetJob, tt.mockJobResponse, &jobStatusCode)
			}

			epGetNetworkServices := endpoints.GetEdgeGatewayServices()
			if tt.mockGetNetworkServicesResponse != nil || tt.mockGetNetworkServicesResponseStatus != 0 {
				statusCode := tt.mockGetNetworkServicesResponseStatus
				ms.CleanResponse(epGetNetworkServices)
				ms.SetResponse(epGetNetworkServices, tt.mockGetNetworkServicesResponse, &statusCode)
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockGetNetworkServicesResponse, &statusCode)
			}

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
			client, ms := newClient(t)
			epDisable := endpoints.DisableCloudavenueServices()
			if tt.mockResponseStatus != 0 {
				statusCode := tt.mockResponseStatus
				ms.SetResponse(epDisable, tt.mockResponse, &statusCode)
			}

			err := client.DeletePublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}
