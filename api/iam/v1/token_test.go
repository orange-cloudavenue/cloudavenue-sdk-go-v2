/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListTokensDoesNotRequireRemovedPathParam(t *testing.T) {
	client, ms := newClient(t)
	assert.Equal(t, "/cloudapi/1.0.0/tokens", endpoints.ListTokens().PathTemplate)

	ms.CleanResponse(endpoints.ListTokens())
	ms.SetResponseFunc(endpoints.ListTokens(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/cloudapi/1.0.0/tokens", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(itypes.APIResponseListTokens{
			Tokens: []itypes.APIResponseToken{{ID: "token-1", Name: "first"}},
		}))
	})

	result, err := client.ListTokens(t.Context())

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "token-1", result[0].ID)
}

func TestTokenCollectionEndpointTemplates(t *testing.T) {
	assert.Equal(t, "/cloudapi/1.0.0/tokens", endpoints.ListTokens().PathTemplate)
	assert.Equal(t, "/cloudapi/1.0.0/tokens", endpoints.CreateToken().PathTemplate)
	assert.Equal(t, "/cloudapi/1.0.0/tokens/{id}", endpoints.GetToken().PathTemplate)
	assert.Equal(t, "/cloudapi/1.0.0/sessions/{id}/token", cav.MustGetEndpoint("GetSessionToken").PathTemplate)
}

func TestTokenIDIsInterpolatedIntoRoute(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		request  func(*Client) error
		endpoint func() *cav.Endpoint
		method   string
	}{
		{
			name: "get token",
			path: "/cloudapi/1.0.0/tokens/token-123",
			request: func(client *Client) error {
				_, err := client.GetToken(t.Context(), types.ParamsGetToken{ID: "token-123"})
				return err
			},
			endpoint: endpoints.GetToken,
			method:   http.MethodGet,
		},
		{
			name: "get session token",
			path: "/cloudapi/1.0.0/sessions/session-123/token",
			request: func(client *Client) error {
				_, err := client.GetSessionToken(t.Context(), types.ParamsGetToken{ID: "session-123"})
				return err
			},
			endpoint: func() *cav.Endpoint { return cav.MustGetEndpoint("GetSessionToken") },
			method:   http.MethodGet,
		},
		{
			name: "update",
			path: "/cloudapi/1.0.0/tokens/token-123",
			request: func(client *Client) error {
				_, err := client.UpdateToken(t.Context(), types.ParamsUpdateToken{ID: "token-123"})
				return err
			},
			endpoint: endpoints.UpdateToken,
			method:   http.MethodPut,
		},
		{
			name: "delete",
			path: "/cloudapi/1.0.0/tokens/token-123",
			request: func(client *Client) error {
				return client.DeleteToken(t.Context(), types.ParamsDeleteToken{ID: "token-123"})
			},
			endpoint: endpoints.DeleteToken,
			method:   http.MethodDelete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ms := newClient(t)
			ms.CleanResponse(tt.endpoint())
			ms.SetResponseFunc(tt.endpoint(), func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				if tt.method != http.MethodDelete {
					_ = json.NewEncoder(w).Encode(itypes.APIResponseToken{ID: "token-123"})
				}
			})

			require.NoError(t, tt.request(client))
		})
	}
}

func TestCreateAPIToken(t *testing.T) {
	client, ms := newClient(t)
	registerEp := endpoints.CreateOAuthClient()
	tokenEp := endpoints.CreateOAuthToken()
	ms.CleanResponse(registerEp)
	ms.CleanResponse(tokenEp)

	ms.SetResponseFunc(registerEp, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/oauth/tenant/tenant-123/register", r.URL.Path)
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "my-token", body["client_name"])
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(itypes.APIResponseOAuthClient{
			ClientName: "my-token", ClientID: "client-123",
		}))
	})
	ms.SetResponseFunc(tokenEp, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/oauth/tenant/tenant-123/token", r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		form, err := url.ParseQuery(string(body))
		require.NoError(t, err)
		assert.Equal(t, "client-123", form.Get("client_id"))
		assert.Equal(t, "urn:ietf:params:oauth:grant-type:jwt-bearer", form.Get("grant_type"))
		assert.Equal(t, "mock-access-token", form.Get("assertion"))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(itypes.APIResponseOAuthToken{ //nolint: gosec
			AccessToken: "initial-access", TokenType: "Bearer", ExpiresIn: 3600, RefreshToken: "refresh-token",
		}))
	})

	result, err := client.CreateAPIToken(t.Context(), types.ParamsCreateAPIToken{
		TenantID: "tenant-123", Name: "my-token",
	})
	require.NoError(t, err)
	assert.Equal(t, "urn:vcloud:token:client-123", result.ID)
	assert.Equal(t, "my-token", result.Name)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.Equal(t, "initial-access", result.InitialAccessToken)
}

func TestCreateAPITokenValidation(t *testing.T) {
	client, _ := newClient(t)

	_, err := client.CreateAPIToken(t.Context(), types.ParamsCreateAPIToken{Name: "my-token"})
	require.ErrorContains(t, err, "tenant_id is required")

	_, err = client.CreateAPIToken(t.Context(), types.ParamsCreateAPIToken{TenantID: "tenant-123"})
	require.ErrorContains(t, err, "name is required")
}
