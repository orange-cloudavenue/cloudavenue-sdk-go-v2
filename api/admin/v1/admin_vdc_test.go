/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package admin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	mockVDCName = "vdc1"
	mockVDCID   = "urn:vcloud:vdc:11111111-1111-1111-1111-111111111111"
	mockVDCHref = "https://example.com/api/admin/vdc/" + mockVDCID
)

func TestListAdminVDCs(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
		expectedLen        int
	}{
		{
			name:        "List Admin VDCs Success",
			expectedErr: false,
			expectedLen: 1,
		},
		{
			name:               "List Admin VDCs Error 500",
			mockResponseStatus: 500,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListAdminVDCs())
				ms.SetResponseFunc(endpoints.ListAdminVDCs(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.ListAdminVDCs())
				ms.SetResponseFunc(endpoints.ListAdminVDCs(), func(w http.ResponseWriter, r *http.Request) {
					vdcs := itypes.AdminVDCs{
						VDCs: []itypes.AdminVDC{
							{
								Name:        mockVDCName,
								ID:          mockVDCID,
								Href:        mockVDCHref,
								DisplayName: "VDC One",
								IsEnabled:   true,
							},
						},
					}
					xmlResponse(w, vdcs)
				})
			}

			result, err := client.ListAdminVDCs(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, mockVDCName, result[0].Name)
			assert.Equal(t, mockVDCID, result[0].ID)
			assert.Equal(t, mockVDCHref, result[0].Href)
			assert.Equal(t, "VDC One", result[0].DisplayName)
			assert.True(t, result[0].IsEnabled)
		})
	}
}

func TestListAdminVDCs_Empty(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminVDCs())
	ms.SetResponseFunc(endpoints.ListAdminVDCs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminVDCs{})
	})

	result, err := client.ListAdminVDCs(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetAdminVDCByID(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAdminVDC())
	ms.SetResponseFunc(endpoints.GetAdminVDC(), func(w http.ResponseWriter, r *http.Request) {
		vdc := itypes.AdminVDC{
			Name:              mockVDCName,
			ID:                mockVDCID,
			Href:              mockVDCHref,
			DisplayName:       "VDC One",
			Description:       "A VDC for testing",
			IsEnabled:         true,
			IsFullProtected:   false,
			VMCount:           10,
			VMRunningCount:    5,
			VappCount:         3,
			VappTemplateCount: 2,
		}
		xmlResponse(w, vdc)
	})

	result, err := client.GetAdminVDC(t.Context(), types.ParamsGetAdminVDC{
		ID: mockVDCID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockVDCID, result.ID)
	assert.Equal(t, mockVDCName, result.Name)
	assert.Equal(t, mockVDCHref, result.Href)
	assert.Equal(t, "VDC One", result.DisplayName)
	assert.Equal(t, "A VDC for testing", result.Description)
	assert.True(t, result.IsEnabled)
	assert.False(t, result.IsFullProtected)
	assert.Equal(t, 10, result.VMCount)
	assert.Equal(t, 5, result.VMRunningCount)
	assert.Equal(t, 3, result.VappCount)
	assert.Equal(t, 2, result.VappTemplateCount)
}

func TestGetAdminVDCByName(t *testing.T) {
	client, ms := newClient(t)

	// First call: list VDCs to resolve by name.
	ms.CleanResponse(endpoints.ListAdminVDCs())
	ms.SetResponseFunc(endpoints.ListAdminVDCs(), func(w http.ResponseWriter, r *http.Request) {
		vdcs := itypes.AdminVDCs{
			VDCs: []itypes.AdminVDC{
				{
					Name:        mockVDCName,
					ID:          mockVDCID,
					Href:        mockVDCHref,
					DisplayName: "VDC One",
					IsEnabled:   true,
				},
			},
		}
		xmlResponse(w, vdcs)
	})

	// Second call: get the resolved VDC by ID.
	ms.CleanResponse(endpoints.GetAdminVDC())
	ms.SetResponseFunc(endpoints.GetAdminVDC(), func(w http.ResponseWriter, r *http.Request) {
		vdc := itypes.AdminVDC{
			Name:        mockVDCName,
			ID:          mockVDCID,
			Href:        mockVDCHref,
			DisplayName: "VDC One",
			IsEnabled:   true,
		}
		xmlResponse(w, vdc)
	})

	result, err := client.GetAdminVDC(t.Context(), types.ParamsGetAdminVDC{
		Name: mockVDCName,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockVDCName, result.Name)
	assert.Equal(t, mockVDCID, result.ID)
	assert.Equal(t, mockVDCHref, result.Href)
	assert.Equal(t, "VDC One", result.DisplayName)
	assert.True(t, result.IsEnabled)
}

func TestGetAdminVDC_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetAdminVDC())

	_, err := client.GetAdminVDC(t.Context(), types.ParamsGetAdminVDC{})
	assert.Error(t, err)
}

func TestGetAdminVDC_NotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminVDCs())
	ms.SetResponseFunc(endpoints.ListAdminVDCs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminVDCs{})
	})

	_, err := client.GetAdminVDC(t.Context(), types.ParamsGetAdminVDC{
		Name: mockVDCName,
	})
	assert.Error(t, err)
}
