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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetEdgeGateway(t *testing.T) {
	tests := []struct {
		name                    string
		params                  *types.ParamsEdgeGateway
		mockQueryResponse       any
		mockQueryResponseStatus int
		mockResponse            any
		mockResponseStatus      int
		expectedErr             bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr:             false,
			mockQueryResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{
				Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
					ID:   generator.MustGenerate("{urn:edgegateway}"),
					HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				}},
			},
			mockResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponseStatus: 200,
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockQueryResponseStatus: 404,
			expectedErr:             true,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: &types.ParamsEdgeGateway{
				ID: "urn:vcloud:vm:invalid-id",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr:        true,
			mockResponse:       nil,
			mockResponseStatus: 500,
		},
		{
			name:        "Error validation params",
			params:      &types.ParamsEdgeGateway{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.SetResponse(endpoints.GetEdgeGateway(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				ms.SetResponse(endpoints.QueryEdgeGateway(), tt.mockQueryResponse, &tt.mockQueryResponseStatus)
				ms.SetResponse(endpoints.ListVDC(), tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			if edgeResp, ok := tt.mockResponse.(*itypes.APIResponseEdgegateway); ok && tt.mockResponseStatus == 200 {
				ms.SetResponse(endpoints.GetEdgeGateway(), edgeResp, &tt.mockResponseStatus)
			}

			// Call the GetEdgeGateway method
			result, err := eC.GetEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error: %v", tt.params)
				assert.Nil(t, result, "Result should be nil: %v", tt.params)
			} else {
				assert.Nil(t, err, "Expected no error: %v", tt.params)
				assert.NotNil(t, result, "Result should not be nil: %v", tt.params)
			}
		})
	}
}

func TestGetEdgeGateway_ContextDeadlineExceeded(t *testing.T) {
	mC, _, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err = eC.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestRetrieveEdgeGatewayIDByName(t *testing.T) {
	mC, ms, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	tests := []struct {
		name        string
		edgeName    string
		queryResp   *itypes.APIResponseQueryEdgeGateway
		queryStatus int
		expectedID  string
		expectedErr bool
	}{
		{
			name:     "Valid Edge Gateway Name",
			edgeName: generator.MustGenerate("{resource_name:edgegateway}"),
			queryResp: &itypes.APIResponseQueryEdgeGateway{
				Record: []itypes.APIResponseQueryEdgeGatewayRecord{
					{ID: "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c", HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c"},
				},
			},
			queryStatus: 200,
			expectedID:  "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.SetResponse(endpoints.QueryEdgeGateway(), tt.queryResp, &tt.queryStatus)
			ms.SetResponse(endpoints.ListVDC(), tt.queryResp, &tt.queryStatus)

			id, err := eC.retrieveEdgeGatewayIDByName(t.Context(), tt.edgeName)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Empty(t, id, "Expected empty ID but got %s", id)
			} else {
				assert.Nil(t, err, "Expected no error but got %v", err)
				assert.Equal(t, tt.expectedID, id, "Expected ID %s but got %s", tt.expectedID, id)
			}
		})
	}
}

func TestDeleteEdgeGateway(t *testing.T) {
	tests := []struct {
		name                    string
		params                  *types.ParamsEdgeGateway
		mockResponse            any
		mockResponseStatus      int
		mockQueryResponseStatus int
		expectedErr             bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse: nil,
			// mockResponseStatus: 202,
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponse: nil,
			// mockResponseStatus: 204,
			expectedErr: false,
		},
		{
			name: "Invalid Edge Gateway Name",
			params: &types.ParamsEdgeGateway{
				Name: "invalidEdgeGateway",
			},
			mockResponse:       nil,
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: &types.ParamsEdgeGateway{
				ID: "urn:vcloud:gateway:invalid-id",
			},
			mockResponse:       nil,
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Error 500",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       nil,
			mockResponseStatus: 500,
			expectedErr:        true,
		},
		{
			name: "Error 401",
			params: &types.ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       nil,
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "error 404 edge gateway name and id not found",
			params: &types.ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponse:            nil,
			mockResponseStatus:      404,
			mockQueryResponseStatus: 404,
			expectedErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.SetResponse(endpoints.DeleteEdgeGateway(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponseStatus != 0 {
				ms.SetResponse(endpoints.QueryEdgeGateway(), nil, &tt.mockQueryResponseStatus)
				ms.SetResponse(endpoints.ListVDC(), nil, &tt.mockQueryResponseStatus)
			}

			err := eC.DeleteEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error for params: %v", tt.params)
			} else {
				assert.Nil(t, err, "Expected no error for params: %v", tt.params)
			}
		})
	}
}

