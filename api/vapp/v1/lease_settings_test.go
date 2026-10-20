/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetVAppLeaseSettings(t *testing.T) {
	orgID := "urn:vcloud:org:12345678-1234-4b8d-89ab-123456789012"
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetVAppLeaseSettings())
	ms.SetResponse(endpoints.GetVAppLeaseSettings(), &itypes.APIResponseOrgLeaseSettings{
		DeploymentLeaseInSeconds: 3600,
		StorageLeaseInSeconds:    7200,
	}, nil)

	resp, err := client.GetVAppLeaseSettings(t.Context(), types.ParamsGetVAppLeaseSettings{OrgID: orgID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 3600, resp.DeploymentLeaseInSeconds)
	assert.Equal(t, 7200, resp.StorageLeaseInSeconds)

	ms.CleanResponse(endpoints.GetVAppLeaseSettings())
}

func TestGetVAppLeaseSettings_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.GetVAppLeaseSettings())

	_, err := client.GetVAppLeaseSettings(t.Context(), types.ParamsGetVAppLeaseSettings{})
	assert.Error(t, err)
}

func TestUpdateVAppLeaseSettings(t *testing.T) {
	orgID := "urn:vcloud:org:12345678-1234-4b8d-89ab-123456789012"
	client, ms := newClient(t)

	ms.CleanResponse(endpoints.UpdateVAppLeaseSettings())
	ms.SetResponse(endpoints.UpdateVAppLeaseSettings(), &itypes.APIResponseOrgLeaseSettings{
		DeploymentLeaseInSeconds: 1800,
		StorageLeaseInSeconds:    3600,
	}, nil)

	resp, err := client.UpdateVAppLeaseSettings(t.Context(), types.ParamsUpdateVAppLeaseSettings{
		OrgID: orgID,
		LeaseSettings: types.ModelVAppLeaseSettings{
			DeploymentLeaseInSeconds: 1800,
			StorageLeaseInSeconds:    3600,
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1800, resp.DeploymentLeaseInSeconds)
	assert.Equal(t, 3600, resp.StorageLeaseInSeconds)

	ms.CleanResponse(endpoints.UpdateVAppLeaseSettings())
}

func TestUpdateVAppLeaseSettings_ValidateError(t *testing.T) {
	client, ms := newClient(t)
	defer ms.CleanResponse(endpoints.UpdateVAppLeaseSettings())

	_, err := client.UpdateVAppLeaseSettings(t.Context(), types.ParamsUpdateVAppLeaseSettings{})
	assert.Error(t, err)
}
