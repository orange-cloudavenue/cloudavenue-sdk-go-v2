/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/orange-cloudavenue/common-go/validators"
)

func TestListDraasOnPremiseIP(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name:        "ListDraasOnPremiseIP OK",
			expectedErr: false,
		},
		{
			name:               "ListDraasOnPremiseIP Not Found",
			mockResponseStatus: http.StatusNotFound,
			mockResponse:       nil,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.ListDraasOnPremiseIp().CleanMockResponse()
				endpoints.ListDraasOnPremiseIp().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			resp, err := client.ListOnPremiseIp(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Response should not be nil")
			assert.NotEmpty(t, resp.IPs, "OnPremise IPs should not be empty")
			for _, ip := range resp.IPs {
				assert.NoError(t, validators.New().Var(ip, "required,ipv4"))
			}
		})
	}
}

func TestAddDraasOnPremiseIP(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsAddDraasOnPremiseIP
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "AddDraasOnPremiseIP OK",
			params: types.ParamsAddDraasOnPremiseIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			expectedErr: false,
		},
		{
			name:        "No IP params provided",
			expectedErr: true,
		},
		{
			name: "Invalid IP format",
			params: types.ParamsAddDraasOnPremiseIP{
				IP: "invalid-ip",
			},
			expectedErr: true,
		},
		{
			name: "Bad request",
			params: types.ParamsAddDraasOnPremiseIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			mockResponseStatus: http.StatusBadRequest,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.AddDraasOnPremiseIp().CleanMockResponse()
				endpoints.AddDraasOnPremiseIp().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			err := client.AddOnPremiseIp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestRemoveDraasOnPremiseIP(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsRemoveDraasOnPremiseIP
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
	}{
		{
			name: "RemoveDraasOnPremiseIP OK",
			params: types.ParamsRemoveDraasOnPremiseIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			expectedErr: false,
		},
		{
			name:        "No IP params provided",
			expectedErr: true,
		},
		{
			name: "Invalid IP format",
			params: types.ParamsRemoveDraasOnPremiseIP{
				IP: "invalid-ip",
			},
			expectedErr: true,
		},
		{
			name: "Bad request",
			params: types.ParamsRemoveDraasOnPremiseIP{
				IP: generator.MustGenerate("{ipv4address}"),
			},
			mockResponseStatus: http.StatusBadRequest,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponseStatus != 0 {
				endpoints.RemoveDraasOnPremiseIp().CleanMockResponse()
				endpoints.RemoveDraasOnPremiseIp().SetMockResponse(tt.mockResponse, &tt.mockResponseStatus)
			}

			client := newClient(t)

			err := client.RemoveOnPremiseIp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err, "Unexpected error: %v", err)
		})
	}
}
