/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
)

func TestDemoRequest(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	oC, err := New(mC)
	assert.Nil(t, err, "Error creating org client")

	ep, err := mock.GetEndpoint("GetOrganization", cav.MethodGET)
	if err != nil {
		t.Fatalf("Error getting endpoint: %v", err)
	}
	defer ep.CleanMockResponse()

	tests := []struct {
		name           string
		orgID          string
		expectedErr    bool
		expectedResp   *OrgResponse
		expectedStatus int
	}{
		{
			name:        "Valid Org ID",
			orgID:       "urn:vcloud:org:fd0847c0-1e81-4eb6-a3b9-b24a7aff6121",
			expectedErr: false,
			expectedResp: &OrgResponse{
				Name:          "Test Org",
				ID:            "urn:vcloud:org:fd0847c0-1e81-4eb6-a3b9-b24a7aff6121",
				Description:   "This is a test organization",
				CanManageOrgs: true,
				CanPublish:    true,
				CatalogCount:  5,
				DiskCount:     10,
				DisplayName:   "Test Organization",
				IsEnabled:     true,
			},
			expectedStatus: 200,
		},
		{
			name:        "Invalid Org ID",
			orgID:       "badid",
			expectedErr: true,
		},
		{
			name:           "Error 500",
			orgID:          "urn:vcloud:org:fd0847c0-1e81-4eb6-a3b9-b24a7aff6121",
			expectedErr:    true,
			expectedResp:   nil,
			expectedStatus: 500,
		},
		{
			name:           "Error 404",
			orgID:          "urn:vcloud:org:nonexistent",
			expectedErr:    true,
			expectedResp:   nil,
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Default().Println("Running test case:", tt.name)
			if tt.expectedResp != nil || tt.expectedStatus != 0 {
				// If we expect a valid response, we need to set the mock response
				mock.SetMockResponse(ep, tt.expectedResp, &tt.expectedStatus)
			}

			resp, err := oC.DemoRequest(t.Context(), tt.orgID)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error for orgID %s", tt.orgID)
				assert.Nil(t, resp, "Response should be nil for orgID %s", tt.orgID)
				t.Log("Error received:", err)
			} else {
				assert.Nil(t, err, "Expected no error for orgID %s", tt.orgID)
				assert.NotNil(t, resp, "Response should not be nil for orgID %s", tt.orgID)
				assert.Equal(t, tt.expectedResp.Name, resp.Name, "Expected Name to match for orgID %s", tt.orgID)
				assert.Equal(t, tt.expectedResp.ID, resp.ID, "Expected ID to match for orgID %s", tt.orgID)
				assert.Equal(t, tt.expectedResp.Description, resp.Description, "Expected Description to match for orgID %s", tt.orgID)
			}

			ep.CleanMockResponse()
		})
	}
}
