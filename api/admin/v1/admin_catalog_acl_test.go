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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	mockACLUserURN     = "urn:vcloud:orguser:22222222-2222-2222-2222-222222222222"
	mockACLUserName    = "user1"
	mockACLAccessLevel = "Change"
)

func TestGetAdminCatalogACLByID(t *testing.T) {
	assert.Equal(t, "/api/catalog/{catalogId}/controlAccess", endpoints.GetAdminCatalogACL().PathTemplate)
	assert.Equal(t, "/api/catalog/{catalogId}/action/controlAccess", endpoints.SetAdminCatalogACL().PathTemplate)

	tests := []struct {
		name         string
		idOrName     string
		expectErr    bool
		expectShared bool
	}{
		{
			name:         "Get Admin Catalog ACL by ID Success",
			idOrName:     mockCatalogID,
			expectErr:    false,
			expectShared: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)

			ms.CleanResponse(endpoints.GetAdminCatalogACL())
			ms.SetResponseFunc(endpoints.GetAdminCatalogACL(), func(w http.ResponseWriter, r *http.Request) {
				xmlResponse(w, itypes.ControlAccessParams{
					Xmlns:               "http://www.vmware.com/vcloud/v1.5",
					IsSharedToEveryone:  true,
					EveryoneAccessLevel: stringPtr("ReadOnly"),
					AccessSettings: &itypes.AccessSettingList{
						AccessSetting: []*itypes.AccessSetting{
							{
								Subject: &itypes.LocalSubject{
									HREF: "https://api.cloudavenue.org/api/admin/user/22222222-2222-2222-2222-222222222222",
									Name: mockACLUserName,
									Type: "urn:vcloud:orguser",
								},
								AccessLevel: mockACLAccessLevel,
							},
						},
					},
				})
			})

			result, err := client.GetAdminCatalogACL(t.Context(), types.ParamsGetAdminCatalogACL{ID: tt.idOrName})
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectShared, result.IsSharedToEveryone)
			assert.NotNil(t, result.EveryoneAccessLevel)
			assert.Equal(t, "ReadOnly", *result.EveryoneAccessLevel)
			assert.Len(t, result.SharedWith, 1)
			assert.Equal(t, mockACLUserName, result.SharedWith[0].SubjectName)
			assert.Equal(t, mockACLAccessLevel, result.SharedWith[0].AccessLevel)
		})
	}
}

func TestGetAdminCatalogACL_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetAdminCatalogACL())

	_, err := client.GetAdminCatalogACL(t.Context(), types.ParamsGetAdminCatalogACL{})
	assert.Error(t, err)
}

func TestSetAdminCatalogACL(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.SetAdminCatalogACL())
	ms.SetResponseFunc(endpoints.SetAdminCatalogACL(), func(w http.ResponseWriter, r *http.Request) {
		var body itypes.ControlAccessParams
		assert.NoError(t, xml.NewDecoder(r.Body).Decode(&body))

		assert.Equal(t, false, body.IsSharedToEveryone)
		assert.NotNil(t, body.AccessSettings)
		assert.Len(t, body.AccessSettings.AccessSetting, 1)
		setting := body.AccessSettings.AccessSetting[0]
		assert.NotNil(t, setting.Subject)
		assert.Equal(t, "https://api.cloudavenue.org/api/admin/user/22222222-2222-2222-2222-222222222222", setting.Subject.HREF)
		assert.Equal(t, "urn:vcloud:orguser", setting.Subject.Type)
		assert.Equal(t, mockACLAccessLevel, setting.AccessLevel)

		xmlResponse(w, itypes.ControlAccessParams{
			Xmlns:               body.Xmlns,
			IsSharedToEveryone:  false,
			EveryoneAccessLevel: nil,
			AccessSettings: &itypes.AccessSettingList{
				AccessSetting: []*itypes.AccessSetting{
					{
						Subject: &itypes.LocalSubject{
							HREF: setting.Subject.HREF,
							Name: setting.Subject.Name,
							Type: setting.Subject.Type,
						},
						ExternalSubject: &itypes.ExternalSubject{
							IsUser:    true,
							SubjectID: mockACLUserURN,
						},
						AccessLevel: setting.AccessLevel,
					},
				},
			},
		})
	})

	result, err := client.SetAdminCatalogACL(t.Context(), types.ParamsSetAdminCatalogACL{
		ID:                 mockCatalogID,
		IsSharedToEveryone: false,
		SharedWith: []types.ParamsAdminCatalogACLSubject{
			{
				UserID:      mockACLUserURN,
				AccessLevel: mockACLAccessLevel,
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsSharedToEveryone)
	assert.Nil(t, result.EveryoneAccessLevel)
	assert.Len(t, result.SharedWith, 1)
	assert.Equal(t, mockACLAccessLevel, result.SharedWith[0].AccessLevel)
	assert.Equal(t, mockACLUserURN, result.SharedWith[0].UserID)
}

func TestSetAdminCatalogACL_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.SetAdminCatalogACL())

	_, err := client.SetAdminCatalogACL(t.Context(), types.ParamsSetAdminCatalogACL{})
	assert.Error(t, err)
}

func stringPtr(s string) *string {
	return &s
}