func TestDeleteEdgeGateway_ContextDeadlineExceeded(t *testing.T) {
	eC, _ := newClient(t)

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	err := eC.DeleteEdgeGateway(ctx, types.ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestCreateEdgeGateway(t *testing.T) {
	tests := []struct {
		name   string
		params *types.ParamsCreateEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockJobResponse       any
		mockJobResponseStatus int

		mockGetEdgeGatewayResponse       any
		mockGetEdgeGatewayResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		mockListVDCResponse       any
		mockListVDCResponseStatus int

		mockListVDCGroupResponse       any
		mockListVDCGroupResponseStatus int

		mockListT0Response       any
		mockListT0ResponseStatus int

		mockUpdateEdgeGatewayBandwidthResponse       any
		mockUpdateEdgeGatewayBandwidthResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid Edge Gateway Creation with VDC",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdc",
			},
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				Name: "owner-vdc",
				HREF: "https://api.example.com/api/vdc/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   "urn:vcloud:vdc:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockListVDCResponseStatus: 200,
			mockGetEdgeGatewayResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockGetEdgeGatewayResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   generator.MustGenerate("{urn:edgegateway}"),
			}}},
			mockQueryResponseStatus: 200,
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: generator.MustGenerate("{resource_name:edgegateway}"),
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockListT0ResponseStatus: 200,
			mockJobResponseStatus:    200,
			expectedErr:              false,
		},
		{
			name: "Valid Edge Gateway Creation with VDCGroup",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdcgroup",
			},
			mockListVDCResponseStatus: 200,
			mockListVDCResponse:       &itypes.APIResponseListVDC{},
			mockGetEdgeGatewayResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: "tn-edge-created-vdcgroup",
			},
			mockGetEdgeGatewayResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				Name: "tn-edge-created-vdcgroup",
				ID:   generator.MustGenerate("{urn:edgegateway}"),
			}}},
			mockQueryResponseStatus: 200,
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{Values: []itypes.APIResponseListVDCGroupDetails{{
				Name: "owner-vdcgroup",
			}}},
			mockListVDCGroupResponseStatus: 200,
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: "tn-edge-created-vdcgroup",
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockListT0ResponseStatus: 200,
			mockJobResponseStatus:    200,
			expectedErr:              false,
		},
		{
			name: "Valid Edge Gateway Creation with Bandwidth",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdc-bandwidth",
				Bandwidth: 25,
			},
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				Name: "owner-vdc-bandwidth",
				HREF: "https://api.example.com/api/vdc/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   "urn:vcloud:vdc:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockListVDCResponseStatus: 200,
			mockGetEdgeGatewayResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: "tn-edge-created-bandwidth",
			},
			mockGetEdgeGatewayResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				Name: "tn-edge-created-bandwidth",
				ID:   generator.MustGenerate("{urn:edgegateway}"),
			}}},
			mockQueryResponseStatus: 200,
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: "tn-edge-created-bandwidth",
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockJobResponseStatus:          200,
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockListT0ResponseStatus: 200,
			expectedErr:              false,
		},
		{
			name: "Failed to list VDC",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListVDCResponseStatus: http.StatusNotFound,
			expectedErr:               true,
		},
		{
			name: "Failed to list VDC Group",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListVDCGroupResponseStatus: http.StatusNotFound,
			expectedErr:                    true,
		},
		{
			name: "Failed no Vdc or VDCGroup Found",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockListVDCResponse:            &itypes.APIResponseListVDC{},
			mockListVDCResponseStatus:      200,
			expectedErr:                    true,
		},
		{
			name: "Failed Both VDCs and VDC Groups found",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListVDCGroupResponse: &itypes.APIResponseListVDCGroup{
				Values: []itypes.APIResponseListVDCGroupDetails{
					{
						Name: generator.MustGenerate("{word}"),
					},
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCGroupResponseStatus: 200,
			mockListVDCResponse: &itypes.APIResponseListVDC{
				Records: []itypes.APIResponseListVDCRecord{
					{
						Name: generator.MustGenerate("{word}"),
					},
					{
						Name: generator.MustGenerate("{word}"),
					},
				},
			},
			mockListVDCResponseStatus: 200,
			expectedErr:               true,
		},

		{
			name: "Failed to list T0",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListT0Response:       nil,
			mockListT0ResponseStatus: 401,
			expectedErr:              true,
		},
		{
			name: "T0 return 0 T0s",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListT0Response:       &itypes.APIResponseT0s{},
			mockListT0ResponseStatus: 200,
			expectedErr:              true,
		},
		{
			name: "Failed T0s > 1",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListT0Response: &itypes.APIResponseT0s{
				itypes.APIResponseT0{
					Type: "tier-0-vrf",
					Name: generator.MustGenerate("{resource_name:t0}"),
				},
				itypes.APIResponseT0{
					Type: "tier-0-vrf",
					Name: generator.MustGenerate("{resource_name:t0}"),
				},
			},
			mockListT0ResponseStatus: 200,
			expectedErr:              true,
		},
		{
			name: "Create Edge Gateway with SHARED T0",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdc-shared-t0",
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				Name: "owner-vdc-shared-t0",
				HREF: "https://api.example.com/api/vdc/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   "urn:vcloud:vdc:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockListVDCResponseStatus: 200,
			mockGetEdgeGatewayResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockGetEdgeGatewayResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   generator.MustGenerate("{urn:edgegateway}"),
			}}},
			mockQueryResponseStatus: 200,
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: generator.MustGenerate("{resource_name:edgegateway}"),
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockJobResponseStatus:          200,
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockListT0ResponseStatus:       200,
			expectedErr:                    false,
		},
		{
			name: "Create Edge Gateway with T0 not found",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    generator.MustGenerate("{resource_name:t0}"),
				Bandwidth: 25,
			},
			expectedErr: true,
		},
		{
			name: "Create Edge Gateway with invalid bandwidth values",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 500,
			},
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
		{
			name: "Failed extract job response",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockJobResponse:                &cav.CerberusJobAPIResponse{},
			mockJobResponseStatus:          200,
			expectedErr:                    true,
		},
		{
			name: "Failed to retrieve edge gateway by name after creation",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdc-get-fail",
				Bandwidth: 5,
			},

			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				Name: "owner-vdc-get-fail",
				HREF: "https://api.example.com/api/vdc/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   "urn:vcloud:vdc:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockListVDCResponseStatus:      200,
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockListT0ResponseStatus: 200,

			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: generator.MustGenerate("{resource_name:edgegateway}"),
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockJobResponseStatus: 200,

			mockGetEdgeGatewayResponseStatus: 404,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockQueryResponseStatus: 200,
			expectedErr:             true,
		},
		{
			name: "Failed to update edge gateway bandwidth after creation",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: "owner-vdc-bandwidth-update-fail",
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListVDCResponse: &itypes.APIResponseListVDC{Records: []itypes.APIResponseListVDCRecord{{
				Name: "owner-vdc-bandwidth-update-fail",
				HREF: "https://api.example.com/api/vdc/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
				ID:   "urn:vcloud:vdc:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockListVDCResponseStatus: 200,
			mockGetEdgeGatewayResponse: &itypes.APIResponseEdgegateway{
				ID:   generator.MustGenerate("{urn:edgegateway}"),
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockGetEdgeGatewayResponseStatus: 200,
			mockQueryResponse: &itypes.APIResponseQueryEdgeGateway{Record: []itypes.APIResponseQueryEdgeGatewayRecord{{
				HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			}}},
			mockQueryResponseStatus: 200,
			mockListT0Response: func() itypes.APIResponseT0s {
				return itypes.APIResponseT0s{{
					Type: "tier-0-vrf",
					Name: "prvrf01eocb0001234allsp01",
					Properties: itypes.APIResponseT0Properties{
						ClassOfService: "SHARED_STANDARD",
					},
				}}
			}(),
			mockJobResponse: &cav.CerberusJobAPIResponse{
				{
					Actions: []cav.CerberusJobAPIResponseAction{
						{
							Details: generator.MustGenerate("{resource_name:edgegateway}"),
							Name:    "Create Edge Gateway",
							Status:  "DONE",
						},
					},
					Name:        "Create Edge Gateway Job",
					Status:      "DONE",
					Description: "Edge Gateway created successfully",
				},
			},
			mockJobResponseStatus:                        200,
			mockListT0ResponseStatus:                     200,
			mockUpdateEdgeGatewayBandwidthResponseStatus: http.StatusBadRequest,
			mockListVDCGroupResponse:                     &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus:               200,
			expectedErr:                                  true,
		},
		{
			name: "Exceeding maximum edge gateways for T0",
			params: &types.ParamsCreateEdgeGateway{
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListT0Response: func() itypes.APIResponseT0s {
				countOfT0s := 5
				var t0 itypes.APIResponseT0
				t0.Name = "prvrf01eocb0001234allsp01"
				for range countOfT0s {
					edge := itypes.APIResponseT0Children{}
					_ = generator.Struct(&edge)
					t0.Children = append(t0.Children, edge)
				}
				return itypes.APIResponseT0s{
					t0,
				}
			}(),
			mockListT0ResponseStatus:       200,
			mockListVDCGroupResponse:       &itypes.APIResponseListVDCGroup{},
			mockListVDCGroupResponseStatus: 200,
			expectedErr:                    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockListVDCResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
			}

			if tt.mockListVDCGroupResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVDCGroup())
				ms.SetResponse(endpoints.ListVDCGroup(), tt.mockListVDCGroupResponse, &tt.mockListVDCGroupResponseStatus)
			}

			if tt.mockListT0ResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListT0())
				ms.SetResponse(endpoints.ListT0(), tt.mockListT0Response, &tt.mockListT0ResponseStatus)
				ms.CleanResponse(endpoints.GetEdgeGatewayServices())
				ms.SetResponse(endpoints.GetEdgeGatewayServices(), tt.mockListT0Response, &tt.mockListT0ResponseStatus)
			}

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.CreateEdgeGateway())
				ms.SetResponse(endpoints.CreateEdgeGateway(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockJobResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetJobCerberus())
				ms.SetResponse(endpoints.GetJobCerberus(), tt.mockJobResponse, &tt.mockJobResponseStatus)
			}

			if tt.mockGetEdgeGatewayResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetEdgeGateway())
				ms.SetResponse(endpoints.GetEdgeGateway(), tt.mockGetEdgeGatewayResponse, &tt.mockGetEdgeGatewayResponseStatus)
			}

			if tt.mockListVDCResponse != nil || tt.mockListVDCResponseStatus != 0 || tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				handler := func(w http.ResponseWriter, r *http.Request) {
					statusCode := http.StatusOK
					data := tt.mockQueryResponse

					if r.URL.Query().Get("type") == "edgeGateway" {
						if tt.mockQueryResponseStatus != 0 {
							statusCode = tt.mockQueryResponseStatus
						}
						data = tt.mockQueryResponse
					} else if r.URL.Query().Get("type") == "orgVdc" {
						if tt.mockListVDCResponseStatus != 0 {
							statusCode = tt.mockListVDCResponseStatus
						}
						data = tt.mockListVDCResponse
					}

					if statusCode >= http.StatusMultipleChoices {
						http.Error(w, http.StatusText(statusCode), statusCode)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("X-Cloud-Avenue-Mock", "true")
					if data == nil {
						data = map[string]any{}
					}
					if err := json.NewEncoder(w).Encode(data); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}

				ms.CleanResponse(endpoints.QueryEdgeGateway())
				ms.SetResponseFunc(endpoints.QueryEdgeGateway(), handler)
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponseFunc(endpoints.ListVDC(), handler)
			}

			if tt.mockUpdateEdgeGatewayBandwidthResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateEdgeGatewayBandwidth())
				ms.SetResponse(endpoints.UpdateEdgeGatewayBandwidth(), tt.mockUpdateEdgeGatewayBandwidthResponse, &tt.mockUpdateEdgeGatewayBandwidthResponseStatus)
			}

			result, err := eC.CreateEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error for params: %v", tt.params)
				assert.Nil(t, result, "Result should be nil for params: %v", tt.params)
			} else {
				assert.Nil(t, err, "Expected no error for params: %v", tt.params)
				assert.NotNil(t, result, "Result should not be nil for params: %v", tt.params)
			}
		})
	}
}

