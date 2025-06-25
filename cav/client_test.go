// SPDX-FileCopyrightText: Copyright (c) 2025 Orange
// SPDX-License-Identifier: Mozilla Public License 2.0
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.
//

package cav

import (
	"fmt"
	"testing"
)

func Test_NewClient(t *testing.T) {
	// Example test case for NewClient
	client, err := NewClient("mockorg001", WithMock())
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = client.NewRequest(t.Context(), mock)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
}

func Test_NewClient_InvalidOrganization(t *testing.T) {
	// Example test case for NewClient with an invalid organization
	_, err := NewClient("invalid_org")
	if err == nil {
		t.Fatal("Expected error for invalid organization, got nil")
	}
}

func Test_NewClient_WithMock(t *testing.T) {
	client, err := NewClient("mockorg001", WithMock())
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}

	_, err = client.NewRequest(t.Context(), mock)
	if err != nil {
		t.Fatalf("Error creating request with mock: %v", err)
	}
}

func Test_NewClient_WithInvalidMock(t *testing.T) {
	// Example test case for NewClient with an invalid mock
	_, err := NewClient("mockorg001", func(_ *settings) error {
		return fmt.Errorf("invalid client mock")
	})

	if err == nil {
		t.Fatal("Expected error for invalid mock client, got nil")
	}
	if err.Error() != "invalid client mock" {
		t.Fatalf("Expected error message 'invalid client mock', got '%v'", err.Error())
	}
}

func Test_NewRequest_WithoutAuth(t *testing.T) {
	client, err := NewClient("mockorg001")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = client.NewRequest(t.Context(), mock)
	if err == nil {
		t.Fatal("Expected error for request without authentication, got nil")
	}
	if err.Error() != "invalid client mock" {
		t.Fatalf("Expected error message 'invalid client mock', got '%v'", err.Error())
	}
}
func Test_NewRequest(t *testing.T) {
	client, err := NewClient("mockorg001", WithMock())
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}

	_, err = client.NewRequest(t.Context(), mock)
	if err != nil {
		t.Fatalf("Error creating request with mock: %v", err)
	}
}
