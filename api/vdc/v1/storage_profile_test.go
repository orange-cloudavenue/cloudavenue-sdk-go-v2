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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
)

func TestListStorageProfiles(t *testing.T) {
	tests := []struct {
		name   string
		params types.ParamsListStorageProfile
		// Mock response for ListStorageProfile endpoint
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name:        "List All Storage Profiles",
			params:      types.ParamsListStorageProfile{},
			expectedErr: false,
		},
		{
			name: "List Storage Profiles by Storage Profile Name",
			params: types.ParamsListStorageProfile{
				Name: "gold",
			},
			expectedErr: false,
		},
		{
			name: "List Storage Profiles by Storage Profile ID",
			params: types.ParamsListStorageProfile{
				ID: generator.MustGenerate("{urn:vdcstorageProfile}"),
			},
			expectedErr: false,
		},
		{
			name: "Error to combine Storage Profile ID and Name",
			params: types.ParamsListStorageProfile{
				ID:   generator.MustGenerate("{urn:vdcstorageProfile}"),
				Name: "gold",
			},
			expectedErr: false,
		},
		{
			name: "List Storage Profiles by VDC Name",
			params: types.ParamsListStorageProfile{
				VdcName: "my-vdc",
			},
			expectedErr: false,
		},
		{
			name: "List Storage Profiles by VDC ID",
			params: types.ParamsListStorageProfile{
				VdcId: generator.MustGenerate("{urn:vdc}"),
			},
			expectedErr: false,
		},
		{
			name: "Error wrong Storage Profile ID",
			params: types.ParamsListStorageProfile{
				ID: "urn:vcloud:vdcstorageProfile:f98f6819-2355-478e-a8ee-4442a9dafdcg",
			},
			expectedErr: true,
		},
		{
			name: "Error wrong VDC ID",
			params: types.ParamsListStorageProfile{
				VdcId: "urn:vcloud:vdc:f98f6819-2355-478e-a8ee-4442a9dafdcg",
			},
			expectedErr: true,
		},
		{
			name: "Error api response return an empty HREF for Storage Profile ID",
			params: types.ParamsListStorageProfile{
				ID: generator.MustGenerate("{urn:vdcstorageProfile}"),
			},
			mockResponse: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:      "", // Empty HREF to simulate error
						Name:      "platinum3k_r1",
						IsEnabled: true,
					},
				},
			},
			mockResponseStatus: 200,
			expectedErr:        true,
		},
		{
			name: "Error api response return an empty HREF for VDC ID",
			params: types.ParamsListStorageProfile{
				VdcId: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponse: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:      generator.MustGenerate("{href_uuid}"),
						VdcId:     "", // Empty VdcId to simulate error
						Name:      "platinum3k_r1",
						IsEnabled: true,
					},
				},
			},
			mockResponseStatus: 200,
			expectedErr:        true,
		},

		{
			name: "Error 400 Bad Request",
			params: types.ParamsListStorageProfile{
				VdcId: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 400,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.ListStorageProfile().CleanMockResponse()
				endpoints.ListStorageProfile().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListStorageProfile(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.VDCS, "Storage profiles should not be empty")
			for _, spVDC := range resp.VDCS {
				assert.NotEmpty(t, spVDC.ID, "VDC ID should not be empty")
				assert.NotEmpty(t, spVDC.Name, "VDC Name should not be empty")
				for i := range spVDC.StorageProfiles {
					assert.NotEmpty(t, spVDC.StorageProfiles[i].ID, "Storage profile ID should not be empty")
					assert.NotEmpty(t, spVDC.StorageProfiles[i].Class, "Storage profile Class should not be empty")
					assert.NotEmpty(t, spVDC.StorageProfiles[i].Limit, "Storage profile Limit should not be empty")
					assert.GreaterOrEqual(t, spVDC.StorageProfiles[i].Used, 0, "Storage profile Used should be greater than or equal to 0")
				}
			}
		})
	}
}

