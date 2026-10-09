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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestLDAPConfigToAPIRequestMatchesVCDShape(t *testing.T) {
	port := 636
	ssl := true

	body, err := LDAPConfigToAPIRequest(types.ParamsTestLDAP{
		Host:                    "ldap.example.com",
		Port:                    &port,
		SSLEnabled:              &ssl,
		BaseDN:                  "dc=example,dc=com",
		BindUser:                "cn=bind,dc=example,dc=com",
		BindPassword:            "secret",
		AuthenticationMechanism: "SIMPLE",
		GroupSearchBase:         "ou=groups,dc=example,dc=com",
		UserObjectClass:         "inetOrgPerson",
		UserNameAttribute:       "uid",
		GroupObjectClass:        "groupOfNames",
		GroupNameAttribute:      "cn",
		GroupMemberAttribute:    "member",
	})
	require.NoError(t, err)

	//nolint:gosec // Test verifies VMware LDAP payload shape; password field name/value are synthetic fixtures.
	encoded, err := json.Marshal(body)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(encoded, &payload))
	assert.Equal(t, "ldap.example.com", payload["hostName"])
	assert.EqualValues(t, 636, payload["port"])
	assert.Equal(t, true, payload["isSsl"])
	assert.Equal(t, "dc=example,dc=com", payload["searchBase"])
	assert.Equal(t, "cn=bind,dc=example,dc=com", payload["userName"])
	assert.Equal(t, "secret", payload["password"])
	assert.Equal(t, true, payload["isGroupSearchBaseEnabled"])
	assert.NotContains(t, payload, "config")

	userAttributes, ok := payload["userAttributes"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "inetOrgPerson", userAttributes["objectClass"])
	assert.Equal(t, "uid", userAttributes["userName"])

	groupAttributes, ok := payload["groupAttributes"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "groupOfNames", groupAttributes["objectClass"])
	assert.Equal(t, "cn", groupAttributes["groupName"])
	assert.Equal(t, "member", groupAttributes["membership"])
}

func TestLDAPRequestSerializesExplicitFalseBooleans(t *testing.T) {
	ssl := false
	acceptAll := false
	pagedDisabled := false
	groupSearchEnabled := false
	externalKerberos := false

	body := APIRequestLDAPTest{
		IsSSL:                    &ssl,
		IsSSLAcceptAll:           &acceptAll,
		PagedSearchDisabled:      &pagedDisabled,
		IsGroupSearchBaseEnabled: &groupSearchEnabled,
		UseExternalKerberos:      &externalKerberos,
	}

	//nolint:gosec // Test verifies explicit false serialization on LDAP request booleans; fixture is non-sensitive.
	encoded, err := json.Marshal(body)
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(encoded, &payload))
	assert.Equal(t, false, payload["isSsl"])
	assert.Equal(t, false, payload["isSslAcceptAll"])
	assert.Equal(t, false, payload["pagedSearchDisabled"])
	assert.Equal(t, false, payload["isGroupSearchBaseEnabled"])
	assert.Equal(t, false, payload["useExternalKerberos"])
}
