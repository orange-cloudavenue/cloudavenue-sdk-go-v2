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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetToken(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetToken())
	ms.SetResponseFunc(endpoints.GetToken(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cloudapi/1.0.0/tokens/token-1", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"token-1","name":"test-token","token":"secret","enabled":true,"role":{"id":"role-1","name":"Organization Administrator"}}`))
	})

	result, err := client.GetToken(t.Context(), types.ParamsGetToken{ID: "token-1"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "token-1", result.ID)
	assert.Equal(t, "test-token", result.Name)
	assert.Equal(t, "secret", result.Token)
	assert.Equal(t, "Organization Administrator", result.RoleName)
}

func TestGetSessionToken(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetSessionToken())
	ms.SetResponseFunc(endpoints.GetSessionToken(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cloudapi/1.0.0/sessions/session-1/token", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"token-1","name":"session-token"}`))
	})

	result, err := client.GetSessionToken(t.Context(), types.ParamsGetToken{ID: "session-1"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "token-1", result.ID)
}

func TestListTokens(t *testing.T) {
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListTokens())
	ms.SetResponseFunc(endpoints.ListTokens(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/cloudapi/1.0.0/tokens", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tokens":[{"id":"token-1","name":"test-token","enabled":true}]}`))
	})

	result, err := client.ListTokens(t.Context())
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "token-1", result[0].ID)
	assert.Equal(t, "test-token", result[0].Name)
}
