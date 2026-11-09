/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package cav

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

func TestCerberusCredentialRefresh(t *testing.T) {
	var called bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, cerberusAuthPath, r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "client_credentials", r.Form.Get("grant_type"))
		assert.Equal(t, "mockuser", r.Form.Get("client_id"))
		assert.Equal(t, "mockpassword", r.Form.Get("client_secret"))
		assert.Equal(t, "tenant:"+mockOrg, r.Form.Get("scope"))
		w.Header().Set("Content-Type", "application/json")
		//nolint:gosec // deterministic fake token for transport/auth unit test
		require.NoError(t, json.NewEncoder(w).Encode(cerberusAuthResponse{
			AccessToken: "cerberus-access-token",
			TokenType:   cerberusTokenTypeBearer,
			ExpiresIn:   3600,
		}))
	}))
	t.Cleanup(srv.Close)

	cred, err := newCerberusAuthProvider(consoles.ConsoleName(mockOrg), mockOrg, "mockuser", "mockpassword")
	require.NoError(t, err)
	c := cred.(*cerberusAuthProvider)
	c.httpC.SetBaseURL(srv.URL)

	require.NoError(t, c.Refresh(t.Context()))
	assert.True(t, called)
	assert.Equal(t, map[string]string{"Authorization": "Bearer cerberus-access-token"}, c.Headers())

	authorization, err := c.CurrentAuthorization(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "Bearer cerberus-access-token", authorization)
	assert.True(t, c.IsInitialized())
}

func TestCerberusCredentialRefreshRejectsRedirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMultipleChoices)
	}))
	t.Cleanup(srv.Close)

	cred, err := newCerberusAuthProvider(consoles.ConsoleName(mockOrg), mockOrg, "mockuser", "mockpassword")
	require.NoError(t, err)
	c := cred.(*cerberusAuthProvider)
	c.httpC.SetBaseURL(srv.URL)
	c.accessToken = "expired-token"
	c.tokenType = cerberusTokenTypeBearer
	c.expiresAt = time.Now().Add(-time.Minute)

	err = c.Refresh(t.Context())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "300")
	assert.False(t, c.IsInitialized())
}

func TestCerberusCredentialRestoreSessionIsAtomicOnInvalidExpiry(t *testing.T) {
	c := &cerberusAuthProvider{
		organization: "original-organization",
		accessToken:  "original-token",
		tokenType:    cerberusTokenTypeBearer,
		expiresAt:    time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	err := c.restoreSession(map[string]string{
		"organization": "replacement-organization",
		"accessToken":  "replacement-token",
		"tokenType":    "Replacement",
		"expiresAt":    "not-a-timestamp",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "restore session: parse expiresAt")
	assert.Equal(t, "original-organization", c.organization)
	assert.Equal(t, "original-token", c.accessToken)
	assert.Equal(t, cerberusTokenTypeBearer, c.tokenType)
	assert.Equal(t, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), c.expiresAt)
}

func TestWithCloudAvenueCredentialUsesCerberusCredentialForInfrapi(t *testing.T) {
	cl, err := NewClient(mockOrg, WithCloudAvenueCredential("mockuser", "mockpassword"))
	require.NoError(t, err)

	c := cl.(*client)
	_, ok := c.clientsInitialized[ClientCerberus].getAuthProvider().(*cerberusAuthProvider)
	assert.True(t, ok)
	_, ok = c.clientsInitialized[ClientVmware].getAuthProvider().(*vmwareAuthProvider)
	assert.True(t, ok)
}
