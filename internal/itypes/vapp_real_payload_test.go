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

// realVAppPayload is a trimmed GET /api/vApp/vapp-<uuid> response with the
// vendor vApp+json media type: no top-level vdc/org/leaseSettings, numeric
// status, VDC in the rel="up" link, lease in the LeaseSettingsSectionType
// section and an owner user without id.
const realVAppPayload = `{
  "link": [
    {"href": "https://console.example.org/api/vApp/vapp-0704fe88-a09b-474c-8f89-fde3bd2736fb/action/deploy", "type": "application/vnd.vmware.vcloud.deployVAppParams+json", "rel": "deploy"},
    {"href": "https://console.example.org/api/vdc/cfb722cd-3a8b-4c0a-930e-8fe116d9ab5d", "type": "application/vnd.vmware.vcloud.vdc+xml", "rel": "up"},
    {"href": "https://console.example.org/api/vdc/cfb722cd-3a8b-4c0a-930e-8fe116d9ab5d", "type": "application/vnd.vmware.vcloud.vdc+json", "rel": "up"}
  ],
  "href": "https://console.example.org/api/vApp/vapp-0704fe88-a09b-474c-8f89-fde3bd2736fb",
  "type": "application/vnd.vmware.vcloud.vApp+json",
  "id": "urn:vcloud:vapp:0704fe88-a09b-474c-8f89-fde3bd2736fb",
  "description": null,
  "name": "test",
  "status": 1,
  "section": [
    {"_type": "LeaseSettingsSectionType", "required": false, "deploymentLeaseInSeconds": 0, "storageLeaseInSeconds": 3600},
    {"_type": "NetworkSectionType", "required": true, "network": []},
    {"_type": "NetworkConfigSectionType", "required": false, "networkConfig": []}
  ],
  "deployed": false,
  "owner": {
    "href": null,
    "type": "application/vnd.vmware.vcloud.owner+json",
    "user": {"href": "https://console.example.org/api/admin/user/e05d2dcb-95f1-4af8-8166-9f8a16b6d8d3", "id": null, "type": "application/vnd.vmware.admin.user+xml", "name": "system"}
  },
  "vAppParent": null
}`

func TestAPIResponseGetVApp_RealPayload(t *testing.T) {
	var resp APIResponseGetVApp
	require.NoError(t, json.Unmarshal([]byte(realVAppPayload), &resp))

	model := resp.ToModel()

	assert.Equal(t, "urn:vcloud:vapp:0704fe88-a09b-474c-8f89-fde3bd2736fb", model.ID)
	assert.Equal(t, "test", model.Name)
	assert.Equal(t, "RESOLVED", model.Status)
	assert.False(t, model.Deployed)

	// VDC derived from the rel="up" link
	assert.Equal(t, "urn:vcloud:vdc:cfb722cd-3a8b-4c0a-930e-8fe116d9ab5d", model.VDC)
}

func TestAPIResponseGetVApp_TopLevelFieldsStillWork(t *testing.T) {
	var resp APIResponseGetVApp
	require.NoError(t, json.Unmarshal([]byte(`{
		"id": "urn:vcloud:vapp:0704fe88-a09b-474c-8f89-fde3bd2736fb",
		"vdc": "urn:vcloud:vdc:11111111-1111-4111-8111-111111111111",
		"status": "POWERED_ON",
		"leaseSettings": {"deploymentLeaseInSeconds": 100, "storageLeaseInSeconds": 200},
		"owner": "urn:vcloud:user:22222222-2222-4222-8222-222222222222"
	}`), &resp))

	model := resp.ToModel()
	assert.Equal(t, "urn:vcloud:vdc:11111111-1111-4111-8111-111111111111", model.VDC)
	assert.Equal(t, "POWERED_ON", model.Status)
	require.NotNil(t, model.DeploymentLeaseInSeconds)
	assert.Equal(t, 100, *model.DeploymentLeaseInSeconds)
	assert.Equal(t, "urn:vcloud:user:22222222-2222-4222-8222-222222222222", model.Owner)
}

func TestAPIResponseGetVDCStorageProfiles_ArrayOrObject(t *testing.T) {
	const profile = `{"id": "urn:vcloud:vdcstorageProfile:33333333-3333-4333-8333-333333333333", "name": "gold", "limit": 500}`

	t.Run("plain array (real API)", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`[`+profile+`]`), &r))
		require.Len(t, r.StorageProfiles, 1)
		assert.Equal(t, "gold", r.StorageProfiles[0].Name)
	})

	t.Run("wrapped object", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`{"vdcStorageProfile": [`+profile+`]}`), &r))
		require.Len(t, r.StorageProfiles, 1)
		assert.Equal(t, "gold", r.StorageProfiles[0].Name)
	})

	t.Run("null", func(t *testing.T) {
		var r APIResponseGetVDCStorageProfiles
		require.NoError(t, json.Unmarshal([]byte(`null`), &r))
		assert.Empty(t, r.StorageProfiles)
	})
}
