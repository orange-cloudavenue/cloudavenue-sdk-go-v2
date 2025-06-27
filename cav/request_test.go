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

	"github.com/stretchr/testify/assert"
)

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
func Test_NewRequest_RequestOptionsError(t *testing.T) {
	client, err := NewClient("mockorg001", WithMock())
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}

	// Simulate a RequestOption that returns an error
	badOpt := func(_ *requestOption) error {
		return assert.AnError
	}

	_, err = client.NewRequest(t.Context(), mock, badOpt)
	if err == nil {
		t.Fatal("Expected error from bad request option, got nil")
	}
}

func Test_NewRequest_WithJobOpts_SubClientDoesNotImplementJobsClient(t *testing.T) {
	client, err := NewClient("mockorg001", WithMock())
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}
	_, err = client.NewRequest(t.Context(), "mock", WithJob())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not support job options")
}

func Test_NewRequest_WithJobOpts_SubClientImplementsJobsClient(t *testing.T) {
	client, err := NewClient("mockorg001", WithMockJob())
	if err != nil {
		t.Fatalf("Error creating client with mock job: %v", err)
	}

	// Create a request with job options
	req, err := client.NewRequest(t.Context(), mockJob, WithJob())
	if err != nil {
		t.Fatalf("Error creating request with job options: %v", err)
	}

	// Check if the request is created successfully
	assert.NotNil(t, req)
}
