package httpclient

import "testing"

func Test_NewMockHTTPClient(t *testing.T) {
	client := NewMockHTTPClient()
	if client == nil {
		t.Error("NewMockHTTPClient() returned nil")
	}
	if client.Client() == nil {
		t.Error("NewMockHTTPClient() did not create a valid HTTP client")
	}
}
