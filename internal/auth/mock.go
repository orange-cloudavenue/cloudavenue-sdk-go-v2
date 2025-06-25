/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package auth

import "context"

var _ Auth = (*MockAuth)(nil)

// MockAuth is a mock implementation of the Auth interface for testing purposes.
type MockAuth struct {
	headers map[string]string
}

// NewMockAuth creates a new instance of MockAuth with the given headers.
func NewMockAuth(headers map[string]string) Auth {
	return &MockAuth{
		headers: headers,
	}
}

// Headers returns the headers for the mock authentication.
func (m *MockAuth) Headers() map[string]string {
	return m.headers
}

// Refresh is a no-op for the mock authentication.
func (m *MockAuth) Refresh(_ context.Context) error {
	// No operation for mock authentication
	return nil
}

// IsInitialized always returns true for the mock authentication.
func (m *MockAuth) IsInitialized() bool {
	return true
}
