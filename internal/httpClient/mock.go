package httpclient

// NewMockHTTPClient creates a new mock HTTP client for testing purposes.
import (
	"github.com/jarcoal/httpmock"
	"resty.dev/v3"
)

// NewMockHTTPClient creates a new mock HTTP client for testing purposes.
func NewMockHTTPClient() *resty.Client {
	client := NewHTTPClient()
	httpmock.ActivateNonDefault(client.Client())
	return client
}
