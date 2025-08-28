/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetOrganization(t *testing.T) {
	tests := []struct {
		name string
		// mock from infraAPI
		mockGetOrgResponse any
		mockGetOrgStatus   int
		// mock from VMware Cloud Director
		mockGetOrgsResponse any
		mockGetOrgsStatus   int

		expectErr bool
	}{
		{
			name: "Success Get Organization",
		},
		{
			name:             "Fail Get Organization from infraAPI",
			mockGetOrgStatus: 404,
			expectErr:        true,
		},
		{
			name:              "Fail Get Organization from VMware Cloud Director",
			mockGetOrgsStatus: 404,
			expectErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock from infraAPI
			if tt.mockGetOrgResponse != nil || tt.mockGetOrgStatus != 0 {
				// Clean all default mock responses
				endpoints.GetOrganization().CleanMockResponse()
				endpoints.GetOrganization().SetMockResponse(tt.mockGetOrgResponse, &tt.mockGetOrgStatus)
			}

			// Mock from VMware Cloud Director
			if tt.mockGetOrgsResponse != nil || tt.mockGetOrgsStatus != 0 {
				// Clean all default mock responses
				endpoints.GetOrganizationDetails().CleanMockResponse()
				endpoints.GetOrganizationDetails().SetMockResponse(tt.mockGetOrgsResponse, &tt.mockGetOrgsStatus)
			}

			client := newClient(t)

			resp, err := client.GetOrganization(t.Context())
			if tt.expectErr {
				assert.NotNil(t, err, "expected an error but got nil")
				return
			}
			assert.Nil(t, err, "expected no error but got: %v", err)
			assert.NotNil(t, resp, "expected a response but got nil")
			assert.NotEmpty(t, resp, "expected a non-empty response but got empty")
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	tests := []struct {
		name string
		// mock from infraAPI
		mockUpdateOrgResponse any
		mockUpdateOrgStatus   int

		// mock from infrAPI
		mockGetOrgResponse any
		mockGetOrgStatus   int

		// mock from VMware Cloud Director
		mockGetOrgsResponse any
		mockGetOrgsStatus   int

		params *types.ParamsUpdateOrganization

		expectErr bool
	}{
		{
			name: "Success Update Organization",
			params: &types.ParamsUpdateOrganization{
				FullName:            "New Org Name",
				Description:         func(s string) *string { return &s }("New Org Description"),
				CustomerMail:        "user@email.com",
				InternetBillingMode: "PAYG",
			},
		},
		{
			name: "Success Update Organization - Empty params",
			params: &types.ParamsUpdateOrganization{
				FullName:            "",
				Description:         func(s string) *string { return &s }(""),
				CustomerMail:        "",
				InternetBillingMode: "",
			},
			expectErr: false,
		},
		{
			name: "Fail - Do not retrieve Organization",
			params: &types.ParamsUpdateOrganization{
				FullName: "New Org Name",
			},
			mockGetOrgStatus: 400,
			expectErr:        true,
		},
		{
			name: "Fail - Empty Params",
			params: &types.ParamsUpdateOrganization{
				FullName:            "",
				Description:         nil,
				CustomerMail:        "",
				InternetBillingMode: "",
			},
			expectErr: true,
		},
		{
			name: "Fail Update Organization - Invalid InternetBillingMode",
			params: &types.ParamsUpdateOrganization{
				InternetBillingMode: "INVALID",
			},
			expectErr: true,
		},
		{
			name: "Fail Update Organization",
			params: &types.ParamsUpdateOrganization{
				FullName:            "New Org Name",
				Description:         func(s string) *string { return &s }("New Org Description"),
				CustomerMail:        "user@email.com",
				InternetBillingMode: "PAYG",
			},
			mockUpdateOrgStatus: 404,
			expectErr:           true,
		},
		{
			name: "Fail Get Organization Error after Update",
			params: &types.ParamsUpdateOrganization{
				FullName: "New Org Name",
			},
			mockGetOrgsStatus: 404,
			expectErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock update
			if tt.mockUpdateOrgResponse != nil || tt.mockUpdateOrgStatus != 0 {
				// Clean all default mock responses
				endpoints.UpdateOrganization().CleanMockResponse()
				endpoints.UpdateOrganization().SetMockResponse(tt.mockUpdateOrgResponse, &tt.mockUpdateOrgStatus)
			}

			// Mock get values from infraAPI
			if tt.mockGetOrgResponse != nil || tt.mockGetOrgStatus != 0 {
				// Clean all default mock responses
				endpoints.GetOrganization().CleanMockResponse()
				endpoints.GetOrganization().SetMockResponse(tt.mockGetOrgResponse, &tt.mockGetOrgStatus)
			}

			// Mock get values from VMware Cloud Director
			if tt.mockGetOrgsResponse != nil || tt.mockGetOrgsStatus != 0 {
				// Clean all default mock responses
				endpoints.GetOrganizationDetails().CleanMockResponse()
				endpoints.GetOrganizationDetails().SetMockResponse(tt.mockGetOrgsResponse, &tt.mockGetOrgsStatus)
			}

			client := newClient(t)

			resp, err := client.UpdateOrganization(t.Context(), *tt.params)
			if tt.expectErr {
				assert.NotNil(t, err, "expected an error but got nil")
				return
			}
			assert.Nil(t, err, "expected no error but got: %v", err)
			assert.NotNil(t, resp, "expected a response but got nil")
			assert.NotEmpty(t, resp, "expected a non-empty response but got empty")
		})
	}
}
