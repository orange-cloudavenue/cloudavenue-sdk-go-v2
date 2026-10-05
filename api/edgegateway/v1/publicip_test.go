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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	sdkerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	// realPublicIPServiceID is a realistic CloudAvenue service identifier, as
	// returned by the network hierarchy for an allocated public IP. Real
	// identifiers are opaque and cannot be derived from the IP address.
	realPublicIPServiceID = "tn01i01ocb1010314spt102-cav-services"

	// testPublicIP is the public IP address used across the delete tests.
	testPublicIP = "195.25.101.7"
)

// recordedRequest captures a request the SDK actually put on the wire.
type recordedRequest struct {
	Method string
	Path   string
}

// newInternetService builds an "internet" network service entry as returned by
// the network hierarchy endpoint.
func newInternetService(serviceID, ip string) itypes.APIResponseNetworkServicesSubChildren {
	return itypes.APIResponseNetworkServicesSubChildren{
		Type:      "service",
		Name:      networkTypeInternet,
		ServiceID: serviceID,
		Properties: struct {
			ClassOfService     string   `json:"classOfService,omitempty"`
			MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
			IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
			Announced          bool     `json:"announced,omitempty" fake:"true"`
			Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
		}{
			IP:        ip,
			Announced: true,
		},
	}
}

// newNamedService builds a non-internet "service" entry, used to pin that only
// internet services are matched during the lookup.
func newNamedService(name, serviceID, ip string) itypes.APIResponseNetworkServicesSubChildren {
	service := newInternetService(serviceID, ip)
	service.Name = name
	return service
}

// newEdgeGatewayChild builds the edge gateway level of the network hierarchy.
func newEdgeGatewayChild(name, edgeUUID string, services ...itypes.APIResponseNetworkServicesSubChildren) itypes.APIResponseNetworkServicesChildren {
	return itypes.APIResponseNetworkServicesChildren{
		Type: "edge-gateway",
		Name: name,
		Properties: struct {
			RateLimit int    `json:"rateLimit,omitempty"`
			EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
		}{
			EdgeUUID: edgeUUID,
		},
		Children: services,
	}
}

// newNetworkHierarchy wraps services into an org-wide network hierarchy.
func newNetworkHierarchy(services ...itypes.APIResponseNetworkServicesSubChildren) *itypes.APIResponseNetworkServices {
	return &itypes.APIResponseNetworkServices{
		{
			Type:     "tier-0-vrf",
			Children: []itypes.APIResponseNetworkServicesChildren{newEdgeGatewayChild("edge-gateway-test", "ed0a243a-374b-4306-ab25-9c3787cbdb4c", services...)},
		},
	}
}

// setNetworkHierarchy mocks the org-wide network hierarchy read.
func setNetworkHierarchy(t *testing.T, ms *mock.Server, response any, statusCode int) {
	t.Helper()

	if response == nil && statusCode == 0 {
		return
	}

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ep := endpoints.GetEdgeGatewayServices()
	ms.CleanResponse(ep)
	ms.SetResponse(ep, response, &statusCode)
}

