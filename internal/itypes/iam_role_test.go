/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIResponseTokenDecodesVCDFields(t *testing.T) {
	var response APIResponseToken
	require.NoError(t, json.Unmarshal([]byte(`{
		"id":"urn:vcloud:token:11111111-1111-1111-1111-111111111111",
		"name":"proxy-token",
		"token":"secret-token",
		"expirationTimeUtc":"2026-10-09T12:00:00Z",
		"owner":{"id":"urn:vcloud:user:22222222-2222-2222-2222-222222222222","name":"alice"},
		"org":{"id":"urn:vcloud:org:33333333-3333-3333-3333-333333333333","name":"example"},
		"type":"PROXY",
		"enabled":true,
		"role":{"id":"urn:vcloud:role:44444444-4444-4444-4444-444444444444","name":"Organization Administrator"}
	}`), &response))

	model := response.ToModel()
	require.NotNil(t, model)
	assert.Equal(t, "PROXY", model.Type)
	assert.Equal(t, "secret-token", model.Token)
	assert.Equal(t, "2026-10-09T12:00:00Z", model.ExpirationTimeUTC)
	assert.Equal(t, "urn:vcloud:user:22222222-2222-2222-2222-222222222222", model.OwnerID)
	assert.Equal(t, "alice", model.OwnerName)
	assert.Equal(t, "urn:vcloud:org:33333333-3333-3333-3333-333333333333", model.OrgID)
	assert.Equal(t, "example", model.OrgName)
}

func TestAPIResponseTokenPreservesExplicitFalse(t *testing.T) {
	var response APIResponseToken
	require.NoError(t, json.Unmarshal([]byte(`{"enabled":false,"requireRotation":false}`), &response))

	model := response.ToModel()
	require.NotNil(t, model)
	assert.False(t, model.Enabled)
	require.NotNil(t, model.RequireRotation)
	assert.False(t, *model.RequireRotation)
}

func TestAPIResponseGlobalRoleDecodesPublishableRoleFields(t *testing.T) {
	var response APIResponseGlobalRole
	require.NoError(t, json.Unmarshal([]byte(`{
		"id":"urn:vcloud:globalRole:11111111-1111-1111-1111-111111111111",
		"name":"Sub-Provider Administrator",
		"description":"Provider role",
		"bundleKey":"sub-provider-administrator",
		"readOnly":true,
		"publishAll":true,
		"numberOfTenants":3,
		"canPublish":true,
		"isInternalRole":false,
		"rights":[{"id":"urn:vcloud:right:22222222-2222-2222-2222-222222222222","name":"Org: View","rightType":"VIEW","isPublishable":true}]
	}`), &response))

	model := response.ToModel()
	require.NotNil(t, model)
	assert.Equal(t, "sub-provider-administrator", model.BundleKey)
	assert.True(t, model.ReadOnly)
	assert.True(t, model.PublishAll)
	assert.Len(t, model.Rights, 1)
	assert.True(t, model.Rights[0].IsPublishable)
}

func TestAPIRequestTokenSerializesExtensionFields(t *testing.T) {
	body := APIRequestToken{
		Name:        "extension-token",
		Type:        "EXTENSION",
		ExtensionID: "urn:vcloud:extension:11111111-1111-1111-1111-111111111111",
		Enabled:     true,
	}

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"extension-token","type":"EXTENSION","extensionId":"urn:vcloud:extension:11111111-1111-1111-1111-111111111111","role":{"id":"","name":""},"enabled":true}`, string(encoded))
}

func TestAPIRequestTokenSerializesExplicitFalse(t *testing.T) {
	body := APIRequestToken{Enabled: false}

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"enabled":false,"role":{"id":"","name":""}}`, string(encoded))
}
