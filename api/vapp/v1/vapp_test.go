/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListVApp(t *testing.T) {
	tests := []struct {
		name               string
		vdcID              string
		mockResponse       any
		mockResponseStatus int

		expectedErr bool
	}{
		{
			name:  "List VApps with valid VDC ID",
			vdcID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			mockResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{
					{
						ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
						Name: "test-vapp-1",
					},
				},
			},
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name:        "List VApps with empty VDC ID",
			vdcID:       "",
			expectedErr: true,
		},
		{
			name:  "Error 401 Unauthorized",
			vdcID: "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			mockResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{},
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			resp, err := client.ListVApp(t.Context(), tt.vdcID)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestGetVApp(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsGetVApp
		mockResponse       any
		mockResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Get VApp by ID",
			params: types.ParamsGetVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponse: &itypes.APIResponseGetVApp{
				ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
				Name: "test-vapp-1",
			},
			mockResponseStatus: 200,
			expectedErr:        false,
		},
		{
			name: "Get VApp by Name",
			params: types.ParamsGetVApp{
				Name: "test-vapp-1",
			},
			mockListResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{
					{
						ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
						Name: "test-vapp-1",
					},
				},
			},
			mockListResponseStatus: 200,
			expectedErr:            false,
		},
		{
			name: "Error 401 Unauthorized",
			params: types.ParamsGetVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
		{
			name: "VApp not found by name",
			params: types.ParamsGetVApp{
				Name: "test-vapp-1",
			},
			mockListResponse:       &itypes.APIResponseListVApp{Records: []itypes.APIResponseListVAppRecord{}},
			mockListResponseStatus: 200,
			expectedErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetVApp())
				ms.SetResponse(endpoints.GetVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockListResponse, &tt.mockListResponseStatus)
			}

			resp, err := client.GetVApp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestCreateVApp(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsCreateVApp
		mockResponse       any
		mockResponseStatus int

		mockGetResponse       any
		mockGetResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Create VApp with valid parameters",
			params: types.ParamsCreateVApp{
				Name:        "test-vapp-1",
				Description: "Test VApp",
				VDCID:       "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponse: &cav.Job{
				ID:     "87ab1934-0146-4fb0-80bc-815fea03214d",
				Status: "queued",
			},
			mockResponseStatus: 200,
			mockListResponse: &itypes.APIResponseListVApp{
				Records: []itypes.APIResponseListVAppRecord{
					{
						ID:   "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
						Name: "test-vapp-1",
					},
				},
			},
			mockListResponseStatus: 200,
			expectedErr:            false,
		},
		{
			name: "Create VApp with missing required parameters",
			params: types.ParamsCreateVApp{
				Name: "test-vapp-1",
			},
			expectedErr: true,
		},
		{
			name: "Error 401 Unauthorized",
			params: types.ParamsCreateVApp{
				Name:        "test-vapp-1",
				Description: "Test VApp",
				VDCID:       "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponseStatus: 401,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.CreateVApp())
				ms.SetResponse(endpoints.CreateVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockGetResponse != nil || tt.mockGetResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetVApp())
				ms.SetResponse(endpoints.GetVApp(), tt.mockGetResponse, &tt.mockGetResponseStatus)
			}

			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockListResponse, &tt.mockListResponseStatus)
			}

			resp, err := client.CreateVApp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
			assert.NotNil(t, resp, "Expected response to be not nil")
		})
	}
}

