package edgegateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/common-go/generator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func TestGetEdgeGateway(t *testing.T) {
	tests := []struct {
		name                    string
		params                  *ParamsEdgeGateway
		mockQueryResponse       any
		mockQueryResponseStatus int
		mockResponse            any
		mockResponseStatus      int
		expectedErr             bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			expectedErr: false,
			mockResponse: &ModelEdgeGateway{
				ID:   generator.MustGenerate("{urn:edgeGateway}"),
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			mockResponseStatus: 200,
		},
		{
			name: "Failed to retrieve Edge Gateway ID by name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			mockQueryResponseStatus: 404,
			expectedErr:             true,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: &ParamsEdgeGateway{
				ID: "urn:vcloud:vm:invalid-id",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
			mockResponse:       nil,
			mockResponseStatus: 500,
		},
		{
			name:        "Error validation params",
			params:      &ParamsEdgeGateway{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.GetEdgeGateway()
			epQuery := endpoints.QueryEdgeGateway()

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
				// If we expect a valid response, we need to set the mock response
				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockQueryResponse != nil || tt.mockQueryResponseStatus != 0 {
				// If we expect a query response, we need to set the mock response for the
				mock.SetMockResponse(epQuery, tt.mockQueryResponse, &tt.mockQueryResponseStatus)
			}

			eC := newClient(t)

			// Call the GetEdgeGateway method
			result, err := eC.GetEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error: %v", tt.params)
				assert.Nil(t, result, "Result should be nil: %v", tt.params)
			} else {
				assert.Nil(t, err, "Expected no error: %v", tt.params)
				assert.NotNil(t, result, "Result should not be nil: %v", tt.params)
			}
		})
	}
}

func TestGetEdgeGateway_ContextDeadlineExceeded(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	_, err = eC.GetEdgeGateway(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgeGateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

func TestRetrieveEdgeGatewayIDByName(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Mock the QueryEdgeGateway endpoint
	epQuery := endpoints.QueryEdgeGateway()

	tests := []struct {
		name        string
		edgeName    string
		queryResp   *apiResponseQueryEdgeGateway
		queryStatus int
		expectedID  string
		expectedErr bool
	}{
		{
			name:     "Valid Edge Gateway Name",
			edgeName: generator.MustGenerate("{edgegateway_name}"),
			queryResp: &apiResponseQueryEdgeGateway{
				Record: []apiResponseQueryEdgeGatewayRecord{
					{ID: "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c", HREF: "https://api.example.com/edgegateways/ed0a243a-374b-4306-ab25-9c3787cbdb4c"},
				},
			},
			queryStatus: 200,
			expectedID:  "urn:vcloud:gateway:ed0a243a-374b-4306-ab25-9c3787cbdb4c",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.SetMockResponse(epQuery, tt.queryResp, &tt.queryStatus)

			id, err := eC.retrieveEdgeGatewayIDByName(t.Context(), tt.edgeName)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Empty(t, id, "Expected empty ID but got %s", id)
			} else {
				assert.Nil(t, err, "Expected no error but got %v", err)
				assert.Equal(t, tt.expectedID, id, "Expected ID %s but got %s", tt.expectedID, id)
			}
		})
	}
}

