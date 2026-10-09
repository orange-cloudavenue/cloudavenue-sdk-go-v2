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
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func newClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	mC, ms, err := mock.NewClient()
	assert.Nil(t, err, "Error creating mock client")

	eC, err := New(mC)
	assert.Nil(t, err, "Error creating admin client")
	return eC, ms
}

func xmlResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("X-Cloud-Avenue-Mock", "true")
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	mockOrgID   = "urn:vcloud:org:11111111-1111-1111-1111-111111111111"
	mockOrgName = "org1"
	mockOrgHref = "https://example.com/api/admin/org/" + mockOrgID
)

func TestListAdminOrgs(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
		expectedLen        int
	}{
		{
			name:        "List Admin Orgs Success",
			expectedErr: false,
			expectedLen: 1,
		},
		{
			name:               "List Admin Orgs Error 500",
			mockResponseStatus: 500,
			expectedErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListAdminOrgs())
				ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.ListAdminOrgs())
				ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, r *http.Request) {
					orgs := itypes.AdminOrgs{
						Orgs: []itypes.AdminOrg{
							{
								Name:        mockOrgName,
								ID:          mockOrgID,
								Href:        mockOrgHref,
								DisplayName: "Org One",
								IsEnabled:   true,
							},
						},
					}
					xmlResponse(w, orgs)
				})
			}

			result, err := client.ListAdminOrgs(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, mockOrgName, result[0].Name)
			assert.Equal(t, mockOrgID, result[0].ID)
			assert.Equal(t, mockOrgHref, result[0].Href)
			assert.Equal(t, "Org One", result[0].DisplayName)
			assert.True(t, result[0].IsEnabled)
		})
	}
}

func TestListAdminOrgs_Empty(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminOrgs())
	ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminOrgs{})
	})

	result, err := client.ListAdminOrgs(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetAdminOrgByID(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAdminOrg())
	ms.SetResponseFunc(endpoints.GetAdminOrg(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/org/"+mockOrgID, r.URL.Path)
		org := itypes.AdminOrg{
			Name:            mockOrgName,
			ID:              mockOrgID,
			Href:            mockOrgHref,
			DisplayName:     "Org One",
			Description:     "An org for testing",
			IsEnabled:       true,
			IsFullProtected: false,
		}
		xmlResponse(w, org)
	})

	result, err := client.GetAdminOrg(t.Context(), types.ParamsGetAdminOrg{
		ID: mockOrgID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockOrgID, result.ID)
	assert.Equal(t, mockOrgName, result.Name)
	assert.Equal(t, mockOrgHref, result.Href)
	assert.Equal(t, "Org One", result.DisplayName)
	assert.Equal(t, "An org for testing", result.Description)
	assert.True(t, result.IsEnabled)
	assert.False(t, result.IsFullProtected)
}

func TestGetAdminOrgByName(t *testing.T) {
	client, ms := newClient(t)

	// First call: list orgs to resolve by name.
	ms.CleanResponse(endpoints.ListAdminOrgs())
	ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/orgs", r.URL.Path)
		orgs := itypes.AdminOrgs{
			Orgs: []itypes.AdminOrg{
				{
					Name:        mockOrgName,
					ID:          mockOrgID,
					Href:        mockOrgHref,
					DisplayName: "Org One",
					IsEnabled:   true,
				},
			},
		}
		xmlResponse(w, orgs)
	})

	// Second call: get the resolved org by ID.
	ms.CleanResponse(endpoints.GetAdminOrg())
	ms.SetResponseFunc(endpoints.GetAdminOrg(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/org/"+mockOrgID, r.URL.Path)
		org := itypes.AdminOrg{
			Name:        mockOrgName,
			ID:          mockOrgID,
			Href:        mockOrgHref,
			DisplayName: "Org One",
			IsEnabled:   true,
		}
		xmlResponse(w, org)
	})

	result, err := client.GetAdminOrg(t.Context(), types.ParamsGetAdminOrg{
		Name: mockOrgName,
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockOrgName, result.Name)
	assert.Equal(t, mockOrgID, result.ID)
	assert.Equal(t, mockOrgHref, result.Href)
	assert.Equal(t, "Org One", result.DisplayName)
	assert.True(t, result.IsEnabled)
}

func TestGetAdminOrg_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetAdminOrg())

	_, err := client.GetAdminOrg(t.Context(), types.ParamsGetAdminOrg{})
	assert.Error(t, err)
}

func TestGetAdminOrgByName_NotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminOrgs())
	ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, _ *http.Request) {
		xmlResponse(w, itypes.AdminOrgs{Orgs: []itypes.AdminOrg{{Name: "other-org", ID: mockOrgID}}})
	})

	result, err := client.GetAdminOrg(t.Context(), types.ParamsGetAdminOrg{Name: mockOrgName})
	assert.Nil(t, result)
	require.EqualError(t, err, `Admin.GetOrg: organization with name "org1" not found`)
}

func TestGetAdminOrg_NotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminOrgs())
	ms.SetResponseFunc(endpoints.ListAdminOrgs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminOrgs{})
	})

	_, err := client.GetAdminOrg(t.Context(), types.ParamsGetAdminOrg{
		Name: mockOrgName,
	})
	assert.Error(t, err)
}
