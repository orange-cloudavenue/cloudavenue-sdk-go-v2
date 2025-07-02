package cav

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRequest_Cerberus(t *testing.T) {
	client, err := newMockClient()
	assert.Nil(t, err, "Error creating mock client")

	endpointSessionCerberus, err := GetEndpoint(CategoryAuthentication, VersionV1, "CreateSessionVmware", MethodPOST)
	assert.Nil(t, err, "Error getting endpoint for CreateSessionVmware")
	defer endpointSessionCerberus.CleanMockResponse()

	tests := []struct {
		name           string
		expectedErr    bool
		expectedResp   any
		expectedStatus int
	}{
		{
			name:           "success",
			expectedErr:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad request",
			expectedErr:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal server error",
			expectedErr:    true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "not found",
			expectedErr: true,
			expectedResp: cerberusError{
				Code:    "err-0002",
				Reason:  "unknown-0001",
				Message: "The request you are trying to perform is not valid.",
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "unauthorized",
			expectedErr:    true,
			expectedStatus: http.StatusUnauthorized,
			expectedResp: cerberusError{
				Code:    "err-0001",
				Reason:  "auth-0001",
				Message: "The request you are trying to make does not have sufficient permissions.",
			},
		},
		{
			name:           "unknown error",
			expectedErr:    true,
			expectedStatus: http.StatusBadGateway,
			expectedResp: struct {
				Foo string `json:"foo"`
			}{
				Foo: "Unknown error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer endpointSessionCerberus.CleanMockResponse()

			endpointSessionCerberus.SetMockResponse(tt.expectedResp, &tt.expectedStatus)

			req, err := client.NewRequest(t.Context(), ClientCerberus)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error creating request")
			assert.NotNil(t, req, "Expected request to be created")

			endpointSessionCerberus.CleanMockResponse()
		})
	}
}
