/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestGetVDCNetworkByID(t *testing.T) {
	client, ms := newClient(t)
	networkID := generator.MustGenerate("{urn:network}")

	ms.CleanResponse(endpoints.GetVDCNetwork())
	ms.SetResponse(endpoints.GetVDCNetwork(), &itypes.APIResponseVDCNetwork{
		ID:          networkID,
		Name:        "org-vdc-network",
		NetworkType: types.VDCNetworkTypeRouted,
	}, nil)

	resp, err := client.GetVDCNetwork(t.Context(), types.ParamsGetVDCNetwork{ID: networkID})
	require.NoError(t, err)
	assert.Equal(t, networkID, resp.ID)
	// cloudapi carries no href; the SDK derives the tenant-accessible legacy
	// network URL used by vApp runtime ParentNetwork references.
	assert.Contains(t, resp.HREF, "/api/network/"+networkID[len("urn:vcloud:network:"):])
	ref := resp.VAppParentNetworkReference()
	if assert.NotNil(t, ref) {
		assert.Equal(t, resp.HREF, ref.HREF)
		assert.Equal(t, networkID, ref.ID)
		assert.Equal(t, "org-vdc-network", ref.Name)
		assert.Empty(t, ref.Type)
	}

	ms.CleanResponse(endpoints.GetVDCNetwork())
}

func TestGetVDCNetworkByNamePaginates(t *testing.T) {
	client, ms := newClient(t)
	vdcURN := generator.MustGenerate("{urn:vdc}")

	var filters []string
	var pages []string
	ms.CleanResponse(endpoints.ListVDCNetwork())
	ms.SetResponseFunc(endpoints.ListVDCNetwork(), func(w http.ResponseWriter, r *http.Request) {
		filters = append(filters, r.URL.Query().Get("filter"))
		pages = append(pages, r.URL.Query().Get("page"))

		list := &itypes.APIResponseListVDCNetwork{ResultTotal: 32, PageCount: 2, Page: 1}
		if r.URL.Query().Get("page") == "2" {
			list.Page = 2
		} else {
			list.Values = []itypes.APIResponseVDCNetwork{{ID: "urn:vcloud:network:aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", Name: "other", NetworkType: types.VDCNetworkTypeRouted}}
		}
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(list))
	})

	resp, err := client.GetVDCNetwork(t.Context(), types.ParamsGetVDCNetwork{Name: "target", OwnerID: vdcURN})
	require.NoError(t, err)
	assert.Equal(t, "urn:vcloud:network:aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", resp.ID)
	assert.Contains(t, resp.HREF, "/api/network/aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa")
	assert.Equal(t, []string{
		"name==target;ownerRef.id==" + vdcURN,
		"name==target;ownerRef.id==" + vdcURN,
	}, filters)
	assert.Equal(t, []string{"1", "2"}, pages)

	ms.CleanResponse(endpoints.ListVDCNetwork())
}

func TestGetVDCNetworkByNameRequiresOwner(t *testing.T) {
	client, _ := newClient(t)
	_, err := client.GetVDCNetwork(t.Context(), types.ParamsGetVDCNetwork{Name: "org-vdc-network"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner id or owner name is required")
}
