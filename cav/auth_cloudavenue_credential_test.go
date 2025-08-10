/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"testing"

	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/xlog"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

func Test_newCloudavenueCredential(t *testing.T) {
	console, _ := consoles.FindByOrganizationName(mockOrg)

	// Test cases for newCloudavenueCredential
	tests := []struct {
		name         string
		console      consoles.ConsoleName
		organization string
		username     string
		password     string
		expectError  bool
	}{
		{
			name:         "Valid credentials",
			console:      console,
			organization: mockOrg,
			username:     "test-user",
			password:     "test-pass",
			expectError:  false,
		},
		{
			name:         "Empty organization",
			console:      console,
			organization: "",
			username:     "test-user",
			password:     "test-pass",
			expectError:  true,
		},
		{
			name:         "Empty username",
			console:      console,
			organization: mockOrg,
			username:     "",
			password:     "test-pass",
			expectError:  true,
		},
		{
			name:         "Empty password",
			console:      console,
			organization: mockOrg,
			username:     "test-user",
			password:     "",
			expectError:  true,
		},
		{
			name:         "Bad org format",
			console:      console,
			organization: "bad-org-format",
			username:     "test-user",
			password:     "test-pass",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := newCloudavenueCredential(tt.console, tt.organization, tt.username, tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("newCloudavenueCredential() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && auth == nil {
				t.Error("Expected non-nil auth object")
			}
		})
	}
}

func Test_CloudavenueCredential_Refresh_with_Bearer(t *testing.T) {
	console, _ := consoles.FindByOrganizationName(mockOrg)

	auth := &cloudavenueCredential{
		logger:       xlog.GetGlobalLogger(),
		httpC:        httpclient.NewHTTPClient().SetBaseURL(console.GetAPIVCDEndpoint()),
		console:      console,
		organization: mockOrg,
		username:     "test-user",
		password:     "test-pass",
		bearer:       "test-bearer-token",
	}

	// Simulate the refresh process
	// Ignore the error for this test case, as we are just testing the method call
	_ = auth.Refresh(t.Context())
}
