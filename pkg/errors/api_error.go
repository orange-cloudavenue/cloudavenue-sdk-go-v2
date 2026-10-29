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
	"fmt"
	"net/http"
	"time"
)

type APIError struct {
	// Operation is a short description of the operation that caused the error.
	// It helps identify which API operation failed.
	Operation string

	// StatusCode is the HTTP status code returned by the API.
	// It is used to determine the type of error that occurred.
	StatusCode int

	// StatusMessage is the HTTP status message returned by the API.
	// It provides additional context about the error.
	StatusMessage string

	// Message is the error message returned by the API.
	// It contains the specific error details provided by the API.
	Message string

	// Duration is the time taken for the API request to complete.
	Duration time.Duration

	// Endpoint is the API endpoint that was called when the error occurred.
	// It helps identify which specific API endpoint was involved in the error.
	// This is useful for debugging and logging purposes.
	Endpoint string

	// Method is the HTTP method used for the API request (e.g., GET, POST).
	Method string

	// Err carries sentinel/root cause information for errors.Is / errors.AsType usage.
	Err error
}

// IsNotFound checks if the APIError indicates a "not found" error.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound || stderrors.Is(e, ErrNotFound)
}

// IsBadRequest checks if the APIError indicates a "bad request" error.
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == http.StatusBadRequest || stderrors.Is(e, ErrBadRequest)
}

// IsUnauthorized checks if the APIError indicates a "unauthorized" error.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized || stderrors.Is(e, ErrUnauthorized)
}

// IsForbidden checks if the APIError indicates a "forbidden" error.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden || stderrors.Is(e, ErrForbidden)
}

// IsMethodNotAllowed checks if the APIError indicates a "method not allowed" error.
func (e *APIError) IsMethodNotAllowed() bool {
	return e.StatusCode == http.StatusMethodNotAllowed || stderrors.Is(e, ErrMethodNotAllowed)
}

// IsRequestTimeout checks if the APIError indicates a "request timeout" error.
func (e *APIError) IsRequestTimeout() bool {
	return e.StatusCode == http.StatusRequestTimeout || stderrors.Is(e, ErrRequestTimeout)
}

// IsConflict checks if the APIError indicates a "conflict" error.
func (e *APIError) IsConflict() bool {
	return e.StatusCode == http.StatusConflict || stderrors.Is(e, ErrConflict)
}

// IsTooManyRequests checks if the APIError indicates a "too many requests" error.
func (e *APIError) IsTooManyRequests() bool {
	return e.StatusCode == http.StatusTooManyRequests || stderrors.Is(e, ErrTooManyRequests)
}

// IsInternalServerError checks if the APIError indicates an "internal server error".
func (e *APIError) IsInternalServerError() bool {
	return e.StatusCode == http.StatusInternalServerError || stderrors.Is(e, ErrInternalServerError)
}

// IsBadGateway checks if the APIError indicates a "bad gateway" error.
func (e *APIError) IsBadGateway() bool {
	return e.StatusCode == http.StatusBadGateway || stderrors.Is(e, ErrBadGateway)
}

// IsServiceUnavailable checks if the APIError indicates a "service unavailable" error.
func (e *APIError) IsServiceUnavailable() bool {
	return e.StatusCode == http.StatusServiceUnavailable || stderrors.Is(e, ErrServiceUnavailable)
}

// IsGatewayTimeout checks if the APIError indicates a "gateway timeout" error.
func (e *APIError) IsGatewayTimeout() bool {
	return e.StatusCode == http.StatusGatewayTimeout || stderrors.Is(e, ErrGatewayTimeout)
}

// Unwrap returns underlying sentinel/root cause when present.
func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Error returns the error message for APIError.
func (e *APIError) Error() string {
	if e == nil {
		return "nil APIError"
	}

	return fmt.Sprintf(
		"[%s] request API error: %s (method:%s status code: %d, duration: %s, endpoint: %s)",
		e.Operation, e.Message, e.Method, e.StatusCode, e.Duration, e.Endpoint,
	)
}
