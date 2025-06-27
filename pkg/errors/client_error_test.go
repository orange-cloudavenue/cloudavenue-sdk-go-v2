package errors

import "testing"

func TestClientError_Error(t *testing.T) {
	err := &ClientError{Message: "client error"}
	if got := err.Error(); got != "client error" {
		t.Errorf("ClientError.Error() = %q, want %q", got, "client error")
	}
	var nilErr *ClientError
	if got := nilErr.Error(); got != "nil ClientError" {
		t.Errorf("ClientError.Error() nil = %q, want %q", got, "nil ClientError")
	}
}