func TestDeleteEdgeGateway(t *testing.T) {
	tests := []struct {
		name                        string
		params                      *ParamsEdgeGateway
		mockResponse                any
		mockResponseStatus          int
		mockMockQueryResponseStatus int
		expectedErr                 bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponse: nil,
			// mockResponseStatus: 202,
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			mockResponse: nil,
			// mockResponseStatus: 204,
			expectedErr: false,
		},
		{
			name: "Invalid Edge Gateway Name",
			params: &ParamsEdgeGateway{
				Name: "invalidEdgeGateway",
			},
			mockResponse:       nil,
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: &ParamsEdgeGateway{
				ID: "urn:vcloud:gateway:invalid-id",
			},
			mockResponse:       nil,
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Error 500",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponse:       nil,
			mockResponseStatus: 500,
			expectedErr:        false,
		},
		{
			name: "Error 401",
			params: &ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgeGateway}"),
			},
			mockResponse:       nil,
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "error 404 edge gateway name and id not found",
			params: &ParamsEdgeGateway{
				Name: generator.MustGenerate("{edgegateway_name}"),
			},
			mockResponse:                nil,
			mockResponseStatus:          404,
			mockMockQueryResponseStatus: 404,
			expectedErr:                 true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			epDelete := endpoints.DeleteEdgeGateway()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Logf("Setting mock response for endpoint %s with status %d", epDelete.Name, tt.mockResponseStatus)
				// If we expect a valid response, we need to set the mock response
				mock.SetMockResponse(epDelete, tt.mockResponse, &tt.mockResponseStatus)
			}

			epQuery := endpoints.QueryEdgeGateway()
			if tt.mockMockQueryResponseStatus != 0 {
				t.Logf("Setting mock response for query endpoint %s with status %d", epQuery.Name, tt.mockMockQueryResponseStatus)
				// If we expect a query response, we need to set the mock response for the
				mock.SetMockResponse(epQuery, nil, &tt.mockMockQueryResponseStatus)
			}

			eC := newClient(t)
			err := eC.DeleteEdgeGateway(t.Context(), *tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error for params: %v", tt.params)
			} else {
				assert.Nil(t, err, "Expected no error for params: %v", tt.params)
			}
		})
	}
}

func TestDeleteEdgeGateway_ContextDeadlineExceeded(t *testing.T) {
	mC, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating edgegateway client")

	// Simulate a context deadline exceeded error
	ctx, cancel := context.WithTimeout(t.Context(), 0)
	defer cancel()

	err = eC.DeleteEdgeGateway(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgeGateway}")})
	assert.NotNil(t, err, "Expected context deadline exceeded error")
	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
}

// func TestCreateEdgeGateway(t *testing.T) {
// 	tests := []struct {
// 		name               string
// 		params             *ParamsCreateEdgeGateway
// 		mockResponse       any
// 		mockResponseStatus int
// 		expectedErr        bool
// 	}{
// 		{
// 			name: "Valid Edge Gateway Creation",
// 			params: &ParamsCreateEdgeGateway{
// 				OwnerType: "vdc",
// 				OwnerName: generator.MustGenerate("{word}"),
// 			},
// 			// mockResponse: &ModelEdgeGateway{
// 			// 	ID:   generator.MustGenerate("{urn:edgeGateway}"),
// 			// 	Name: generator.MustGenerate("{edgegateway_name}"),
// 			// },
// 			// mockResponseStatus: 201,
// 			expectedErr: false,
// 		},
// 		// {
// 		// 	name: "Invalid Edge Gateway Name",
// 		// 	params: &ParamsCreateEdgeGateway{
// 		// 		Name: "",
// 		// 	},
// 		// 	mockResponseStatus: 400,
// 		// 	expectedErr: true,
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ep, _ := mock.GetEndpoint("EdgeGateway", cav.MethodPOST)

// 			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
// 				t.Logf("Setting mock response for endpoint %s with status %d", ep.Name, tt.mockResponseStatus)
// 				mock.SetMockResponse(ep, tt.mockResponse, &tt.mockResponseStatus)
// 			}

// 			eC := newClient(t)

// 			result, err := eC.CreateEdgeGateway(t.Context(), *tt.params)
// 			if tt.expectedErr {
// 				assert.NotNil(t, err, "Expected error for params: %v", tt.params)
// 				assert.Nil(t, result, "Result should be nil for params: %v", tt.params)
// 			} else {
// 				assert.Nil(t, err, "Expected no error for params: %v", tt.params)
// 				assert.NotNil(t, result, "Result should not be nil for params: %v", tt.params)
// 			}
// 		})
// 	}
// }
