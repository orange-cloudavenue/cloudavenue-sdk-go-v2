/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func TestCreateIPSet(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	createdID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.SetResponse(endpoints.CreateFirewallGroup(), &itypes.APIResponseFirewallGroup{ID: createdID}, nil)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          createdID,
		Name:        "ipset-1",
		Description: "desc",
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		IPAddresses: []string{"10.0.0.1", "10.0.0.0/24"},
	}, nil)

	resp, err := client.CreateIPSet(t.Context(), types.ParamsCreateEdgeGatewayIPSet{
		EdgeGatewayID: edgeGatewayID,
		Name:          "ipset-1",
		Description:   "desc",
		IPAddresses:   []string{"10.0.0.1", "10.0.0.0/24"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.0/24"}, resp.IPAddresses)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.CleanResponse(endpoints.GetFirewallGroup())
}

func TestCreateIPSetRequiresVDCGroupOwner(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcID := generator.MustGenerate("{urn:vdc}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcID, Name: "vdc-1"},
	}, nil)

	resp, err := client.CreateIPSet(t.Context(), types.ParamsCreateEdgeGatewayIPSet{
		EdgeGatewayID: edgeGatewayID,
		Name:          "ipset-1",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "vdc group")

	ms.CleanResponse(endpoints.GetEdgeGateway())
}

func TestListIPSet(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.APIResponseListFirewallGroup{
		Values: []itypes.APIResponseFirewallGroup{{
			ID:          groupID,
			Name:        "ipset-1",
			TypeValue:   itypes.FirewallGroupTypeIPSet,
			OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
			IPAddresses: []string{"10.0.0.1"},
		}},
	}, nil)

	resp, err := client.ListIPSet(t.Context(), types.ParamsListEdgeGatewayIPSet{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.FirewallGroups, 1)
	assert.Equal(t, groupID, resp.FirewallGroups[0].ID)
	assert.Equal(t, []string{"10.0.0.1"}, resp.FirewallGroups[0].IPAddresses)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListFirewallGroup())
}

func TestGetIPSetByName(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.APIResponseListFirewallGroup{
		Values: []itypes.APIResponseFirewallGroup{{
			ID:          groupID,
			Name:        "ipset-1",
			TypeValue:   itypes.FirewallGroupTypeIPSet,
			OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
			IPAddresses: []string{"10.0.0.1"},
		}},
	}, nil)

	resp, err := client.GetIPSet(t.Context(), types.ParamsGetEdgeGatewayIPSet{
		Name:          "ipset-1",
		EdgeGatewayID: edgeGatewayID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, itypes.FirewallGroupTypeIPSet, resp.TypeValue)
	assert.Equal(t, []string{"10.0.0.1"}, resp.IPAddresses)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListFirewallGroup())
}

func TestUpdateIPSet(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          groupID,
		Name:        "ipset-1",
		Description: "old",
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		IPAddresses: []string{"10.0.0.1"},
	}, nil)
	ms.CleanResponse(endpoints.UpdateFirewallGroup())

	resp, err := client.UpdateIPSet(t.Context(), types.ParamsUpdateEdgeGatewayIPSet{
		ID:            groupID,
		EdgeGatewayID: edgeGatewayID,
		Description:   "updated",
		IPAddresses:   []string{"10.0.0.1", "10.0.0.2"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, "updated", resp.Description)
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, resp.IPAddresses)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.UpdateFirewallGroup())
}

func TestDeleteIPSet(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:        groupID,
		Name:      "ipset-1",
		TypeValue: itypes.FirewallGroupTypeIPSet,
		OwnerRef:  &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)
	ms.CleanResponse(endpoints.DeleteFirewallGroup())

	err := client.DeleteIPSet(t.Context(), types.ParamsDeleteEdgeGatewayIPSet{
		ID:            groupID,
		EdgeGatewayID: edgeGatewayID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.DeleteFirewallGroup())
}
