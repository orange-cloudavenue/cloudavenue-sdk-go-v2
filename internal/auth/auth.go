package auth

import (
	"context"
)

// Auth implements methods required for authentication.
type Auth interface {
	// Headers returns headers that must be included in the http request.
	Headers() map[string]string

	// Refresh refreshes the authentication token.
	Refresh(context.Context) error

	// IsInitialized checks if the authentication is initialized.
	IsInitialized() bool
}
