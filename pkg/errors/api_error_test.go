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
	stderrors "errors"
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
		Method:     "POST",
	}
	want := "[API operation] request API error: internal error (method:POST status code: 500, duration: 2s, endpoint: /test)"
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

	apiErr.Err = ErrNotFound
	if !apiErr.IsNotFound() {
		t.Error("APIError.IsNotFound() with wrapped ErrNotFound = false, want true")
	}
}

func TestAPIError_IsUnauthorized(t *testing.T) {
	apiErr := &APIError{StatusCode: 401}
	if !apiErr.IsUnauthorized() {
		t.Error("APIError.IsUnauthorized() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsUnauthorized() {
		t.Error("APIError.IsUnauthorized() = true, want false")
	}
	apiErr.Err = ErrUnauthorized
	if !apiErr.IsUnauthorized() {
		t.Error("APIError.IsUnauthorized() with wrapped ErrUnauthorized = false, want true")
	}
}

func TestAPIError_IsForbidden(t *testing.T) {
	apiErr := &APIError{StatusCode: 403}
	if !apiErr.IsForbidden() {
		t.Error("APIError.IsForbidden() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsForbidden() {
		t.Error("APIError.IsForbidden() = true, want false")
	}
	apiErr.Err = ErrForbidden
	if !apiErr.IsForbidden() {
		t.Error("APIError.IsForbidden() with wrapped ErrForbidden = false, want true")
	}
}

func TestAPIError_IsConflict(t *testing.T) {
	apiErr := &APIError{StatusCode: 409}
	if !apiErr.IsConflict() {
		t.Error("APIError.IsConflict() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsConflict() {
		t.Error("APIError.IsConflict() = true, want false")
	}
	apiErr.Err = ErrConflict
	if !apiErr.IsConflict() {
		t.Error("APIError.IsConflict() with wrapped ErrConflict = false, want true")
	}
}

func TestAPIError_IsInternalServerError(t *testing.T) {
	apiErr := &APIError{StatusCode: 500}
	if !apiErr.IsInternalServerError() {
		t.Error("APIError.IsInternalServerError() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsInternalServerError() {
		t.Error("APIError.IsInternalServerError() = true, want false")
	}
	apiErr.Err = ErrInternalServerError
	if !apiErr.IsInternalServerError() {
		t.Error("APIError.IsInternalServerError() with wrapped ErrInternalServerError = false, want true")
	}
}

func TestAPIError_IsBadRequest(t *testing.T) {
	apiErr := &APIError{StatusCode: 400}
	if !apiErr.IsBadRequest() {
		t.Error("APIError.IsBadRequest() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsBadRequest() {
		t.Error("APIError.IsBadRequest() = true, want false")
	}
	apiErr.Err = ErrBadRequest
	if !apiErr.IsBadRequest() {
		t.Error("APIError.IsBadRequest() with wrapped ErrBadRequest = false, want true")
	}
}

func TestAPIError_IsMethodNotAllowed(t *testing.T) {
	apiErr := &APIError{StatusCode: 405}
	if !apiErr.IsMethodNotAllowed() {
		t.Error("APIError.IsMethodNotAllowed() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsMethodNotAllowed() {
		t.Error("APIError.IsMethodNotAllowed() = true, want false")
	}
	apiErr.Err = ErrMethodNotAllowed
	if !apiErr.IsMethodNotAllowed() {
		t.Error("APIError.IsMethodNotAllowed() with wrapped ErrMethodNotAllowed = false, want true")
	}
}

func TestAPIError_IsRequestTimeout(t *testing.T) {
	apiErr := &APIError{StatusCode: 408}
	if !apiErr.IsRequestTimeout() {
		t.Error("APIError.IsRequestTimeout() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsRequestTimeout() {
		t.Error("APIError.IsRequestTimeout() = true, want false")
	}
	apiErr.Err = ErrRequestTimeout
	if !apiErr.IsRequestTimeout() {
		t.Error("APIError.IsRequestTimeout() with wrapped ErrRequestTimeout = false, want true")
	}
}

func TestAPIError_IsTooManyRequests(t *testing.T) {
	apiErr := &APIError{StatusCode: 429}
	if !apiErr.IsTooManyRequests() {
		t.Error("APIError.IsTooManyRequests() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsTooManyRequests() {
		t.Error("APIError.IsTooManyRequests() = true, want false")
	}
	apiErr.Err = ErrTooManyRequests
	if !apiErr.IsTooManyRequests() {
		t.Error("APIError.IsTooManyRequests() with wrapped ErrTooManyRequests = false, want true")
	}
}

func TestAPIError_IsBadGateway(t *testing.T) {
	apiErr := &APIError{StatusCode: 502}
	if !apiErr.IsBadGateway() {
		t.Error("APIError.IsBadGateway() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsBadGateway() {
		t.Error("APIError.IsBadGateway() = true, want false")
	}
	apiErr.Err = ErrBadGateway
	if !apiErr.IsBadGateway() {
		t.Error("APIError.IsBadGateway() with wrapped ErrBadGateway = false, want true")
	}
}

func TestAPIError_IsServiceUnavailable(t *testing.T) {
	apiErr := &APIError{StatusCode: 503}
	if !apiErr.IsServiceUnavailable() {
		t.Error("APIError.IsServiceUnavailable() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsServiceUnavailable() {
		t.Error("APIError.IsServiceUnavailable() = true, want false")
	}
	apiErr.Err = ErrServiceUnavailable
	if !apiErr.IsServiceUnavailable() {
		t.Error("APIError.IsServiceUnavailable() with wrapped ErrServiceUnavailable = false, want true")
	}
}

func TestAPIError_IsGatewayTimeout(t *testing.T) {
	apiErr := &APIError{StatusCode: 504}
	if !apiErr.IsGatewayTimeout() {
		t.Error("APIError.IsGatewayTimeout() = false, want true")
	}
	apiErr.StatusCode = 200
	if apiErr.IsGatewayTimeout() {
		t.Error("APIError.IsGatewayTimeout() = true, want false")
	}
	apiErr.Err = ErrGatewayTimeout
	if !apiErr.IsGatewayTimeout() {
		t.Error("APIError.IsGatewayTimeout() with wrapped ErrGatewayTimeout = false, want true")
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	apiErr := &APIError{Err: ErrJobFailed}
	if !stderrors.Is(apiErr, ErrJobFailed) {
		t.Error("errors.Is(apiErr, ErrJobFailed) = false, want true")
	}
}
