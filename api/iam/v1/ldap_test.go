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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	sdkerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestTestLDAP(t *testing.T) {
	t.Run("portal flat payload and query", func(t *testing.T) {
		client, ms := newClient(t)
		port := 636
		sslEnabled := true

		ms.CleanResponse(endpoints.TestLDAP())
		ms.SetResponseFunc(endpoints.TestLDAP(), func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.Query().Get("username"))

			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.NotContains(t, body, "config")
			assert.Equal(t, "ldap.example.com", body["hostName"])
			assert.EqualValues(t, 636, body["port"])
			assert.Equal(t, true, body["isSsl"])
			assert.Equal(t, "dc=example,dc=com", body["searchBase"])
			assert.Equal(t, "bind-user", body["userName"])
			assert.Equal(t, "secret", body["password"])
			assert.Equal(t, "SIMPLE", body["authenticationMechanism"])
			assert.Equal(t, "ou=groups,dc=example,dc=com", body["groupSearchBase"])
			assert.Equal(t, true, body["isGroupSearchBaseEnabled"])

			userAttributes, ok := body["userAttributes"].(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "inetOrgPerson", userAttributes["objectClass"])
			assert.Equal(t, "uid", userAttributes["userName"])

			groupAttributes, ok := body["groupAttributes"].(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "groupOfNames", groupAttributes["objectClass"])
			assert.Equal(t, "cn", groupAttributes["groupName"])
			assert.Equal(t, "member", groupAttributes["membership"])

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"connectionTest":{"successful":true},"settingsTest":[{"attribute":"userName","result":"ok","successful":true}]}`))
		})

		result, err := client.TestLDAP(t.Context(), types.ParamsTestLDAP{
			Host:                    "ldap.example.com",
			Port:                    &port,
			BindUser:                "bind-user",
			BindPassword:            "secret",
			BaseDN:                  "dc=example,dc=com",
			GroupSearchBase:         "ou=groups,dc=example,dc=com",
			UserObjectClass:         "inetOrgPerson",
			GroupObjectClass:        "groupOfNames",
			UserNameAttribute:       "uid",
			GroupNameAttribute:      "cn",
			GroupMemberAttribute:    "member",
			AuthenticationMechanism: "SIMPLE",
			SSLEnabled:              &sslEnabled,
		})
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, "ok", result.Message)
	})

	t.Run("nil Enabled accepted", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.TestLDAP())
		ms.SetResponseFunc(endpoints.TestLDAP(), func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"connectionTest":{"successful":true}}`))
		})

		result, err := client.TestLDAP(t.Context(), types.ParamsTestLDAP{
			Host:              "ldap.example.com",
			Enabled:           nil,
			BaseDN:            "dc=example,dc=com",
			UserNameAttribute: "uid",
		})
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
	})

	t.Run("reject unsupported params", func(t *testing.T) {
		trueVal := true
		falseVal := false
		timeout := 5
		zero := 0

		tests := []struct {
			name   string
			params types.ParamsTestLDAP
			want   string
		}{
			{
				name: "userSearchBase",
				params: types.ParamsTestLDAP{
					Host:           "ldap.example.com",
					UserSearchBase: "ou=people,dc=example,dc=com",
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: userSearchBase",
			},
			{
				name: "sslTrustCertificate",
				params: types.ParamsTestLDAP{
					Host:                "ldap.example.com",
					SSLTrustCertificate: "-----BEGIN CERTIFICATE-----",
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: sslTrustCertificate",
			},
			{
				name: "connectionTimeout",
				params: types.ParamsTestLDAP{
					Host:              "ldap.example.com",
					ConnectionTimeout: &timeout,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: connectionTimeout",
			},
			{
				name: "readTimeout",
				params: types.ParamsTestLDAP{
					Host:        "ldap.example.com",
					ReadTimeout: &timeout,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: readTimeout",
			},
			{
				name: "enabled true",
				params: types.ParamsTestLDAP{
					Host:    "ldap.example.com",
					Enabled: &trueVal,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: enabled",
			},
			{
				name: "enabled false",
				params: types.ParamsTestLDAP{
					Host:    "ldap.example.com",
					Enabled: &falseVal,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: enabled",
			},
			{
				name: "all five",
				params: types.ParamsTestLDAP{
					Host:                "ldap.example.com",
					UserSearchBase:      "ou=people,dc=example,dc=com",
					SSLTrustCertificate: "-----BEGIN CERTIFICATE-----",
					ConnectionTimeout:   &timeout,
					ReadTimeout:         &timeout,
					Enabled:             &trueVal,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: userSearchBase, sslTrustCertificate, connectionTimeout, readTimeout, enabled",
			},
			{
				name: "zero value pointers rejected",
				params: types.ParamsTestLDAP{
					Host:              "ldap.example.com",
					ConnectionTimeout: &zero,
					ReadTimeout:       &zero,
				},
				want: "IAM.TestLDAP: validate: unsupported params for endpoint: connectionTimeout, readTimeout",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				client, _ := newClient(t)

				result, err := client.TestLDAP(t.Context(), tc.params)
				assert.Nil(t, result)
				assert.EqualError(t, err, tc.want)
				assert.ErrorIs(t, err, types.ErrUnsupportedLDAPTestParams)
			})
		}
	})
}

func TestSyncLDAP(t *testing.T) {
	t.Run("no content success", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SyncLDAP())
		ms.SetResponseFunc(endpoints.SyncLDAP(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Empty(t, body)
			w.WriteHeader(http.StatusNoContent)
		})

		err := client.SyncLDAP(t.Context(), types.ParamsSyncLDAP{})
		assert.NoError(t, err)
	})

	t.Run("propagates API error", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SyncLDAP())
		ms.SetResponseFunc(endpoints.SyncLDAP(), func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		})

		err := client.SyncLDAP(t.Context(), types.ParamsSyncLDAP{})
		var apiErr *sdkerrors.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}

func TestSearchLDAPUsers(t *testing.T) {
	t.Run("bare array response", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SearchLDAPUsers())
		ms.SetResponseFunc(endpoints.SearchLDAPUsers(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "john", r.URL.Query().Get("q"))

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"username":"john","fullname":"John Doe","email":"john@example.com"}]`))
		})

		result, err := client.SearchLDAPUsers(t.Context(), types.ParamsSearchLDAP{Filter: "john"})
		assert.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "john", result[0].Name)
		assert.Equal(t, "John Doe", result[0].FullName)
		assert.Equal(t, "john@example.com", result[0].Email)
	})
}

func TestSearchLDAPGroups(t *testing.T) {
	t.Run("bare array response", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SearchLDAPGroups())
		ms.SetResponseFunc(endpoints.SearchLDAPGroups(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "admins", r.URL.Query().Get("q"))

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"name":"admins","id":"cn=admins,ou=groups,dc=example,dc=com"}]`))
		})

		result, err := client.SearchLDAPGroups(t.Context(), types.ParamsSearchLDAP{Filter: "admins"})
		assert.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "admins", result[0].Name)
	})
}
