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

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
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
			expectedErr: false,
			mockResponse: &types.ModelEdgeGateway{
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
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
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
			ep := endpoints.GetEdgeGateway()
			epQuery := endpoints.QueryEdgeGateway()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
				// If we expect a valid response, we need to set the mock response
				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				// If we expect a query response, we need to set the mock response for the
				mock.SetMockResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)

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
	mC, err := mock.NewClient()
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
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Mock the QueryEdgeGateway endpoint
	epQuery := endpoints.QueryEdgeGateway()

	tests := []struct {
		name        string
		edgeName    string
		queryResp   *itypes.ApiResponseQueryEdgeGateway
		queryStatus int
		expectedID  string
		expectedErr bool
	}{
		{
			name:     "Valid Edge Gateway Name",
			edgeName: generator.MustGenerate("{resource_name:edgegateway}"),
			queryResp: &itypes.ApiResponseQueryEdgeGateway{
				Record: []itypes.ApiResponseQueryEdgeGatewayRecord{
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
			mock.SetMockResponse(epQuery, tt.queryResp, &tt.queryStatus)

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
			expectedErr:        false,
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
			epDelete := endpoints.DeleteEdgeGateway()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", epDelete.Name, tt.mockResponseStatus)
				// If we expect a valid response, we need to set the mock response
				mock.SetMockResponse(epDelete, tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockQueryResponseStatus != 0 {
				t.Logf("Setting mock response for query endpoint %s with status %d", epQuery.Name, tt.mockQueryResponseStatus)
				// If we expect a query response, we need to set the mock response for the
				mock.SetMockResponse(epQuery, nil, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)
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
	eC := newClient(t)

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

		mockListT0Response       any
		mockListT0ResponseStatus int

		mockUpdateEdgeGatewayBandwidthResponse       any
		mockUpdateEdgeGatewayBandwidthResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid Edge Gateway Creation",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
			},
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
			expectedErr:           false,
		},
		{
			name: "Valid Edge Gateway Creation with Bandwidth",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
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
			expectedErr:           false,
		},
		{
			name: "Failed to list T0",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
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
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListT0Response:       &itypes.ApiResponseT0s{},
			mockListT0ResponseStatus: 200,
			expectedErr:              true,
		},
		{
			name: "Failed T0s > 1",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockListT0Response: &itypes.ApiResponseT0s{
				itypes.ApiResponseT0{
					Type: "tier-0-vrf",
					Name: generator.MustGenerate("{resource_name:t0}"),
				},
				itypes.ApiResponseT0{
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
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListT0Response: func() itypes.ApiResponseT0s {
				t0 := itypes.ApiResponseT0{}
				_ = generator.Struct(&t0)
				t0.Name = "prvrf01eocb0001234allsp01"
				return itypes.ApiResponseT0s{
					t0,
				}
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
			mockJobResponseStatus:    200,
			mockListT0ResponseStatus: 200,
			expectedErr:              false,
		},
		{
			name: "Create Edge Gateway with T0 not found",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    generator.MustGenerate("{resource_name:t0}"),
				Bandwidth: 25,
			},
			expectedErr: true,
		},
		{
			name: "Create Edge Gateway with invalid bandwidth values",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 500,
			},
			expectedErr: true,
		},
		{
			name: "Failed extract job response",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 25,
			},
			mockJobResponse:       &cav.CerberusJobAPIResponse{},
			mockJobResponseStatus: 200,
			expectedErr:           true,
		},
		{
			name: "Failed to retrieve edge gateway by name after creation",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				Bandwidth: 5,
			},
			mockGetEdgeGatewayResponseStatus: 404,
			expectedErr:                      true,
		},
		{
			name: "Failed to update edge gateway bandwidth after creation",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListT0Response: func() itypes.ApiResponseT0s {
				t0 := itypes.ApiResponseT0{}
				_ = generator.Struct(&t0)
				t0.Name = "prvrf01eocb0001234allsp01"
				return itypes.ApiResponseT0s{
					t0,
				}
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
			expectedErr: true,
		},
		{
			name: "Exceeding maximum edge gateways for T0",
			params: &types.ParamsCreateEdgeGateway{
				OwnerType: "vdc",
				OwnerName: generator.MustGenerate("{word}"),
				T0Name:    "prvrf01eocb0001234allsp01",
				Bandwidth: 25,
			},
			mockListT0Response: func() itypes.ApiResponseT0s {
				countOfT0s := 5
				var t0 itypes.ApiResponseT0
				t0.Name = "prvrf01eocb0001234allsp01"
				for i := 0; i < countOfT0s; i++ {
					edge := itypes.ApiResponseT0Children{}
					_ = generator.Struct(&edge)
					t0.Children = append(t0.Children, edge)
				}
				return itypes.ApiResponseT0s{
					t0,
				}
			}(),
			mockListT0ResponseStatus: 200,
			expectedErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.CreateEdgeGateway()
			defer ep.RestoreMockResponse()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			epJob := endpoints.GetJobCerberus()
			defer epJob.RestoreMockResponse()
			if tt.mockJobResponse != nil || tt.mockJobResponseStatus != 0 {
				t.Logf("Setting mock job response for endpoint %s with status %d", epJob.Name, tt.mockJobResponseStatus)
				mock.SetMockResponse(epJob, tt.mockJobResponse, &tt.mockJobResponseStatus)
			}

			epGetEdge := endpoints.GetEdgeGateway()
			defer epGetEdge.RestoreMockResponse()
			if tt.mockGetEdgeGatewayResponse != nil || tt.mockGetEdgeGatewayResponseStatus != 0 {
				t.Logf("Setting mock get edge gateway response for endpoint %s with status %d", epGetEdge.Name, tt.mockGetEdgeGatewayResponseStatus)
				mock.SetMockResponse(epGetEdge, tt.mockGetEdgeGatewayResponse, &tt.mockGetEdgeGatewayResponseStatus)
			}

			epListT0 := endpoints.ListT0()
			defer epListT0.RestoreMockResponse()
			if tt.mockListT0Response != nil || tt.mockListT0ResponseStatus != 0 {
				t.Logf("Setting mock list T0 response for endpoint %s with status %d", epListT0.Name, tt.mockListT0ResponseStatus)
				mock.SetMockResponse(epListT0, tt.mockListT0Response, &tt.mockListT0ResponseStatus)
			}

			epUpdateBandwidth := endpoints.UpdateEdgeGatewayBandwidth()
			defer epUpdateBandwidth.RestoreMockResponse()
			if tt.mockUpdateEdgeGatewayBandwidthResponse != nil || tt.mockUpdateEdgeGatewayBandwidthResponseStatus != 0 {
				t.Logf("Setting mock update edge gateway bandwidth response for endpoint %s with status %d", epUpdateBandwidth.Name, tt.mockUpdateEdgeGatewayBandwidthResponseStatus)
				mock.SetMockResponse(epUpdateBandwidth, tt.mockUpdateEdgeGatewayBandwidthResponse, &tt.mockUpdateEdgeGatewayBandwidthResponseStatus)
			}

			eC := newClient(t)

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
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
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
			ep := endpoints.ListEdgeGateway()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

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
		{
			name: "Error 500",
			params: &types.ParamsUpdateEdgeGateway{
				ID:        generator.MustGenerate("{urn:edgegateway}"),
				Bandwidth: 10,
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 500,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Error 401",
			params: &types.ParamsUpdateEdgeGateway{
				ID:        generator.MustGenerate("{urn:edgegateway}"),
				Bandwidth: 10,
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.UpdateEdgeGatewayBandwidth()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockQueryResponseStatus != 0 {
				t.Logf("Setting mock response for query endpoint %s with status %d", epQuery.Name, tt.mockQueryResponseStatus)
				// If we expect a query response, we need to set the mock response for the
				mock.SetMockResponse(epQuery, nil, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)

			_, err := eC.UpdateEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error for params: %v", tt.params)
			}
		})
	}
}
