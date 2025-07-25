package edgegateway

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/common-go/generator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func TestGetNetworkServices(t *testing.T) {
	tests := []struct {
		name                    string
		params                  *ParamsEdgeGateway
		mockResponse            any
		mockResponseStatus      int
		mockQueryResponse       any
		mockQueryResponseStatus int
		expectedErr             bool
	}{
		{
			name: "Valid Edge Gateway services",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway services with name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			expectedErr: false,
		},
		{
			name:        "Simulate empty params",
			params:      &ParamsEdgeGateway{},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Simulate empty response",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponse:       &apiResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.GetNetworkServices()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				// Clean all default mock responses
				ep.CleanMockResponse()
				// Set the mock response
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			// Set up mock query response
			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				// Clean all default mock responses
				epQuery.CleanMockResponse()
				// Set the mock query response
				epQuery.SetMockResponse(tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)

			// Call the GetNetworkServices method
			result, err := eC.GetNetworkServices(t.Context(), *tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil: %v", result)
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
				assert.NotNil(t, result, "Result should not be nil: %v", tt.params)
			}
		})
	}
}

func TestGetNetworkServices_ContextDeadlineExceeded(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err = eC.GetNetworkServices(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgeGateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestEnableCloudAvenueServices(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Enable network services with valid ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable network services with valid Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable network services with empty params",
			params: ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.EnableCloudAvenueServices()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

			err := eC.EnableCloudAvenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
			}
		})
	}
}

func TestEnableCloudavenueServices_ContextDeadlineExceeded(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	err = eC.EnableCloudAvenueServices(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgeGateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestDisableCloudAvenueServices(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Disable network services with valid ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Disable network services with valid Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			expectedErr: false,
		},
		{
			name: "Disable network services with empty params",
			params: ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.DisableCloudAvenueServices()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

			err := eC.DisableCloudAvenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
			}

			ep.CleanMockResponse()
		})
	}
}
