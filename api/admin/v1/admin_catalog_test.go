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
	mockCatalogID   = "urn:vcloud:catalog:11111111-1111-1111-1111-111111111111"
	mockCatalogName = "cat1"
	mockCatalogHref = "https://example.com/api/admin/catalog/" + mockCatalogID
)

func TestListAdminCatalogs(t *testing.T) {
	tests := []struct {
		name               string
		mockResponseStatus int
		mockResponse       any
		expectedErr        bool
		expectedLen        int
	}{
		{name: "List Admin Catalogs Success", expectedErr: false, expectedLen: 1},
		{name: "List Admin Catalogs Error 500", mockResponseStatus: 500, expectedErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			if tt.mockResponse != nil || tt.mockResponseStatus != 0 {
				ms.CleanResponse(endpoints.ListAdminCatalogs())
				ms.SetResponseFunc(endpoints.ListAdminCatalogs(), func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.mockResponseStatus)
					if tt.mockResponse != nil {
						xmlResponse(w, tt.mockResponse)
					}
				})
			} else {
				ms.CleanResponse(endpoints.ListAdminCatalogs())
				ms.SetResponseFunc(endpoints.ListAdminCatalogs(), func(w http.ResponseWriter, r *http.Request) {
					xmlResponse(w, itypes.AdminCatalogs{
						Catalogs: []itypes.AdminCatalog{
							{Name: mockCatalogName, ID: mockCatalogID, Href: mockCatalogHref, IsEnabled: true},
						},
					})
				})
			}

			result, err := client.ListAdminCatalogs(t.Context())
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, tt.expectedLen)
			assert.Equal(t, mockCatalogName, result[0].Name)
			assert.Equal(t, mockCatalogID, result[0].ID)
			assert.Equal(t, mockCatalogHref, result[0].Href)
			assert.True(t, result[0].IsEnabled)
		})
	}
}

func TestListAdminCatalogs_Empty(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminCatalogs())
	ms.SetResponseFunc(endpoints.ListAdminCatalogs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminCatalogs{})
	})

	result, err := client.ListAdminCatalogs(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestGetAdminCatalogByID(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetAdminCatalog())
	ms.SetResponseFunc(endpoints.GetAdminCatalog(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminCatalog{
			Name: mockCatalogName, ID: mockCatalogID, Href: mockCatalogHref,
			Description: "A catalog for testing",
			IsEnabled:   true, IsPublished: true, IsTrusted: true,
		})
	})

	result, err := client.GetAdminCatalog(t.Context(), types.ParamsGetAdminCatalog{ID: mockCatalogID})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockCatalogID, result.ID)
	assert.Equal(t, mockCatalogName, result.Name)
	assert.Equal(t, mockCatalogHref, result.Href)
	assert.Equal(t, "A catalog for testing", result.Description)
	assert.True(t, result.IsEnabled)
	assert.True(t, result.IsPublished)
	assert.True(t, result.IsTrusted)
}

func TestGetAdminCatalogByName(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminCatalogs())
	ms.SetResponseFunc(endpoints.ListAdminCatalogs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminCatalogs{
			Catalogs: []itypes.AdminCatalog{
				{Name: mockCatalogName, ID: mockCatalogID, Href: mockCatalogHref, IsEnabled: true},
			},
		})
	})

	ms.CleanResponse(endpoints.GetAdminCatalog())
	ms.SetResponseFunc(endpoints.GetAdminCatalog(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminCatalog{
			Name: mockCatalogName, ID: mockCatalogID, Href: mockCatalogHref,
			Description: "A catalog for testing", IsEnabled: true,
		})
	})

	result, err := client.GetAdminCatalog(t.Context(), types.ParamsGetAdminCatalog{Name: mockCatalogName})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockCatalogName, result.Name)
	assert.Equal(t, mockCatalogID, result.ID)
	assert.Equal(t, mockCatalogHref, result.Href)
	assert.Equal(t, "A catalog for testing", result.Description)
	assert.True(t, result.IsEnabled)
}

func TestGetAdminCatalog_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetAdminCatalog())

	_, err := client.GetAdminCatalog(t.Context(), types.ParamsGetAdminCatalog{})
	assert.Error(t, err)
}

func TestGetAdminCatalog_NotFound(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListAdminCatalogs())
	ms.SetResponseFunc(endpoints.ListAdminCatalogs(), func(w http.ResponseWriter, r *http.Request) {
		xmlResponse(w, itypes.AdminCatalogs{})
	})

	_, err := client.GetAdminCatalog(t.Context(), types.ParamsGetAdminCatalog{Name: mockCatalogName})
	assert.Error(t, err)
}
