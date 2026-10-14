/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestCreateOAuthClient(t *testing.T) {
	client, ms := newClient(t)
	ep := endpoints.CreateOAuthClient()
	ms.CleanResponse(ep)
	ms.SetResponseFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/oauth/tenant/tenant-123/register", r.URL.Path)
		assert.Equal(t, "application/json;version=39.1", r.Header.Get("Accept"))
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "my-client", body["client_name"])
		assert.Equal(t, "software-1", body["software_id"])
		assert.NotContains(t, body, "client_secret")
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(itypes.APIResponseOAuthClient{
			ClientName: "my-client", ClientID: "client-123", GrantTypes: []string{"client_credentials"},
			TokenEndpointAuthMethod: "client_secret_basic",
		}))
	})

	softwareID := "software-1"
	result, err := client.CreateOAuthClient(t.Context(), types.ParamsCreateOAuthClient{
		TenantID: "tenant-123", ClientName: "my-client", SoftwareID: &softwareID,
	})
	require.NoError(t, err)
	assert.Equal(t, "client-123", result.ClientID)
	assert.Equal(t, []string{"client_credentials"}, result.GrantTypes)
	assert.Nil(t, result.ClientURI)
}

func TestCreateOAuthToken(t *testing.T) {
	client, ms := newClient(t)
	ep := endpoints.CreateOAuthToken()
	ms.CleanResponse(ep)
	ms.SetResponseFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/oauth/tenant/tenant-123/token", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json;version=39.1", r.Header.Get("Accept"))
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		form, err := url.ParseQuery(string(body))
		require.NoError(t, err)
		assert.Equal(t, "client-123", form.Get("client_id"))
		assert.Equal(t, "urn:ietf:params:oauth:grant-type:jwt-bearer", form.Get("grant_type"))
		assert.Equal(t, "jwt-value", form.Get("assertion"))
		assert.Empty(t, form.Get("device_code"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(itypes.APIResponseOAuthToken{ //nolint: gosec
			AccessToken: "access-token", TokenType: "Bearer", ExpiresIn: 3600, RefreshToken: "refresh-token",
		}))
	})

	result, err := client.CreateOAuthToken(t.Context(), types.ParamsCreateOAuthToken{
		TenantID: "tenant-123", ClientID: "client-123",
		GrantType: "urn:ietf:params:oauth:grant-type:jwt-bearer", Assertion: "jwt-value",
	})
	require.NoError(t, err)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "Bearer", result.TokenType)
	assert.Equal(t, 3600, result.ExpiresIn)
	assert.Equal(t, "refresh-token", result.RefreshToken)
}
