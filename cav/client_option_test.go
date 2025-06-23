package cav

import (
	"testing"

	"resty.dev/v3"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	subclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/subClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

var getMockAuth = auth.NewMockAuth(map[string]string{
	"Authorization": "Bearer mock-token",
})

var getMockConsole = func() consoles.Console {
	c, _ := consoles.FindBySiteID("mock")
	return c
}

func TestWithCloudAvenueCredential_SetsCredentials(t *testing.T) {
	// Mock the NewCloudavenueCredential function to return a mock auth
	origNewCloudavenueCredential := auth.NewCloudavenueCredential
	auth.NewCloudavenueCredential = func(_ *resty.Client, _ consoles.Console, _ string, _ string, _ string) (auth.Auth, error) {
		return getMockAuth, nil
	}
	defer func() { auth.NewCloudavenueCredential = origNewCloudavenueCredential }()

	s := &settings{
		Organization: "mockorg001",
		SubClients:   make(map[subclient.Name]subclient.Client),
		Console:      getMockConsole(),
	}

	opt := WithCloudAvenueCredential("user", "pass")
	err := opt(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if s.SubClients[Vmware] == nil || s.SubClients[Cerberus] == nil {
		t.Fatalf("expected subclients to be initialized")
	}
}

func TestWithCloudAvenueCredential_SetsCredentialsInvalid(t *testing.T) {
	s := &settings{
		Organization: "test-org",
		SubClients:   make(map[subclient.Name]subclient.Client),
		Console:      getMockConsole(),
	}

	opt := WithCloudAvenueCredential("user", "")
	err := opt(s)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestWithConsole_SetsConsole(t *testing.T) {
	s := &settings{
		Organization: "mockorg001",
		SubClients:   make(map[subclient.Name]subclient.Client),
	}

	opt := withConsole()
	err := opt(s)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, s.Console.GetSiteName(), "Console Mock")
	assert.Equal(t, s.Console.GetSiteID(), consoles.Console("mock"))
	assert.Equal(t, s.Console.GetAPIVmwareEndpoint(), "http://mock.api/cloudapi")
}
