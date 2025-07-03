package errors

import (
	"fmt"
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
}

// IsNotFound checks if the APIError indicates a "not found" error.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// Error returns the error message for APIError.
func (e *APIError) Error() string {
	if e == nil {
		return "nil APIError"
	}

	return fmt.Sprintf("[%s] request API error: %s (status code: %d, duration: %s, endpoint: %s)",
		e.Operation, e.Message, e.StatusCode, e.Duration, e.Endpoint,
	)
}
