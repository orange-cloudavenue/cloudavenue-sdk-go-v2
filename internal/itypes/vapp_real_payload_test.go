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
	"encoding/xml"
	"strings"
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

	// Lease read from the section, a 0 lease is a real value (never expires)
	require.NotNil(t, model.DeploymentLeaseInSeconds)
	assert.Equal(t, 0, *model.DeploymentLeaseInSeconds)
	require.NotNil(t, model.StorageLeaseInSeconds)
	assert.Equal(t, 3600, *model.StorageLeaseInSeconds)

	// Owner user href converted to a user URN
	assert.Equal(t, "urn:vcloud:user:e05d2dcb-95f1-4af8-8166-9f8a16b6d8d3", model.Owner)
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

func TestVAppOwnerRef_UserWithoutIDOrUUID(t *testing.T) {
	var resp APIResponseGetVApp
	require.NoError(t, json.Unmarshal([]byte(`{"owner": {"user": {"href": "https://console.example.org/api/admin/user/not-a-uuid"}}}`), &resp))
	assert.Equal(t, "https://console.example.org/api/admin/user/not-a-uuid", string(resp.Owner))
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

func TestAPIResponseVAppNetworkConfig_LinkChildGivesTheHREF(t *testing.T) {
	const uuid = "44444444-4444-4444-8444-444444444444"

	t.Run("no id nor href attribute: the Link href is used", func(t *testing.T) {
		var section APIResponseVAppNetworkConfigSection
		require.NoError(t, xml.Unmarshal([]byte(`<NetworkConfigSection xmlns="http://www.vmware.com/vcloud/v1.5">
			<NetworkConfig networkName="net1">
				<Link rel="repair" href="https://console.example.org/api/network/`+uuid+`"/>
				<Configuration><FenceMode>isolated</FenceMode></Configuration>
			</NetworkConfig>
		</NetworkConfigSection>`), &section))

		model := section.ToModel()
		require.Len(t, model.NetworkConfig, 1)
		assert.Equal(t, "net1", model.NetworkConfig[0].NetworkName)
		assert.Empty(t, model.NetworkConfig[0].ID)
		assert.Equal(t, "https://console.example.org/api/network/"+uuid, model.NetworkConfig[0].HREF)
	})

	t.Run("the href attribute wins over the Link", func(t *testing.T) {
		var section APIResponseVAppNetworkConfigSection
		require.NoError(t, xml.Unmarshal([]byte(`<NetworkConfigSection xmlns="http://www.vmware.com/vcloud/v1.5">
			<NetworkConfig networkName="net1" href="https://console.example.org/api/vAppNetwork/`+uuid+`">
				<Link rel="repair" href="https://console.example.org/api/other"/>
			</NetworkConfig>
		</NetworkConfigSection>`), &section))

		assert.Equal(t, "https://console.example.org/api/vAppNetwork/"+uuid, section.ToModel().NetworkConfig[0].HREF)
	})
}

// The section is read, modified and sent back as a whole: every element of the
// NetworkConfigSection schema must survive the round trip, in the schema order
// (VMware Cloud Director API, NetworkConfigSection).
func TestAPIResponseVAppNetworkConfigSection_RoundTripKeepsSchemaElements(t *testing.T) {
	const in = `<NetworkConfigSection xmlns="http://www.vmware.com/vcloud/v1.5">
		<NetworkConfig networkName="net1">
			<Link rel="repair" href="https://console.example.org/api/network/44444444-4444-4444-8444-444444444444"/>
			<Description>desc</Description>
			<Configuration>
				<BackwardCompatibilityMode>false</BackwardCompatibilityMode>
				<IpScopes><IpScope><IsInherited>false</IsInherited><Gateway>10.0.0.1</Gateway></IpScope></IpScopes>
				<ParentNetwork href="https://console.example.org/api/network/55555555-5555-4555-8555-555555555555" name="org"/>
				<FenceMode>bridged</FenceMode>
				<RetainNetInfoAcrossDeployments>false</RetainNetInfoAcrossDeployments>
				<SubInterface>false</SubInterface>
				<DistributedInterface>false</DistributedInterface>
				<ServiceInterface>false</ServiceInterface>
				<GuestVlanAllowed>true</GuestVlanAllowed>
				<Connected>true</Connected>
				<DualStackNetwork>false</DualStackNetwork>
			</Configuration>
			<IsDeployed>true</IsDeployed>
		</NetworkConfig>
	</NetworkConfigSection>`

	var section APIResponseVAppNetworkConfigSection
	require.NoError(t, xml.Unmarshal([]byte(in), &section))

	out, err := xml.Marshal(section)
	require.NoError(t, err)

	// Every element is present, in the schema order.
	ordered := []string{
		"<Link ", "<Description>", "<Configuration>",
		"<BackwardCompatibilityMode>", "<IpScopes>", "<ParentNetwork ", "<FenceMode>",
		"<RetainNetInfoAcrossDeployments>", "<SubInterface>", "<DistributedInterface>",
		"<ServiceInterface>", "<GuestVlanAllowed>", "<Connected>", "<DualStackNetwork>",
		"</Configuration>", "<IsDeployed>",
	}
	last := -1
	for _, element := range ordered {
		idx := strings.Index(string(out), element)
		require.GreaterOrEqual(t, idx, 0, "missing %s in %s", element, out)
		assert.Greater(t, idx, last, "%s is out of order in %s", element, out)
		last = idx
	}

	// IsDeployed is read from NetworkConfig into the public model.
	model := section.ToModel()
	require.Len(t, model.NetworkConfig, 1)
	if assert.NotNil(t, model.NetworkConfig[0].Configuration.IsDeployed) {
		assert.True(t, *model.NetworkConfig[0].Configuration.IsDeployed)
	}
}
