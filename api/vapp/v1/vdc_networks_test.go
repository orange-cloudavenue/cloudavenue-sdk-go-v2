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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestListVDCNetworkByOwnerID(t *testing.T) {
	client, ms := newClient(t)
	const vdcURN = "urn:vcloud:vdc:12345678-1234-4b8d-89ab-123456789012"

	var gotFilter string
	ms.CleanResponse(endpoints.ListVDCNetwork())
	ms.SetResponseFunc(endpoints.ListVDCNetwork(), func(w http.ResponseWriter, r *http.Request) {
		gotFilter = r.URL.Query().Get("filter")
		w.Header().Set("Content-Type", "application/json,version=39.1")
		require.NoError(t, json.NewEncoder(w).Encode(&itypes.APIResponseListVDCNetwork{
			Values: []itypes.APIResponseVDCNetwork{{ID: "urn:vcloud:network:aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Name: "test org vdc network", NetworkType: types.VDCNetworkTypeRouted}},
		}))
	})

	result, err := client.ListVDCNetwork(t.Context(), types.ParamsListVDCNetwork{OwnerID: vdcURN})
	require.NoError(t, err)
	assert.Equal(t, "ownerRef.id=="+vdcURN, gotFilter)
	if assert.Len(t, result.VDCNetworks, 1) {
		assert.Equal(t, "test org vdc network", result.VDCNetworks[0].Name)
		assert.Equal(t, types.VDCNetworkTypeRouted, result.VDCNetworks[0].NetworkType)
	}
}

func TestListVDCNetworkRequiresOwner(t *testing.T) {
	client, _ := newClient(t)
	_, err := client.ListVDCNetwork(t.Context(), types.ParamsListVDCNetwork{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner id or owner name is required")
}
