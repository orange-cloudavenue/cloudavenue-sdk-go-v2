package cav

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	_, err := newMockClient()
	assert.Nil(t, err, "Error creating mock client")
}

func Test_NewClient_InvalidOrganization(t *testing.T) {
	// Example test case for NewClient with an invalid organization
	_, err := NewClient("invalid_org")
	if err == nil {
		t.Fatal("Expected error for invalid organization, got nil")
	}
}

func Test_NewRequest_WithoutAuth(t *testing.T) {
	client, err := NewClient(mockOrg)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	_, err = client.NewRequest(t.Context(), ClientVmware)
	if err == nil {
		t.Fatal("Expected error for request without authentication, got nil")
	}
	if err.Error() != "invalid client vmware" {
		t.Fatalf("Expected error message 'invalid client vmware', got '%v'", err.Error())
	}
}

func Test_NewRequest(t *testing.T) {
	client, err := newMockClient()
	if err != nil {
		t.Fatalf("Error creating client with mock: %v", err)
	}

	_, err = client.NewRequest(t.Context(), ClientVmware)
	if err != nil {
		t.Fatalf("Error creating request with mock: %v", err)
	}
}
