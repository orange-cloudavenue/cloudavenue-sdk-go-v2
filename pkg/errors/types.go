package errors

import (
	"fmt"
	"time"
)

type (
	ClientError struct {
		Message string
		APIErr  APIError
	}

	APIError struct {
		StatusCode int
		Message    string
		Duration   time.Duration
		Endpoint   string
	}
)

// Error returns the error message for ClientError.
func (e *ClientError) Error() string {
	if e == nil {
		return "nil ClientError"
	}
	return e.Message
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
	return fmt.Sprintf("request API error: %s (status code: %d, duration: %s, endpoint: %s)",
		e.Message, e.StatusCode, e.Duration, e.Endpoint,
	)
}