func TestUpdateVAppPreservesExistingDescription(t *testing.T) {
	client, ms := newClient(t)
	vappID := "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"
	deploymentLease := 3600
	storageLease := 7200
	getCalls := 0

	ms.CleanResponse(endpoints.GetVApp())
	ms.SetResponseFunc(endpoints.GetVApp(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &itypes.APIResponseGetVApp{
			ID:          vappID,
			Name:        "test-vapp-1",
			Description: "current description",
			LeaseSettings: itypes.APIResponseLeaseSettings{
				DeploymentLeaseInSeconds: deploymentLease,
				StorageLeaseInSeconds:    storageLease,
			},
		}
		if getCalls > 0 {
			resp.LeaseSettings.DeploymentLeaseInSeconds = 1800
		}
		getCalls++
		assert.NoError(t, json.NewEncoder(w).Encode(resp))
	})

	ms.CleanResponse(endpoints.UpdateVApp())
	ms.SetResponseFunc(endpoints.UpdateVApp(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.APIRequestUpdateVApp
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if assert.NotNil(t, body.Description) {
			assert.Equal(t, "current description", *body.Description)
		}
		if assert.NotNil(t, body.LeaseSettings) && assert.NotNil(t, body.LeaseSettings.DeploymentLeaseInSeconds) {
			assert.Equal(t, 1800, *body.LeaseSettings.DeploymentLeaseInSeconds)
		}
		if assert.NotNil(t, body.LeaseSettings) && assert.NotNil(t, body.LeaseSettings.StorageLeaseInSeconds) {
			assert.Equal(t, storageLease, *body.LeaseSettings.StorageLeaseInSeconds)
		}
		w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.UpdateVApp(t.Context(), types.ParamsUpdateVApp{ID: vappID, DeploymentLeaseInSeconds: new(1800)})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "current description", resp.Description)
	assert.NotNil(t, resp.DeploymentLeaseInSeconds)
	assert.Equal(t, 1800, *resp.DeploymentLeaseInSeconds)

	ms.CleanResponse(endpoints.GetVApp())
	ms.CleanResponse(endpoints.UpdateVApp())
}

func TestUpdateVAppCanClearDescription(t *testing.T) {
	client, ms := newClient(t)
	vappID := "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012"
	empty := ""
	getCalls := 0

	ms.CleanResponse(endpoints.GetVApp())
	ms.SetResponseFunc(endpoints.GetVApp(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &itypes.APIResponseGetVApp{
			ID:          vappID,
			Name:        "test-vapp-1",
			Description: "current description",
		}
		if getCalls > 0 {
			resp.Description = empty
		}
		getCalls++
		assert.NoError(t, json.NewEncoder(w).Encode(resp))
	})

	ms.CleanResponse(endpoints.UpdateVApp())
	ms.SetResponseFunc(endpoints.UpdateVApp(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.APIRequestUpdateVApp
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		if assert.NotNil(t, body.Description) {
			assert.Equal(t, empty, *body.Description)
		}
		w.Header().Set("Location", "/api/task/87ab1934-0146-4fb0-80bc-815fea03214d")
		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.UpdateVApp(t.Context(), types.ParamsUpdateVApp{ID: vappID, Description: &empty})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, empty, resp.Description)

	ms.CleanResponse(endpoints.GetVApp())
	ms.CleanResponse(endpoints.UpdateVApp())
}

func TestDeleteVApp(t *testing.T) {
	tests := []struct {
		name               string
		params             types.ParamsDeleteVApp
		mockResponse       any
		mockResponseStatus int

		mockGetResponseStatus int

		mockListResponse       any
		mockListResponseStatus int

		expectedErr bool
	}{
		{
			name: "Delete VApp with valid ID",
			params: types.ParamsDeleteVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			expectedErr: false,
		},
		{
			name: "Error 404 Not Found",
			params: types.ParamsDeleteVApp{
				ID: "urn:vcloud:vapp:12345678-1234-4b8d-89ab-123456789012",
			},
			mockResponseStatus: 404,
			expectedErr:        true,
		},
		{
			name: "Error 404 on Get VApp",
			params: types.ParamsDeleteVApp{
				Name: "test-vapp-1",
			},
			mockListResponse:       &itypes.APIResponseListVApp{Records: []itypes.APIResponseListVAppRecord{}},
			mockListResponseStatus: 404,
			expectedErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.DeleteVApp())
				ms.SetResponse(endpoints.DeleteVApp(), tt.mockResponse, &tt.mockResponseStatus)
			}

			if tt.mockGetResponseStatus != 0 {
				ms.CleanResponse(endpoints.GetVApp())
				ms.SetResponse(endpoints.GetVApp(), nil, &tt.mockGetResponseStatus)
			}

			if tt.mockListResponse != nil || tt.mockListResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListVApp())
				ms.SetResponse(endpoints.ListVApp(), tt.mockListResponse, &tt.mockListResponseStatus)
			}

			err := client.DeleteVApp(t.Context(), tt.params)
			if tt.expectedErr {
				assert.NotNil(t, err, "Expected error but got none")
				return
			}
			assert.Nil(t, err, "Unexpected error: %v", err)
		})
	}
}

func TestIsVAppOrgNetwork(t *testing.T) {
	assert.True(t, IsVAppOrgNetwork("https://example.com/api/network/1234"))
	assert.True(t, IsVAppOrgNetwork("/api/network/1234"))
	assert.False(t, IsVAppOrgNetwork("https://example.com/api/vAppNetwork/vappNetwork-1234"))
}

func TestIsVAppNetwork(t *testing.T) {
	assert.True(t, IsVAppNetwork("https://example.com/api/vAppNetwork/vappNetwork-1234"))
	assert.True(t, IsVAppNetwork("/api/vAppNetwork/vappNetwork-1234"))
	assert.False(t, IsVAppNetwork("https://example.com/api/network/1234"))
}

func Test_NewClient_ClientNil(t *testing.T) {
	c, err := New(nil)
	assert.Nil(t, c, "Expected nil client when input is nil")
	assert.Error(t, err, "Expected error when input is nil")
}
