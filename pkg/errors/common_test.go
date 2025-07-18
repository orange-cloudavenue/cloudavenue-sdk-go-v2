/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"errors"
	"testing"
)

// Mock types for testing

func TestIsAPIError(t *testing.T) {
	var apiErr error = &APIError{}
	var clientErr error = &ClientError{}
	var stdErr error = errors.New("standard error")
	var nilErr error

	if !IsAPIError(apiErr) {
		t.Errorf("IsAPIError should return true for APIError")
	}
	if IsAPIError(clientErr) {
		t.Errorf("IsAPIError should return false for ClientError")
	}
	if IsAPIError(stdErr) {
		t.Errorf("IsAPIError should return false for standard error")
	}
	if IsAPIError(nilErr) {
		t.Errorf("IsAPIError should return false for nil error")
	}
}

func TestIsClientError(t *testing.T) {
	var apiErr error = &APIError{}
	var clientErr error = &ClientError{}
	var stdErr error = errors.New("standard error")
	var nilErr error

	if !IsClientError(clientErr) {
		t.Errorf("IsClientError should return true for ClientError")
	}
	if IsClientError(apiErr) {
		t.Errorf("IsClientError should return false for APIError")
	}
	if IsClientError(stdErr) {
		t.Errorf("IsClientError should return false for standard error")
	}
	if IsClientError(nilErr) {
		t.Errorf("IsClientError should return false for nil error")
	}
}
