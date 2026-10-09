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
)

func Test_NewRequest_WithoutAuth(t *testing.T) {
	client, err := NewClient(mockOrg)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = client.NewRequest(t.Context(), &Endpoint{Backend: BackendVMware})
	if err == nil {
		t.Fatal("Expected error for request without authentication, got nil")
	}
	if err.Error() != "invalid client vmware" {
		t.Fatalf("Expected error message 'invalid client vmware', got '%v'", err.Error())
	}
}

func Test_NewRequest(t *testing.T) {
	client, err := NewClient(mockOrg, WithCloudAvenueCredential("mockuser", "mockpassword"))
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatalf("Error closing client with mock: %v", err)
		}
	}()

	_, err = client.NewRequest(t.Context(), &Endpoint{Backend: BackendVMware})
	if err != nil {
		t.Fatalf("Error creating request with mock: %v", err)
	}
}

func TestNewRequestAppliesEndpointHeaders(t *testing.T) {
	client, err := NewClient(mockOrg, WithCloudAvenueCredential("mockuser", "mockpassword"))
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Fatalf("Error closing client with mock: %v", err)
		}
	}()

	req, err := client.NewRequest(t.Context(), &Endpoint{
		Backend: BackendVMware,
		Headers: map[string]string{"Accept": "application/xml", "X-Endpoint": "endpoint"},
	})
	if err != nil {
		t.Fatalf("Error creating request with endpoint headers: %v", err)
	}
	if got := req.Header.Get("Accept"); got != "application/xml" {
		t.Fatalf("Accept header = %q, want %q", got, "application/xml")
	}
	if got := req.Header.Get("X-Endpoint"); got != "endpoint" {
		t.Fatalf("X-Endpoint header = %q, want %q", got, "endpoint")
	}
}

func TestNewRawRequest(t *testing.T) {
	c, err := NewClient(mockOrg, WithCloudAvenueCredential("mockuser", "mockpassword"))
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			t.Fatalf("Error closing client with mock: %v", err)
		}
	}()

	_, err = c.(*client).NewRawRequest(t.Context(), BackendVMware)
	if err != nil {
		t.Fatalf("Error creating raw request with mock: %v", err)
	}
}
