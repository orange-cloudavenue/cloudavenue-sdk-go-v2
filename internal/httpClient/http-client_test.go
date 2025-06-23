package httpclient

import "testing"

func Test_NewHTTPClient(t *testing.T) {
	client := NewHTTPClient()
	if client == nil {
		t.Error("NewHTTPClient() returned nil")
	}
	if client.IsDebug() {
		t.Error("NewHTTPClient() should not be in debug mode by default")
	}
}
