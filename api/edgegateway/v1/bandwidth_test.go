package edgegateway

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/common-go/generator"
)

// func TestGetBandwidth(t *testing.T) {
// 	tests := []struct {
// 		name                string
// 		params              ParamsEdgeGateway
// 		queryResponse       any
// 		queryResponseStatus int
// 		expectedErr         bool
// 	}{
// 		{
// 			name: "Valid Edge Gateway ID",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "Valid Edge Gateway Name",
// 			params: ParamsEdgeGateway{
// 				Name: generator.MustGenerate("{resource_name:edgegateway}"),
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "Invalid Edge Gateway ID",
// 			params: ParamsEdgeGateway{
// 				ID: "invalid-id",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Error 500",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			queryResponse:       struct{}{},
// 			queryResponseStatus: 500,
// 			expectedErr:         false, // Error HTTP 500 does not return an error because a retry is performed.
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ep := endpoints.GetEdgeGatewayBandwidth()
// 			if tt.queryResponse != nil || tt.queryResponseStatus != 0 {
// 				mock.SetMockResponse(ep, tt.queryResponse, &tt.queryResponseStatus)
// 			}

// 			mC, err := mock.NewClient()
// 			assert.Nil(t, err, "Error creating mock client")

// 			eC, err := New(mC)
// 			assert.Nil(t, err, "Error creating edgegateway client")

// 			result, err := eC.GetBandwidth(t.Context(), tt.params)
// 			if tt.expectedErr {
// 				assert.NotNil(t, err, "Expected error but got nil")
// 				assert.Nil(t, result, "Result should be nil when error is expected")
// 			} else {
// 				assert.Nil(t, err, "Unexpected error: %v", tt.params)
// 				assert.NotNil(t, result, "Result should not be nil")
// 			}
// 		})
// 	}
// }

// func TestGetBandwidth_ContextDeadlineExceeded(t *testing.T) {
// 	mC, err := mock.NewClient()
// 	assert.Nil(t, err, "Error creating mock client")

// 	eC, err := New(mC)
// 	assert.Nil(t, err, "Error creating edgegateway client")

// 	// Simulate a context deadline exceeded error
// 	ctx, cancel := context.WithTimeout(t.Context(), 0)
// 	defer cancel()

// 	_, err = eC.GetBandwidth(ctx, ParamsEdgeGateway{ID: generator.MustGenerate("{urn:edgegateway}")})
// 	assert.NotNil(t, err, "Expected context deadline exceeded error")
// 	assert.Contains(t, err.Error(), "context deadline exceeded", "Expected error to contain 'context deadline exceeded'")
// }

// func TestGetRemainingBandwidth(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		params      ParamsEdgeGateway
// 		queryResp   any
// 		queryStatus int
// 		expectedErr bool
// 	}{
// 		{
// 			name: "Valid Edge Gateway Name",
// 			params: ParamsEdgeGateway{
// 				Name: generator.MustGenerate("{resource_name:edgegateway}"),
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "Valid Edge Gateway ID",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			expectedErr: false,
// 		},
// 		{
// 			name: "Invalid Edge Gateway Name",
// 			params: ParamsEdgeGateway{
// 				Name: "invalid-edgegateway-name",
// 			},
// 			expectedErr: true,
// 		},
// 		{
// 			name:        "No parameters provided",
// 			params:      ParamsEdgeGateway{},
// 			expectedErr: true,
// 		},
// 		{
// 			name: "EdgeGateway not found by name",
// 			params: ParamsEdgeGateway{
// 				Name: generator.MustGenerate("{resource_name:edgegateway}"),
// 			},
// 			queryResp: &apiResponseT0s{
// 				apiResponseT0{
// 					Type: "tier-0-vrf",
// 					Name: generator.MustGenerate("{resource_name:t0}"),
// 					Children: []apiResponseT0Children{
// 						{
// 							Type: "edge-gateway",
// 							Name: generator.MustGenerate("{resource_name:edgegateway}"),
// 						},
// 					},
// 				},
// 			},
// 			queryStatus: http.StatusOK,
// 			expectedErr: true,
// 		},
// 		{
// 			name: "EdgeGateway not found by ID",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			queryResp: &apiResponseT0s{
// 				apiResponseT0{
// 					Type: "tier-0-vrf",
// 					Name: generator.MustGenerate("{resource_name:t0}"),
// 					Children: []apiResponseT0Children{
// 						{
// 							Type: "edge-gateway",
// 							Name: generator.MustGenerate("{resource_name:edgegateway}"),
// 							Properties: struct {
// 								RateLimit int    "json:\"rateLimit,omitempty\" fake:\"5\""
// 								EdgeUUID  string "json:\"edgeUuid,omitempty\" fake:\"{urn:edgegateway}\""
// 							}{
// 								RateLimit: 5,
// 								EdgeUUID:  must(extractor.ExtractUUID(generator.MustGenerate("{urn:edgegateway}"))),
// 							},
// 						},
// 					},
// 				},
// 			},
// 			queryStatus: http.StatusOK,
// 			expectedErr: true,
// 		},
// 		{
// 			name: "Error 500",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			queryStatus: 500,
// 			expectedErr: false, // Error HTTP 500 does not return an error because a retry is performed.
// 		},
// 		{
// 			name: "Error 404",
// 			params: ParamsEdgeGateway{
// 				ID: generator.MustGenerate("{urn:edgegateway}"),
// 			},
// 			queryStatus: 404,
// 			expectedErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ep := endpoints.GetEdgeGatewayBandwidth()

// 			eC := newClient(t)

// 			if tt.queryResp != nil || tt.queryStatus != 0 {
// 				// If we expect a valid response, we need to set the mock response
// 				mock.SetMockResponse(ep, tt.queryResp, &tt.queryStatus)
// 			}

// 			bandwidth, err := eC.GetEdgeGatewayBandwidth(t.Context(), tt.params)
// 			if tt.expectedErr {
// 				assert.NotNil(t, err, "Expected error but got nil")
// 				assert.Empty(t, bandwidth, "Expected empty bandwidth but got %v", bandwidth)
// 			} else {
// 				assert.Nil(t, err, "Expected no error but got %v", err)
// 			}
// 		})
// 	}
// }

func Test_GetEdgeGatewayBandwidth(t *testing.T) {
	tests := []struct {
		name               string
		params             ParamsEdgeGateway
		mockResponse       any
		mockResponseStatus int
		expectedErr        bool
	}{
		{
			name: "Valid Edge Gateway ID",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Valid Edge Gateway Name",
			params: ParamsEdgeGateway{
				Name: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid Edge Gateway ID",
			params: ParamsEdgeGateway{
				ID: "invalid-id",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 500,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Error 404",
			params: ParamsEdgeGateway{
				ID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       struct{}{},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eC := newClient(t)

			// Set up mock response
			ep := endpoints.ListT0()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				t.Log("Setting up mock response for:", tt.name)
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			result, err := eC.GetBandwidth(t.Context(), tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, result, "Result should be nil when error is expected")
			} else {
				assert.Nil(t, err, "Unexpected error: %v", err)
				assert.NotNil(t, result, "Result should not be nil")
				assert.NotEmpty(t, result.ID, "Expected edge gateway ID to match")
				assert.NotEmpty(t, result.Name, "Expected edge gateway name to match")
			}
		})
	}
}
