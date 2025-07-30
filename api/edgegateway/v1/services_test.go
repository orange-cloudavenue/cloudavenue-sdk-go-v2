package edgegateway

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/common-go/generator"
)

func TestGetEdgeGatewayServices(t *testing.T) {
	tests := []struct {
		name               string
		params             *ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid Edge Gateway services",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway services with name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockQueryResponse:       nil,
			mockQueryResponseStatus: http.StatusNotFound,
			expectedErr:             true,
		},
		{
			name:        "Simulate empty params",
			params:      &ParamsEdgeGateway{},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Simulate empty response",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       &apiResponseNetworkServices{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.GetEdgeGatewayServices()
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
			result, err := eC.GetServices(t.Context(), *tt.params)

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

func TestGetEdgeGatewayServices_ContextDeadlineExceeded(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err = eC.GetServices(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestEnableCloudavenueServices(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int

		mockQueryResponse       any
		mockQueryResponseStatus int

		expectedErr bool
	}{
		{
			name: "Enable network services with valid ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Enable network services with valid Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
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
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockQueryResponse:       nil,
			mockQueryResponseStatus: http.StatusNotFound,
			expectedErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.EnableCloudavenueServices()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			// Set up mock query response
			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				t.Log("Setting up mock query response for:", tt.name)
				epQuery.CleanMockResponse()
				epQuery.SetMockResponse(tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)

			err := eC.EnableCloudavenueServices(t.Context(), tt.params)
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

	err = eC.EnableCloudavenueServices(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestDisableCloudavenueServices(t *testing.T) {
	tests := []struct {
		name   string
		params ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		mockGetNetworkServicesResponse       any
		mockGetNetworkServicesResponseStatus int

		expectedErr bool
	}{
		{
			name: "Disable network services with valid ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Disable network services with valid Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
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
			name: "Failed to get network services",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockGetNetworkServicesResponse:       nil,
			mockGetNetworkServicesResponseStatus: http.StatusNotFound,
			expectedErr:                          true,
		},

		{
			name: "Error 500",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Error 401",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusUnauthorized,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.DisableCloudavenueServices()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			epGetNetworkServices := endpoints.GetEdgeGatewayServices()
			// Set up mock response for GetNetworkServices
			if tt.mockGetNetworkServicesResponse != nil || tt.mockGetNetworkServicesResponseStatus != 0 {
				t.Log("Setting up mock GetNetworkServices response for:", tt.name)
				epGetNetworkServices.CleanMockResponse()
				epGetNetworkServices.SetMockResponse(tt.mockGetNetworkServicesResponse, &tt.mockGetNetworkServicesResponseStatus)
			}

			eC := newClient(t)

			err := eC.DisableCloudavenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", tt.params)
			}

			ep.CleanMockResponse()
		})
	}
}

func TestGetCloudavenueServices(t *testing.T) {
	tests := []struct {
		name   string
		params ParamsEdgeGateway

		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name: "Get Cloud Avenue services with valid ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Get Cloud Avenue services with valid Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Get Cloud Avenue services with empty params",
			params: ParamsEdgeGateway{
				ID: "",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Error 401",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponseStatus: http.StatusUnauthorized,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.GetEdgeGatewayServices()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

			result, err := eC.GetCloudavenueServices(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil")
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}