// recordDeleteRequests installs a handler on the delete endpoint that records
// the method and path of every request the SDK sends to it, then answers with
// statusCode.
//
// The mock router is built from the endpoint path template itself, so this is
// the only way to assert the path the SDK really emitted for a given endpoint.
func recordDeleteRequests(t *testing.T, ms *mock.Server, statusCode int) *[]recordedRequest {
	t.Helper()

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	recorded := &[]recordedRequest{}

	ms.SetResponseFunc(endpoints.DisableCloudavenueServices(), func(w http.ResponseWriter, r *http.Request) {
		*recorded = append(*recorded, recordedRequest{Method: r.Method, Path: r.URL.Path})

		if statusCode >= 300 {
			http.Error(w, http.StatusText(statusCode), statusCode)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jobId":"87ab1934-0146-4fb0-80bc-815fea03214d","message":"Job created successfully"}`)) //nolint:errcheck
	})

	return recorded
}

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
		assertResp  func(*testing.T, *types.ModelEdgeGatewayPublicIPs)
	}{
		{
			name: "Valid request",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
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
			assertResp: func(t *testing.T, resp *types.ModelEdgeGatewayPublicIPs) {
				assert.Len(t, resp.PublicIPs, 1)
			},
		},
		{
			name: "Stable ordering by IP then ID",
			params: types.ParamsEdgeGateway{
				ID: validEdgeGWID,
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
									ServiceID: "b-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{IP: "192.0.2.2"},
								},
								{
									Type:      "service",
									Name:      networkTypeInternet,
									ServiceID: "c-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{IP: "192.0.2.1"},
								},
								{
									Type:      "service",
									Name:      networkTypeInternet,
									ServiceID: "a-id",
									Properties: struct {
										ClassOfService     string   `json:"classOfService,omitempty"`
										MaxVirtualServices int      `json:"maxVirtualServices,omitempty"`
										IP                 string   `json:"ip,omitempty" fake:"{ipv4address}"`
										Announced          bool     `json:"announced,omitempty" fake:"true"`
										Ranges             []string `json:"ranges,omitempty" fake:"{ipv4address}/{intrange:24,32}"`
									}{IP: "192.0.2.1"},
								},
							},
						},
					},
				},
			},
			mockResponseStatus: http.StatusOK,
			expectedErr:        false,
			assertResp: func(t *testing.T, resp *types.ModelEdgeGatewayPublicIPs) {
				assert.Len(t, resp.PublicIPs, 3)
				assert.Equal(t, "192.0.2.1", resp.PublicIPs[0].IP)
				assert.Equal(t, "a-id", resp.PublicIPs[0].ID)
				assert.Equal(t, "192.0.2.1", resp.PublicIPs[1].IP)
				assert.Equal(t, "c-id", resp.PublicIPs[1].ID)
				assert.Equal(t, "192.0.2.2", resp.PublicIPs[2].IP)
				assert.Equal(t, "b-id", resp.PublicIPs[2].ID)
			},
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
				ms.CleanResponse(cav.MustGetEndpoint("GetT0"))
				ms.SetResponse(cav.MustGetEndpoint("GetT0"), tt.mockResponse, &statusCode)
			}

			resp, err := client.ListPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.PublicIPs, "Public IPs should not be empty")
			if tt.assertResp != nil {
				tt.assertResp(t, resp)
			}
			for _, ip := range resp.PublicIPs {
				assert.NotEmpty(t, ip.ID, "Public IP ID should not be empty")
				assert.NotEmpty(t, ip.IP, "Public IP Address should not be empty")
			}

			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
			ms.CleanResponse(cav.MustGetEndpoint("GetT0"))
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
		assertErr   func(*testing.T, error)
	}{
		{
			name: "Valid request",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
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
			mockListResponse: &itypes.APIResponseQueryEdgeGateway{
				Record: []itypes.APIResponseQueryEdgeGatewayRecord{
					{
						ID:   validEdgeGWID,
						HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
						Name: validEdgeGWName,
					},
				},
			},
			mockListResponseStatus: http.StatusOK,
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
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
			assertErr: func(t *testing.T, err error) {
				var apiErr *sdkerrors.APIError
				assert.ErrorAs(t, err, &apiErr)
				assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
			},
		},
		{
			name: "Simulate empty response",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponse:       &itypes.APIResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "No matching edge gateway in hierarchy",
			params: types.ParamsGetEdgeGatewayPublicIP{
				ID: validEdgeGWID,
				IP: validIP,
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: "other-edge",
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "11111111-1111-1111-1111-111111111111",
							},
						},
					},
				},
			},
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
				ms.CleanResponse(cav.MustGetEndpoint("GetT0"))
				ms.SetResponse(cav.MustGetEndpoint("GetT0"), tt.mockResponse, &statusCode)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				listStatusCode := tt.mockListResponseStatus
				ms.CleanResponse(epQuery)
				ms.SetResponse(epQuery, tt.mockListResponse, &listStatusCode)
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), tt.mockListResponse, &listStatusCode)
			}

			resp, err := client.GetPublicIP(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.Equal(t, tt.params.IP, resp.IP, "Public IP Address should match")

			ms.CleanResponse(endpoints.GetEdgeGatewayServices())
			ms.CleanResponse(endpoints.ListT0())
			ms.CleanResponse(cav.MustGetEndpoint("GetT0"))
			ms.CleanResponse(endpoints.QueryEdgeGateway())
			ms.CleanResponse(endpoints.ListVDC())
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
			mockGetNetworkServicesResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
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
			mockListResponse: &itypes.APIResponseQueryEdgeGateway{
				Record: []itypes.APIResponseQueryEdgeGatewayRecord{
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
			mockGetNetworkServicesResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						{
							Type: "edge-gateway",
							Name: validEdgeGWName,
							Properties: struct {
								RateLimit int    `json:"rateLimit,omitempty"`
								EdgeUUID  string `json:"edgeUuid,omitempty" fake:"{urn:edgegateway}"`
							}{
								EdgeUUID: "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							},
							Children: []itypes.APIResponseNetworkServicesSubChildren{
								{
									Type:      "service",
									Name:      networkTypeInternet,
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
			epCreatePublicIP := endpoints.CreatePublicIP()
			if tt.mockResponseStatus != 0 {
				statusCode := tt.mockResponseStatus
				ms.SetResponse(epCreatePublicIP, tt.mockResponse, &statusCode)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				listStatusCode := tt.mockListResponseStatus
				ms.SetResponse(epQuery, tt.mockListResponse, &listStatusCode)
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), tt.mockListResponse, &listStatusCode)
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
	tests := []struct {
		name   string
		params types.ParamsDeleteEdgeGatewayPublicIP

		mockResponse       any
		mockResponseStatus int

		// deleteStatus is the status code answered by the delete endpoint.
		// The delete request is always recorded, whatever the outcome.
		deleteStatus int

		expectedErr         string
		expectedDeletePath  string
		expectedDeleteCount int
	}{
		{
			name: "Valid request resolves the real service ID",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse:        newNetworkHierarchy(newInternetService(realPublicIPServiceID, testPublicIP)),
			mockResponseStatus:  http.StatusOK,
			deleteStatus:        http.StatusOK,
			expectedDeletePath:  "/infrapicustomerproxy/v2.0/services/" + realPublicIPServiceID,
			expectedDeleteCount: 1,
		},
		{
			name: "Lookup spans every edge gateway of the organization",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse: &itypes.APIResponseNetworkServices{
				{
					Type: "tier-0-vrf",
					Children: []itypes.APIResponseNetworkServicesChildren{
						newEdgeGatewayChild("edge-a", "ed0a243a-374b-4306-ab25-9c3787cbdb4c",
							newInternetService("other-service-id", "192.0.2.1")),
						newEdgeGatewayChild("edge-b", "11111111-1111-1111-1111-111111111111",
							newInternetService(realPublicIPServiceID, testPublicIP)),
					},
				},
			},
			mockResponseStatus:  http.StatusOK,
			deleteStatus:        http.StatusOK,
			expectedDeletePath:  "/infrapicustomerproxy/v2.0/services/" + realPublicIPServiceID,
			expectedDeleteCount: 1,
		},
		{
			name: "Only internet services are matched",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse: newNetworkHierarchy(
				newNamedService(cavServicesNetworkType, "cav-services-id", testPublicIP),
			),
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "public IP " + testPublicIP + " not found in network services",
			expectedDeleteCount: 0,
		},
		{
			name: "IP absent from the hierarchy issues no delete",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse:        newNetworkHierarchy(newInternetService(realPublicIPServiceID, "192.0.2.1")),
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "public IP " + testPublicIP + " not found in network services",
			expectedDeleteCount: 0,
		},
		{
			name: "Empty hierarchy issues no delete",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse:        &itypes.APIResponseNetworkServices{},
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "public IP " + testPublicIP + " not found in network services",
			expectedDeleteCount: 0,
		},
		{
			name: "Empty serviceID issues no delete",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse:        newNetworkHierarchy(newInternetService("", testPublicIP)),
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "serviceID is empty, cannot delete public IP " + testPublicIP,
			expectedDeleteCount: 0,
		},
		{
			name: "Network lookup failure issues no delete",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponseStatus:  http.StatusNotFound,
			expectedErr:         "error retrieving network services to resolve public IP " + testPublicIP,
			expectedDeleteCount: 0,
		},
		{
			name: "Invalid request",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: "invalid-ip",
			},
			mockResponse:        newNetworkHierarchy(newInternetService(realPublicIPServiceID, testPublicIP)),
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "invalid IP address:",
			expectedDeleteCount: 0,
		},
		{
			name: "Missing IP",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: "",
			},
			mockResponse:        newNetworkHierarchy(newInternetService(realPublicIPServiceID, testPublicIP)),
			mockResponseStatus:  http.StatusOK,
			expectedErr:         "ip is required",
			expectedDeleteCount: 0,
		},
		{
			name: "Error 404 Not Found on delete",
			params: types.ParamsDeleteEdgeGatewayPublicIP{
				IP: testPublicIP,
			},
			mockResponse:        newNetworkHierarchy(newInternetService(realPublicIPServiceID, testPublicIP)),
			mockResponseStatus:  http.StatusOK,
			deleteStatus:        http.StatusNotFound,
			expectedErr:         "error deleting public IP " + testPublicIP,
			expectedDeletePath:  "/infrapicustomerproxy/v2.0/services/" + realPublicIPServiceID,
			expectedDeleteCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			setNetworkHierarchy(t, ms, tt.mockResponse, tt.mockResponseStatus)
			recorded := recordDeleteRequests(t, ms, tt.deleteStatus)

			err := client.DeletePublicIP(t.Context(), tt.params)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
			}

			require.Len(t, *recorded, tt.expectedDeleteCount,
				"unexpected number of DELETE requests issued to DisableCloudavenueServices")

			if tt.expectedDeletePath != "" {
				assert.Equal(t, http.MethodDelete, (*recorded)[0].Method)
				assert.Equal(t, tt.expectedDeletePath, (*recorded)[0].Path)
			}
		})
	}
}

// TestDeletePublicIPNeverSynthesizesServiceID pins the exact path the SDK puts
// on the wire.
//
// The mock router is built from the endpoint path template itself, so a green
// mock suite proves nothing about the resolved identifier on its own. This test
// records the outgoing request and asserts the real service ID coming from the
// network hierarchy, not the bogus "ip-<dashed ip>" value the SDK used to
// fabricate from the address.
func TestDeletePublicIPNeverSynthesizesServiceID(t *testing.T) {
	client, ms := newClient(t)

	setNetworkHierarchy(t, ms, newNetworkHierarchy(newInternetService(realPublicIPServiceID, testPublicIP)), http.StatusOK)
	recorded := recordDeleteRequests(t, ms, http.StatusOK)

	require.NoError(t, client.DeletePublicIP(t.Context(), types.ParamsDeleteEdgeGatewayPublicIP{
		IP: testPublicIP,
	}))

	require.Len(t, *recorded, 1, "exactly one DELETE must be issued")
	assert.Equal(t, http.MethodDelete, (*recorded)[0].Method)
	assert.Equal(t,
		"/infrapicustomerproxy/v2.0/services/tn01i01ocb1010314spt102-cav-services",
		(*recorded)[0].Path,
		"DELETE must target the real service ID resolved from the network hierarchy",
	)

	// Explicit guard against a regression to the synthesized form.
	assert.NotEqual(t, "/infrapicustomerproxy/v2.0/services/ip-195-25-101-7", (*recorded)[0].Path)
	assert.NotContains(t, (*recorded)[0].Path, "ip-195-25-101-7")
}
