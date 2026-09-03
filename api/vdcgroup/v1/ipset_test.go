/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestCreateIPSet(t *testing.T) {
	createdID := generator.MustGenerate("{urn:firewallGroup}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	ipAddresses := []string{"192.0.2.10", "192.0.2.11-192.0.2.20"}

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.SetResponse(endpoints.CreateFirewallGroup(), &itypes.ApiResponseFirewallGroup{ID: createdID}, nil)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.ApiResponseFirewallGroup{
		ID:          createdID,
		Name:        "ipset-1",
		Description: "desc",
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		IPAddresses: ipAddresses,
	}, nil)

	resp, err := client.CreateIPSet(t.Context(), types.ParamsCreateIPSet{
		Name:         "ipset-1",
		Description:  "desc",
		VdcGroupName: vdcGroupName,
		IPAddresses:  ipAddresses,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)
	assert.Equal(t, ipAddresses, resp.IPAddresses)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.GetFirewallGroup())
}

func TestListIPSet(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.ApiResponseListFirewallGroup{
		Values: []itypes.ApiResponseFirewallGroup{{
			ID:          groupID,
			Name:        "ipset-1",
			TypeValue:   itypes.FirewallGroupTypeIPSet,
			IPAddresses: []string{"192.0.2.10"},
			OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		}},
	}, nil)

	resp, err := client.ListIPSet(t.Context(), types.ParamsListIPSet{VdcGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.FirewallGroups, 1)
	assert.Equal(t, groupID, resp.FirewallGroups[0].ID)
	assert.Equal(t, []string{"192.0.2.10"}, resp.FirewallGroups[0].IPAddresses)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.CleanResponse(endpoints.ListVdcGroup())
}

func TestGetIPSet(t *testing.T) {
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.ApiResponseFirewallGroup{
		ID:          groupID,
		Name:        "ipset-1",
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		IPAddresses: []string{"192.0.2.10"},
	}, nil)

	resp, err := client.GetIPSet(t.Context(), types.ParamsGetIPSet{ID: groupID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, []string{"192.0.2.10"}, resp.IPAddresses)

	ms.CleanResponse(endpoints.GetFirewallGroup())
}
