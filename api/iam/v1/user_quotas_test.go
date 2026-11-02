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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
)

func TestGetUserQuotas(t *testing.T) {
	client, ms := newClient(t)
	ep := endpoints.GetAdminUser()
	ms.CleanResponse(ep)
	ms.SetResponseFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/user/22740512-4f4c-4d2a-9d05-2d35294bcbc1", r.URL.Path)
		assert.Equal(t, "application/vnd.vmware.admin.user+xml;version=39.1", r.Header.Get("Accept"))
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, `<User xmlns="http://www.vmware.com/vcloud/v1.5" href="https://host/api/admin/user/22740512-4f4c-4d2a-9d05-2d35294bcbc1" id="urn:vcloud:user:22740512-4f4c-4d2a-9d05-2d35294bcbc1" name="user1"><FullName>One</FullName><IsEnabled>true</IsEnabled><DeployedVmQuota>10</DeployedVmQuota><StoredVmQuota>20</StoredVmQuota></User>`) //nolint:lll
	})

	result, err := client.GetUserQuotas(t.Context(), ParamsGetUserQuotas{ID: "urn:vcloud:user:22740512-4f4c-4d2a-9d05-2d35294bcbc1"})
	require.NoError(t, err)
	assert.Equal(t, "urn:vcloud:user:22740512-4f4c-4d2a-9d05-2d35294bcbc1", result.ID)
	assert.Equal(t, "user1", result.Name)
	assert.Equal(t, 10, result.DeployedVMQuota)
	assert.Equal(t, 20, result.StoredVMQuota)
}

func TestGetUserQuotasAcceptsPlainUUID(t *testing.T) {
	client, ms := newClient(t)
	ep := endpoints.GetAdminUser()
	ms.CleanResponse(ep)
	ms.SetResponseFunc(ep, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/admin/user/22740512-4f4c-4d2a-9d05-2d35294bcbc1", r.URL.Path)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, `<User name="user1" id="urn:vcloud:user:22740512-4f4c-4d2a-9d05-2d35294bcbc1"><DeployedVmQuota>1</DeployedVmQuota><StoredVmQuota>2</StoredVmQuota></User>`)
	})

	result, err := client.GetUserQuotas(t.Context(), ParamsGetUserQuotas{ID: "22740512-4f4c-4d2a-9d05-2d35294bcbc1"})
	require.NoError(t, err)
	assert.Equal(t, 1, result.DeployedVMQuota)
}

func TestGetUserQuotasValidation(t *testing.T) {
	client, _ := newClient(t)

	_, err := client.GetUserQuotas(t.Context(), ParamsGetUserQuotas{})
	require.ErrorContains(t, err, "id is required")
}
