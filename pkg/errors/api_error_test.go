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
	"testing"
	"time"
)

func TestAPIError_Error(t *testing.T) {
	apiErr := &APIError{
		Operation:  "API operation",
		StatusCode: 500,
		Message:    "internal error",
		Duration:   2 * time.Second,
		Endpoint:   "/test",
	}
	want := "[API operation] request API error: internal error (status code: 500, duration: 2s, endpoint: /test)"
	if got := apiErr.Error(); got != want {
		t.Errorf("APIError.Error() = %q, want %q", got, want)
	}
	var nilErr *APIError
	if got := nilErr.Error(); got != "nil APIError" {
		t.Errorf("APIError.Error() nil = %q, want %q", got, "nil APIError")
	}
}

func TestAPIError_IsNotFound(t *testing.T) {
	apiErr := &APIError{StatusCode: 404}
	if !apiErr.IsNotFound() {
		t.Error("APIError.IsNotFound() = false, want true")
	}
	apiErr.StatusCode = 500
	if apiErr.IsNotFound() {
		t.Error("APIError.IsNotFound() = true, want false")
	}
}