func TestListEdgeGateay(t *testing.T) {
	tests := []struct {
		name               string
		mockResponse       any
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name:               "Valid List Edge Gateways",
			mockResponseStatus: 200,
		},
		{
			name:               "Error 500",
			mockResponse:       struct{}{},
			mockResponseStatus: 500,
			expectedErr:        true,
		},
		{
			name:               "Error 401",
			mockResponse:       struct{}{},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.SetResponse(endpoints.ListEdgeGateway(), tt.mockResponse, &tt.mockResponseStatus)
			}

			result, err := eC.ListEdgeGateway(t.Context())
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil when error is expected")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.name)
				assert.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}

func TestUpdateEdgeGateway(t *testing.T) {
	tests := []struct {
		name   string
		params *types.ParamsUpdateEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockQueryResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: &types.ParamsUpdateEdgeGateway{
				ID:        generator.MustGenerate("{urn:edgegateway}"),
				Bandwidth: 25,
			},

			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: &types.ParamsUpdateEdgeGateway{
				Name:      generator.MustGenerate("{resource_name:edgegateway}"),
				Bandwidth: 25,
			},
			expectedErr: false,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: &types.ParamsUpdateEdgeGateway{
				ID:        "invalid-id",
				Bandwidth: 25,
			},
			expectedErr: true,
		},
		{
			name: "Edge Gateway Name not found",
			params: &types.ParamsUpdateEdgeGateway{
				Name:      generator.MustGenerate("{resource_name:edgegateway}"),
				Bandwidth: 25,
			},
			mockQueryResponseStatus: 404,
			expectedErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC, ms := newClient(t)

			if tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.UpdateEdgeGatewayBandwidth())
				ms.SetResponse(endpoints.UpdateEdgeGatewayBandwidth(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponseStatus != 0 {
				ms.CleanResponse(endpoints.QueryEdgeGateway())
				ms.SetResponse(endpoints.QueryEdgeGateway(), nil, &tt.mockQueryResponseStatus)
				ms.CleanResponse(endpoints.ListVDC())
				ms.SetResponse(endpoints.ListVDC(), nil, &tt.mockQueryResponseStatus)
			}

			_, err := eC.UpdateEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error for params: %v", tt.params)
			}
		})
	}
}
