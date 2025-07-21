package edgegateway

import (
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/common-go/generator"
)

func Test_ListT0(t *testing.T) {
	tests := []struct {
		name                string
		queryResponse       any
		queryResponseStatus int
		expectedErr         bool
	}{
		{
			name:        "List T0",
			expectedErr: false,
		},
		{
			name:                "Error 500",
			queryResponseStatus: http.StatusInternalServerError,
			expectedErr:         false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name:                "Error 404",
			queryResponseStatus: http.StatusNotFound,
			expectedErr:         true, // Error HTTP 404 should return an error.
		},
		{
			name: "Simulate unknown class of service",
			queryResponse: &apiResponseT0s{
				{
					Type:       "edge-gateway",
					Name:       generator.MustGenerate("{edgegateway_name}"),
					Properties: apiResponseT0Properties{ClassOfService: "unknown"},
				},
			},
			expectedErr:         false,
			queryResponseStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep, _ := mock.GetEndpoint("T0", cav.MethodGET)
			// Set up mock response
			if tt.queryResponse != nil || tt.queryResponseStatus != 0 {
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.queryResponse, &tt.queryResponseStatus)
			}

			mC, err := mock.NewClient(mock.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, nil))))
			assert.Nil(t, err, "Error creating mock client")

			t0C, err := New(mC)
			assert.Nil(t, err, "Error creating t0 client")

			t0s, err := t0C.ListT0(t.Context())
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
			} else {
				assert.Nil(t, err, "Unexpected error while listing T0s")
				assert.NotNil(t, t0s, "Expected non-nil T0s response")
			}
		})
	}
}

func Test_GetT0(t *testing.T) {
	tests := []struct {
		name                string
		params              ParamsGetT0
		queryResponse       any
		queryResponseStatus int
		expectedErr         bool
	}{
		{
			name: "Valid T0",
			params: ParamsGetT0{
				Name: generator.MustGenerate("{t0_name}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid TO name",
			params: ParamsGetT0{
				Name: "invalid_t0_name",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: ParamsGetT0{
				Name: generator.MustGenerate("{t0_name}"),
			},
			queryResponseStatus: http.StatusInternalServerError,
			expectedErr:         false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Simulate empty response",
			params: ParamsGetT0{
				Name: generator.MustGenerate("{t0_name}"),
			},
			queryResponse:       &apiResponseT0s{},
			queryResponseStatus: http.StatusOK,
			expectedErr:         true,
		},
		{
			name: "Error 404",
			params: ParamsGetT0{
				Name: generator.MustGenerate("{t0_name}"),
			},
			queryResponseStatus: http.StatusNotFound,
			expectedErr:         true, // Error HTTP 404 should return an error.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep, _ := mock.GetEndpoint("T0", cav.MethodGET)
			if tt.queryResponse != nil || tt.queryResponseStatus != 0 {
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.queryResponse, &tt.queryResponseStatus)
			}

			mC, err := mock.NewClient(mock.WithLogger(slog.New(slog.NewTextHandler(os.Stderr, nil))))
			assert.Nil(t, err, "Error creating mock client")

			t0C, err := New(mC)
			assert.Nil(t, err, "Error creating t0 client")

			t0, err := t0C.GetT0(t.Context(), tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, t0, "Expected nil T0 response")
			} else {
				assert.Nil(t, err, "Unexpected error while getting T0")
				assert.NotNil(t, t0, "Expected non-nil T0 response")
				assert.Equal(t, tt.params.Name, t0.Name, "Expected T0 name to match the requested name")
			}
		})
	}
}