func TestAddStorageProfile(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsAddStorageProfile
		mockResponseStatus int
		mockResponse       any
		// Specific mock response for ListVDC endpoint
		mockResponseVDC       any
		mockResponseVDCStatus int
		expectedErr           bool
	}{
		{
			name: "Add Storage Profile",
			params: types.ParamsAddStorageProfile{
				VdcId:   generator.MustGenerate("{urn:vdc}"),
				VdcName: "my-vdc",
				StorageProfiles: []types.ParamsCreateVDCStorageProfile{
					{
						Class:   "gold",
						Limit:   500,
						Default: false,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "Error 401 Unauthorized",
			params: types.ParamsAddStorageProfile{
				VdcId:   generator.MustGenerate("{urn:vdc}"),
				VdcName: "my-vdc",
				StorageProfiles: []types.ParamsCreateVDCStorageProfile{
					{
						Class:   "gold",
						Limit:   500,
						Default: false,
					},
				},
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "Error 404 VDC Not Found",
			params: types.ParamsAddStorageProfile{
				VdcId:   generator.MustGenerate("{urn:vdc}"),
				VdcName: "my-vdc",
			},
			mockResponseVDCStatus: 404,
			expectedErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdc().CleanMockResponse()
				endpoints.UpdateVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockResponseVDCStatus != 0 {
				endpoints.ListVdc().CleanMockResponse()
				endpoints.ListVdc().SetMockResponse(tt.mockResponseVDC, &tt.mockResponseVDCStatus)
			}

			client := newClient(t)

			err := client.AddStorageProfile(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestDeleteStorageProfile(t *testing.T) {
	tests := []struct {
		name                                 string
		params                               types.ParamsDeleteStorageProfile
		mockResponseStatus                   int
		mockResponse                         any
		mockResponseListStorageProfileStatus int
		mockResponseListStorageProfile       any
		expectedErr                          bool
	}{
		// Successful deletion of a Storage Profile (not default and no last storage profile)
		{
			name: "Delete Storage Profile",
			params: types.ParamsDeleteStorageProfile{
				VdcName:         "vdc1",
				VdcId:           "urn:vcloud:vdc:5ec9d15c-dc05-4a0f-8340-b10b18cda038",
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "gold",
						IsDefaultStorageProfile: false,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   500,
						Used:                    0,
					},
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "silver",
						IsDefaultStorageProfile: true,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   1000,
						Used:                    0,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          false,
		},

		// Error in list Storage Profiles
		{
			name: "Error 400 for VDC List",
			params: types.ParamsDeleteStorageProfile{
				VdcName:         generator.MustGenerate("{word}"),
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfileStatus: 400,
			expectedErr:                          true,
		},

		// Error to delete a Storage Profile Class where VDC is not found
		{
			name: "Error 404 Not Found",
			params: types.ParamsDeleteStorageProfile{
				VdcId:           "urn:vcloud:vdc:5ec9d15c-dc05-4a0f-8340-b10b18cda038",
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseStatus: 404,
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "gold",
						IsDefaultStorageProfile: false,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   500,
						Used:                    0,
					},
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "silver",
						IsDefaultStorageProfile: true,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   1000,
						Used:                    0,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete a VDC with no storage profiles
		{
			name: "Error delete an empty list of storage profiles",
			params: types.ParamsDeleteStorageProfile{
				VdcId:           generator.MustGenerate("{urn:vdc}"),
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete a unique Storage Profiles
		{
			name: "Error delete a unique Storage Profile",
			params: types.ParamsDeleteStorageProfile{
				VdcId:           generator.MustGenerate("{urn:vdc}"),
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "gold",
						IsDefaultStorageProfile: true,
						VdcName:                 generator.MustGenerate("{word}"),
						VdcId:                   generator.MustGenerate("{urn:vdc}"),
						Limit:                   500,
						Used:                    0,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete a default Storage Profile
		{
			name: "Error delete a default Storage Profile",
			params: types.ParamsDeleteStorageProfile{
				VdcName:         "vdc1",
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "gold",
						IsDefaultStorageProfile: true,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   500,
						Used:                    0,
					},
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "silver",
						IsDefaultStorageProfile: false,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   1000,
						Used:                    0,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete a Storage Profile Class not empty
		{
			name: "Error delete a Storage Profile Class not empty",
			params: types.ParamsDeleteStorageProfile{
				VdcName:         "vdc1",
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "gold",
						IsDefaultStorageProfile: false,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   500,
						Used:                    100, // Used is not zero, so it cannot be deleted
					},
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "silver",
						IsDefaultStorageProfile: true,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   1000,
						Used:                    100,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete no Storage Profile Class found in VDC
		{
			name: "Error delete no Storage Profile Class found in VDC",
			params: types.ParamsDeleteStorageProfile{
				VdcName:         "vdc1",
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "silver",
						IsDefaultStorageProfile: true,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   1000,
						Used:                    0,
					},
					{
						HREF:                    generator.MustGenerate("{href_uuid}"),
						Name:                    "bronze",
						IsDefaultStorageProfile: false,
						VdcName:                 "vdc1",
						VdcId:                   generator.MustGenerate("{url}/5ec9d15c-dc05-4a0f-8340-b10b18cda038"),
						Limit:                   300,
						Used:                    0,
					},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},

		// Error to delete a Storage Profile Name with Several VDC response List
		{
			name: "Error delete a Storage Profile Name with Several VDC response List",
			params: types.ParamsDeleteStorageProfile{
				VdcId:           generator.MustGenerate("{urn:vdc}"),
				StorageProfiles: []types.ParamsDeleteVDCStorageProfile{{Class: "gold"}},
			},
			mockResponseListStorageProfile: &itypes.ApiResponseListStorageProfiles{
				StorageProfiles: []itypes.ApiResponseListStorageProfile{
					{HREF: generator.MustGenerate("{href_uuid}"), Name: "gold", IsDefaultStorageProfile: true, VdcId: generator.MustGenerate("{urn:vdc}"), VdcName: generator.MustGenerate("{word}")},
					{HREF: generator.MustGenerate("{href_uuid}"), Name: "gold", IsDefaultStorageProfile: false, VdcId: generator.MustGenerate("{urn:vdc}"), VdcName: generator.MustGenerate("{word}")},
				},
			},
			mockResponseListStorageProfileStatus: 200,
			expectedErr:                          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdc().CleanMockResponse()
				endpoints.UpdateVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockResponseListStorageProfileStatus != 0 {
				endpoints.ListStorageProfile().CleanMockResponse()
				endpoints.ListStorageProfile().SetMockResponse(tt.mockResponseListStorageProfile, &tt.mockResponseListStorageProfileStatus)
			}

			client := newClient(t)

			err := client.DeleteStorageProfile(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}
