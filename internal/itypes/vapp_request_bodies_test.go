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
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The expected documents are the ones produced by go-vcloud-director
// (SmallRecomposeVappParams and UpdateLeaseSettingsSection) for the same values:
// they are compared as text, not read back with the same structs.

func TestAPIRequestRecomposeVApp_Body(t *testing.T) {
	body, err := xml.Marshal(APIRequestRecomposeVApp{
		OVF:         ovfEnvelopeXMLNS,
		XSI:         "http://www.w3.org/2001/XMLSchema-instance",
		XMLNS:       vcloudXMLNS,
		Name:        "my-vapp",
		Deploy:      false,
		Description: "new description",
	})
	require.NoError(t, err)

	assert.Equal(t, `<RecomposeVAppParams xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.vmware.com/vcloud/v1.5" name="my-vapp" deploy="false"><Description>new description</Description></RecomposeVAppParams>`, string(body))
}

func TestAPIRequestLeaseSettingsSection_Body(t *testing.T) {
	deployment, storage := 3600, 7200

	body, err := xml.Marshal(APIRequestLeaseSettingsSection{
		OVF:                      ovfEnvelopeXMLNS,
		XMLNS:                    vcloudXMLNS,
		Info:                     "Lease section settings",
		DeploymentLeaseInSeconds: &deployment,
		StorageLeaseInSeconds:    &storage,
	})
	require.NoError(t, err)

	assert.Equal(t, `<LeaseSettingsSection xmlns:ovf="http://schemas.dmtf.org/ovf/envelope/1" xmlns="http://www.vmware.com/vcloud/v1.5"><ovf:Info>Lease section settings</ovf:Info><DeploymentLeaseInSeconds>3600</DeploymentLeaseInSeconds><StorageLeaseInSeconds>7200</StorageLeaseInSeconds></LeaseSettingsSection>`, string(body))
}
