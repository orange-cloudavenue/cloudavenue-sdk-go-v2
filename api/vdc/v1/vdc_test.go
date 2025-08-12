/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/utils"
)

func TestListVDC(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsListVDC
		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name: "List VDCs with no filters",
			params: ParamsListVDC{
				ID:   "",
				Name: "",
			},
			expectedErr: false,
		},
		{
			name: "List VDCs with filter by Name",
			params: ParamsListVDC{
				Name: "test-name",
			},
			expectedErr: false,
		},
		{
			name: "List VDCs with filter by ID",
			params: ParamsListVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			expectedErr: false,
		},
		{
			name: "Error 401 Unauthorized",
			params: ParamsListVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.ListVdc()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListVDC(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
			assert.NotEmpty(t, resp.VDCS, "Expected VDCs to be not empty")
		})
	}
}

func TestGetVDC(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsGetVDC
		mockResponse       any
		mockResponseStatus int

		mockListVDCResponse       any
		mockListVDCResponseStatus int

		mockGetMetadataResponseStatus int

		expectedErr bool
	}{
		{
			name: "Get VDC by ID",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			expectedErr: false,
		},
		{
			name: "Get VDC by Name",
			params: ParamsGetVDC{
				Name: "test-name",
			},
			expectedErr: false,
		},
		{
			name: "Error 401 Unauthorized",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "failed to get VDC details",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "List VDCs response with no VDCs",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockListVDCResponse:       apiResponseListVDC{Records: []apiResponseListVDCRecord{}},
			mockListVDCResponseStatus: 200,
			expectedErr:               true,
		},
		{
			name: "Failed to list VDCs",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockListVDCResponseStatus: 401,
			expectedErr:               true,
		},
		{
			name: "Failed Get VDC Metadata",
			params: ParamsGetVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockGetMetadataResponseStatus: 401,
			expectedErr:                   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				endpoints.GetVdc().CleanMockResponse()
				endpoints.GetVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListVDCResponseStatus != 0 {
				endpoints.ListVdc().CleanMockResponse()
				endpoints.ListVdc().SetMockResponse(tt.mockListVDCResponse, &tt.mockListVDCResponseStatus)
			}

			if tt.mockGetMetadataResponseStatus != 0 {
				endpoints.GetVdcMetadata().CleanMockResponse()
				endpoints.GetVdcMetadata().SetMockResponse(nil, &tt.mockGetMetadataResponseStatus)
			}
			client := newClient(t)

			resp, err := client.GetVDC(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestCreateVDC(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsCreateVDC
		mockResponse       any
		mockResponseStatus int

		mockGetVDCResponse       any
		mockGetVDCResponseStatus int

		expectedErr bool
	}{
		{
			name: "Create VDC with valid parameters",
			params: ParamsCreateVDC{
				Name:                "test-vdc",
				Description:         "Test VDC",
				ServiceClass:        "STD",
				BillingModel:        "PAYG",
				DisponibilityClass:  "ONE-ROOM",
				StorageBillingModel: "PAYG",
				Vcpu:                5,
				Memory:              16,
				StorageProfiles: []ParamsCreateVDCStorageProfile{
					{
						Class:   "silver",
						Limit:   100,
						Default: true,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Create VDC with missing required parameters",
			params: ParamsCreateVDC{
				Name: "test-vdc",
				// Missing other required fields
			},
			expectedErr: true,
		},
		{
			name: "Create VDC with no default storage profile",
			params: ParamsCreateVDC{
				Name:                "test-vdc",
				Description:         "Test VDC",
				ServiceClass:        "STD",
				BillingModel:        "PAYG",
				DisponibilityClass:  "ONE-ROOM",
				StorageBillingModel: "PAYG",
				Vcpu:                5,
				Memory:              16,
				StorageProfiles: []ParamsCreateVDCStorageProfile{
					{
						Class: "silver",
						Limit: 100,
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "Error 401 Unauthorized",
			params: ParamsCreateVDC{
				Name:                "test-vdc",
				Description:         "Test VDC",
				ServiceClass:        "STD",
				BillingModel:        "PAYG",
				DisponibilityClass:  "ONE-ROOM",
				StorageBillingModel: "PAYG",
				Vcpu:                5,
				Memory:              16,
				StorageProfiles: []ParamsCreateVDCStorageProfile{
					{
						Class:   "silver",
						Limit:   100,
						Default: true,
					},
				},
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "Failed to get VDC details after creation",
			params: ParamsCreateVDC{
				Name:                "test-vdc",
				Description:         "Test VDC",
				ServiceClass:        "STD",
				BillingModel:        "PAYG",
				DisponibilityClass:  "ONE-ROOM",
				StorageBillingModel: "PAYG",
				Vcpu:                5,
				Memory:              16,
				StorageProfiles: []ParamsCreateVDCStorageProfile{
					{
						Class:   "silver",
						Limit:   100,
						Default: true,
					},
				},
			},
			mockResponseStatus:       201,
			mockGetVDCResponseStatus: 401,
			expectedErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				endpoints.CreateVdc().CleanMockResponse()
				endpoints.CreateVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockGetVDCResponse != nil || tt.mockGetVDCResponseStatus != 0 {
				endpoints.GetVdc().CleanMockResponse()
				endpoints.GetVdc().SetMockResponse(tt.mockGetVDCResponse, &tt.mockGetVDCResponseStatus)
			}

			client := newClient(t)

			resp, err := client.CreateVDC(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
			assert.NotEmpty(t, resp.ID, "Expected VDC ID to be not empty")
		})
	}
}

func TestUpdateVDC(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsUpdateVDC
		mockResponse       any
		mockResponseStatus int

		mockGetVDCResponseStatus int

		expectedErr bool
	}{
		{
			name: "Update VDC with valid parameters",
			params: ParamsUpdateVDC{
				Name:        "updated-vdc",
				Description: utils.ToPTR("Updated VDC"),
			},
			expectedErr: false,
		},
		{
			name: "Failed to retrieve VDC with valid parameters VCPU",
			params: ParamsUpdateVDC{
				Name: "updated-vdc",
				Vcpu: utils.ToPTR(10),
			},
			expectedErr:              true,
			mockGetVDCResponseStatus: 404,
		},
		{
			name: "Update VDC with valid parameters Memory",
			params: ParamsUpdateVDC{
				Name:   "updated-vdc",
				Memory: utils.ToPTR(16),
			},
			expectedErr: false,
		},
		{
			name:   "Update VDC with missing required parameters",
			params: ParamsUpdateVDC{
				// Missing other required fields
			},
			expectedErr: true,
		},
		{
			name: "Update VDC with valid parameters VCPU",
			params: ParamsUpdateVDC{
				Name: "updated-vdc",
				Vcpu: utils.ToPTR(10),
			},
			expectedErr: false,
		},
		{
			name: "Error 404 Not Found",
			params: ParamsUpdateVDC{
				Name:        "updated-vdc",
				Description: utils.ToPTR("Updated VDC"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				endpoints.UpdateVdc().CleanMockResponse()
				endpoints.UpdateVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockGetVDCResponseStatus != 0 {
				endpoints.GetVdc().CleanMockResponse()
				endpoints.GetVdc().SetMockResponse(nil, &tt.mockGetVDCResponseStatus)
			}

			client := newClient(t)

			err := client.UpdateVDC(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestDeleteVDC(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsDeleteVDC
		mockResponse       any
		mockResponseStatus int

		mockGetVDCResponseStatus int

		expectedErr bool
	}{
		{
			name: "Delete VDC with valid ID",
			params: ParamsDeleteVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			expectedErr: false,
		},
		{
			name: "Error 404 Not Found",
			params: ParamsDeleteVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Error 404 on Get VDC",
			params: ParamsDeleteVDC{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockGetVDCResponseStatus: 404,
			expectedErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				endpoints.DeleteVdc().CleanMockResponse()
				endpoints.DeleteVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockGetVDCResponseStatus != 0 {
				endpoints.ListVdc().CleanMockResponse()
				endpoints.ListVdc().SetMockResponse(nil, &tt.mockGetVDCResponseStatus)
			}

			client := newClient(t)

			err := client.DeleteVDC(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
		})
	}
}
