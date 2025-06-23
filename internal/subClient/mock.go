package subclient

import (
	"context"

	"resty.dev/v3"

	"github.com/jarcoal/httpmock"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	httpclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/httpClient"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

var _ Client = &mockClient{}

type mockClient struct {
	client
}

type MockError struct {
	Message string `json:"message"`
}

var NewMockClient = func() Client {
	return &mockClient{}
}

func (m *mockClient) NewHTTPClient(_ context.Context) (*resty.Client, error) {
	if m.httpClient == nil {
		m.httpClient = httpclient.NewMockHTTPClient().
			SetBaseURL("https://mock-api.cloudavenue.com").
			SetHeaders(m.credential.Headers()).
			SetError(MockError{})
		httpmock.ActivateNonDefault(m.httpClient.Client())
	}
	return m.httpClient, nil
}

func (m *mockClient) SetCredential(a auth.Auth) {
	m.credential = a
}

func (m *mockClient) SetConsole(c consoles.Console) {
	m.console = c
}

func (m *mockClient) ParseAPIError(resp *resty.Response) *errors.APIError {
	if resp == nil || !resp.IsError() {
		return nil
	}

	// Mock error parsing logic
	if err, ok := resp.Error().(*MockError); ok {
		return &errors.APIError{
			StatusCode: resp.StatusCode(),
			Message:    err.Message,
			Duration:   resp.Duration(),
			Endpoint:   resp.Request.URL,
		}
	}

	return nil
}
