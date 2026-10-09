/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iendpoints_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func TestLDAPEndpointPaths(t *testing.T) {
	tests := []struct {
		name string
		get  func() *cav.Endpoint
		path string
	}{
		{name: "TestLDAP", get: endpoints.TestLDAP, path: "/cloudapi/1.0.0/ldap/{orgId}/test"},
		{name: "SyncLDAP", get: endpoints.SyncLDAP, path: "/cloudapi/1.0.0/ldap/{orgId}/sync"},
		{name: "SearchLDAPUsers", get: endpoints.SearchLDAPUsers, path: "/cloudapi/1.0.0/ldap/{orgId}/search/user"},
		{name: "SearchLDAPGroups", get: endpoints.SearchLDAPGroups, path: "/cloudapi/1.0.0/ldap/{orgId}/search/group"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := tt.get()
			require.NotNil(t, ep)
			assert.Equal(t, cav.BackendVMware, ep.Backend)
			assert.Equal(t, tt.path, ep.PathTemplate)
			require.Len(t, ep.PathParams, 1)
			assert.Equal(t, "orgId", ep.PathParams[0].Name)
		})
	}
}
