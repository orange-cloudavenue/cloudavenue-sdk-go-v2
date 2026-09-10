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

func TestCreateSecurityGroup(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	createdID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.SetResponse(endpoints.CreateFirewallGroup(), &itypes.ApiResponseFirewallGroup{ID: createdID}, nil)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.ApiResponseFirewallGroup{
		ID:          createdID,
		Name:        "sg-1",
		Description: "desc",
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	resp, err := client.CreateSecurityGroup(t.Context(), types.ParamsCreateEdgeGatewaySecurityGroup{
		EdgeGatewayID: edgeGatewayID,
		Name:          "sg-1",
		Description:   "desc",
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.CleanResponse(endpoints.GetFirewallGroup())
}

func TestCreateSecurityGroupRequiresVdcGroupOwner(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcID := generator.MustGenerate("{urn:vdc}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcID, Name: "vdc-1"},
	}, nil)

	resp, err := client.CreateSecurityGroup(t.Context(), types.ParamsCreateEdgeGatewaySecurityGroup{
		EdgeGatewayID: edgeGatewayID,
		Name:          "sg-1",
	})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "vdc group")

	ms.CleanResponse(endpoints.GetEdgeGateway())
}

func TestListSecurityGroup(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.ApiResponseListFirewallGroup{
		Values: []itypes.ApiResponseFirewallGroup{{
			ID:        groupID,
			Name:      "sg-1",
			TypeValue: itypes.FirewallGroupTypeSecurityGroup,
			OwnerRef:  &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		}},
	}, nil)

	resp, err := client.ListSecurityGroup(t.Context(), types.ParamsListEdgeGatewaySecurityGroup{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.FirewallGroups, 1)
	assert.Equal(t, groupID, resp.FirewallGroups[0].ID)
	assert.Equal(t, vdcGroupID, resp.FirewallGroups[0].OwnerID)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListFirewallGroup())
}

func TestGetSecurityGroupByName(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.ApiResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.ApiResponseListFirewallGroup{
		Values: []itypes.ApiResponseFirewallGroup{{
			ID:        groupID,
			Name:      "sg-1",
			TypeValue: itypes.FirewallGroupTypeSecurityGroup,
			OwnerRef:  &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		}},
	}, nil)

	resp, err := client.GetSecurityGroup(t.Context(), types.ParamsGetEdgeGatewaySecurityGroup{
		Name:          "sg-1",
		EdgeGatewayID: edgeGatewayID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, itypes.FirewallGroupTypeSecurityGroup, resp.TypeValue)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.ListFirewallGroup())
}

func TestUpdateSecurityGroup(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.ApiResponseFirewallGroup{
		ID:          groupID,
		Name:        "sg-1",
		Description: "old-desc",
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		Members:     []itypes.ApiObjectReference{{ID: generator.MustGenerate("{uuid}"), Name: "net-a"}},
	}, nil)
	ms.CleanResponse(endpoints.UpdateFirewallGroup())

	resp, err := client.UpdateSecurityGroup(t.Context(), types.ParamsUpdateEdgeGatewaySecurityGroup{
		ID:            groupID,
		EdgeGatewayID: edgeGatewayID,
		Description:   "new-desc",
		Members: []types.ParamsFirewallGroupMember{
			{ID: generator.MustGenerate("{uuid}"), Name: "net-a"},
			{ID: generator.MustGenerate("{uuid}"), Name: "net-b"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Equal(t, "new-desc", resp.Description)
	assert.Len(t, resp.Members, 2)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.UpdateFirewallGroup())
}

func TestDeleteSecurityGroup(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.ApiResponseFirewallGroup{
		ID:        groupID,
		Name:      "sg-1",
		TypeValue: itypes.FirewallGroupTypeSecurityGroup,
		OwnerRef:  &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)
	ms.CleanResponse(endpoints.DeleteFirewallGroup())

	err := client.DeleteSecurityGroup(t.Context(), types.ParamsDeleteEdgeGatewaySecurityGroup{
		ID:            groupID,
		EdgeGatewayID: edgeGatewayID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.DeleteFirewallGroup())
}
