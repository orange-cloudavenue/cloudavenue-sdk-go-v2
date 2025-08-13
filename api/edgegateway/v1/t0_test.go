/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
)

func Test_ListT0(t *testing.T) {
	tests := []struct {
		name string

		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name:        "List T0",
			expectedErr: false,
		},
		{
			name:               "Error 500",
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name:               "Error 404",
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true, // Error HTTP 404 should return an error.
		},
		{
			name: "Simulate unknown class of service",
			mockResponse: &itypes.ApiResponseT0s{
				{
					Type:       "edge-gateway",
					Name:       generator.MustGenerate("{resource_name:edgegateway}"),
					Properties: itypes.ApiResponseT0Properties{ClassOfService: "unknown"},
				},
			},
			expectedErr:        false,
			mockResponseStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.ListT0()
			// Set up mock response
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

			t0s, err := eC.ListT0(t.Context())
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
		name   string
		params types.ParamsGetT0

		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name: "Valid T0",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			expectedErr: false,
		},
		{
			name: "Invalid TO name",
			params: types.ParamsGetT0{
				T0Name: "invalid_t0_name",
			},
			expectedErr: true,
		},
		{
			name: "Error 500",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponseStatus: http.StatusInternalServerError,
			expectedErr:        false, // Error HTTP 500 does not return an error because a retry is performed.
		},
		{
			name: "Simulate empty response",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponse:       &itypes.ApiResponseT0s{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "Simulate empty response EdgeGateway Name",
			params: types.ParamsGetT0{
				EdgegatewayName: generator.MustGenerate("{resource_name:edgegateway}"),
			},
			mockResponse:       &itypes.ApiResponseT0s{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "Simulate empty response EdgeGateway ID",
			params: types.ParamsGetT0{
				EdgegatewayID: generator.MustGenerate("{urn:edgegateway}"),
			},
			mockResponse:       &itypes.ApiResponseT0s{},
			mockResponseStatus: http.StatusOK,
			expectedErr:        true,
		},
		{
			name: "Error 404",
			params: types.ParamsGetT0{
				T0Name: generator.MustGenerate("{resource_name:t0}"),
			},
			mockResponseStatus: http.StatusNotFound,
			expectedErr:        true, // Error HTTP 404 should return an error.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := endpoints.ListT0()
			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ep.CleanMockResponse()
				ep.SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			eC := newClient(t)

			t0, err := eC.GetT0(t.Context(), tt.params)

			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got nil")
				assert.Nil(t, t0, "Expected nil T0 response")
			} else {
				assert.Nil(t, err, "Unexpected error while getting T0")
				assert.NotNil(t, t0, "Expected non-nil T0 response")
				if tt.params.T0Name != "" {
					assert.Equal(t, tt.params.T0Name, t0.Name, "Expected T0 name to match the requested name")
				}
				if tt.params.EdgegatewayID != "" || tt.params.EdgegatewayName != "" {
					assert.NotEmpty(t, t0.EdgeGateways, "Expected T0 to have edge gateways")
				}
			}
		})
	}
}
