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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestTestLDAP(t *testing.T) {
	t.Run("portal flat payload and query", func(t *testing.T) {
		client, ms := newClient(t)
		port := 636
		sslEnabled := true

		ms.CleanResponse(endpoints.TestLDAP())
		ms.SetResponseFunc(endpoints.TestLDAP(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "bind-user", r.URL.Query().Get("username"))

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
}

func TestSearchLDAPUsers(t *testing.T) {
	t.Run("bare array response", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SearchLDAPUsers())
		ms.SetResponseFunc(endpoints.SearchLDAPUsers(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "john", r.URL.Query().Get("q"))
			assert.Empty(t, r.URL.Query().Get("maxResults"))
			assert.Empty(t, r.URL.Query().Get("pageSize"))
			assert.Empty(t, r.URL.Query().Get("page"))

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

	t.Run("reject undocumented params", func(t *testing.T) {
		client, _ := newClient(t)

		result, err := client.SearchLDAPUsers(t.Context(), types.ParamsSearchLDAP{MaxResults: "10", PageSize: "20", Page: "2"})
		assert.Nil(t, result)
		assert.EqualError(t, err, "IAM.SearchLDAPUsers: validate: unsupported params for endpoint: maxResults, pageSize, page")
	})
}

func TestSearchLDAPGroups(t *testing.T) {
	t.Run("bare array response", func(t *testing.T) {
		client, ms := newClient(t)

		ms.CleanResponse(endpoints.SearchLDAPGroups())
		ms.SetResponseFunc(endpoints.SearchLDAPGroups(), func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "admins", r.URL.Query().Get("q"))
			assert.Empty(t, r.URL.Query().Get("maxResults"))
			assert.Empty(t, r.URL.Query().Get("pageSize"))
			assert.Empty(t, r.URL.Query().Get("page"))

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"name":"admins","id":"cn=admins,ou=groups,dc=example,dc=com"}]`))
		})

		result, err := client.SearchLDAPGroups(t.Context(), types.ParamsSearchLDAP{Filter: "admins"})
		assert.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "admins", result[0].Name)
		assert.Equal(t, "cn=admins,ou=groups,dc=example,dc=com", result[0].DN)
	})

	t.Run("reject undocumented params", func(t *testing.T) {
		client, _ := newClient(t)

		result, err := client.SearchLDAPGroups(t.Context(), types.ParamsSearchLDAP{MaxResults: "10", PageSize: "20", Page: "2"})
		assert.Nil(t, result)
		assert.EqualError(t, err, "IAM.SearchLDAPGroups: validate: unsupported params for endpoint: maxResults, pageSize, page")
	})
}
