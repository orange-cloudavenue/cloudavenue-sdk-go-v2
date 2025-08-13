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

	"github.com/orange-cloudavenue/common-go/generator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListStorageProfiles(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsListStorageProfiles
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "List Storage Profiles",
			params: types.ParamsListStorageProfiles{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			expectedErr: false,
		},
		{
			name: "Error 401 Unauthorized",
			params: types.ParamsListStorageProfiles{
				ID: generator.MustGenerate("{urn:vdc}"),
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.ListStorageProfiles().CleanMockResponse()
				endpoints.ListStorageProfiles().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListStorageProfile(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.StorageProfiles, "Storage profiles should not be empty")
			for _, sp := range resp.StorageProfiles {
				assert.NotEmpty(t, sp.ID, "Storage profile ID should not be empty")
				assert.NotEmpty(t, sp.Class, "Storage profile Class should not be empty")
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
		expectedErr        bool
	}{
		{
			name: "Add Storage Profile",
			params: types.ParamsAddStorageProfile{
				VdcID:   generator.MustGenerate("{urn:vdc}"),
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
				VdcID:   generator.MustGenerate("{urn:vdc}"),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.UpdateVdc().CleanMockResponse()
				endpoints.UpdateVdc().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
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
