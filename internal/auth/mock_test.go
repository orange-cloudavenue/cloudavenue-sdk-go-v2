package auth

import (
	"reflect"
	"testing"
)

func TestNewMockAuth_Headers(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token",
		"X-Test":        "test-value",
	}
	mock := NewMockAuth(headers)
	got := mock.Headers()
	if !reflect.DeepEqual(got, headers) {
		t.Errorf("Headers() = %v, want %v", got, headers)
	}
}

func TestMockAuth_Refresh(t *testing.T) {
	mock := NewMockAuth(map[string]string{})
	err := mock.Refresh(t.Context())
	if err != nil {
		t.Errorf("Refresh() error = %v, want nil", err)
	}
}

func TestMockAuth_IsInitialized(t *testing.T) {
	mock := NewMockAuth(map[string]string{})
	if !mock.IsInitialized() {
		t.Errorf("IsInitialized() = false, want true")
	}
}
