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

func TestGetFirewallRequiresVDCGroupOwner(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcID := generator.MustGenerate("{urn:vdc}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcID, Name: "vdc-1"},
	}, nil)

	resp, err := client.GetFirewall(t.Context(), types.ParamsGetEdgeGatewayFirewall{EdgeGatewayID: edgeGatewayID})

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "vdc group")

	ms.CleanResponse(endpoints.GetEdgeGateway())
}

func TestCreateFirewall(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledFalse := false
	enabledTrue := true
	ncpID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.SetResponse(endpoints.UpdateDFWPolicies(), &itypes.APIDFWPolicies{Enabled: true}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled:       true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{ID: "default-policy", Name: "Default", Enabled: &enabledFalse},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDFWDefaultPolicy(), &itypes.APIDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledTrue,
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.SetResponse(endpoints.UpdateDFWRules(), &itypes.APIDistributedFirewallRules{
		Values: []itypes.APIDistributedFirewallRule{{
			ID:                     "rule-1",
			Name:                   "allow-web",
			Description:            "allow web traffic",
			Enabled:                true,
			Direction:              types.FirewallRuleDirectionInOut,
			IPProtocol:             types.FirewallRuleIPProtocolIPv4,
			ActionValue:            types.FirewallRuleActionAllow,
			Logging:                true,
			NetworkContextProfiles: []itypes.APIObjectReference{{ID: ncpID, Name: "ncp-1"}},
		}},
	}, nil)

	resp, err := client.CreateFirewall(t.Context(), types.ParamsCreateEdgeGatewayFirewall{
		EdgeGatewayID: edgeGatewayID,
		Enabled:       true,
		Rules: []types.ParamsFirewallRule{{
			Name:                   "allow-web",
			Description:            "allow web traffic",
			Enabled:                true,
			Direction:              types.FirewallRuleDirectionInOut,
			IPProtocol:             types.FirewallRuleIPProtocolIPv4,
			Action:                 types.FirewallRuleActionAllow,
			Logging:                true,
			NetworkContextProfiles: []types.ParamsFirewallGroupMember{{ID: ncpID, Name: "ncp-1"}},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-1", resp.Rules[0].ID)
	assert.Equal(t, types.FirewallRuleActionAllow, resp.Rules[0].Action)
	assert.Len(t, resp.Rules[0].NetworkContextProfiles, 1)
	assert.Equal(t, ncpID, resp.Rules[0].NetworkContextProfiles[0].ID)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWRules())
}

func TestGetFirewall(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	ncpID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled:       true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{ID: "default-policy", Name: "Default", Enabled: &enabledTrue},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWRules())
	ms.SetResponse(endpoints.GetDFWRules(), &itypes.APIDistributedFirewallRules{
		Values: []itypes.APIDistributedFirewallRule{{
			ID:                     "rule-1",
			Name:                   "allow-web",
			Description:            "allow web traffic",
			Enabled:                true,
			Direction:              types.FirewallRuleDirectionInOut,
			IPProtocol:             types.FirewallRuleIPProtocolIPv4,
			ActionValue:            types.FirewallRuleActionAllow,
			Logging:                true,
			NetworkContextProfiles: []itypes.APIObjectReference{{ID: ncpID, Name: "ncp-1"}},
		}},
	}, nil)

	resp, err := client.GetFirewall(t.Context(), types.ParamsGetEdgeGatewayFirewall{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-1", resp.Rules[0].ID)
	assert.Equal(t, types.FirewallRuleActionAllow, resp.Rules[0].Action)
	assert.Len(t, resp.Rules[0].NetworkContextProfiles, 1)
	assert.Equal(t, ncpID, resp.Rules[0].NetworkContextProfiles[0].ID)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.GetDFWRules())
}

func TestUpdateFirewall(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false
	ncpID := generator.MustGenerate("{uuid}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled:       true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{ID: "default-policy", Name: "Default", Enabled: &enabledTrue},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDFWDefaultPolicy(), &itypes.APIDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledFalse,
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.SetResponse(endpoints.UpdateDFWRules(), &itypes.APIDistributedFirewallRules{
		Values: []itypes.APIDistributedFirewallRule{{
			ID:                     "rule-2",
			Name:                   "deny-web",
			Description:            "deny web traffic",
			Enabled:                true,
			Direction:              types.FirewallRuleDirectionOut,
			IPProtocol:             types.FirewallRuleIPProtocolIPv6,
			ActionValue:            types.FirewallRuleActionDrop,
			NetworkContextProfiles: []itypes.APIObjectReference{{ID: ncpID, Name: "ncp-1"}},
		}},
	}, nil)

	resp, err := client.UpdateFirewall(t.Context(), types.ParamsUpdateEdgeGatewayFirewall{
		EdgeGatewayID: edgeGatewayID,
		Enabled:       false,
		Rules: []types.ParamsFirewallRule{{
			Name:                   "deny-web",
			Description:            "deny web traffic",
			Enabled:                true,
			Direction:              types.FirewallRuleDirectionOut,
			IPProtocol:             types.FirewallRuleIPProtocolIPv6,
			Action:                 types.FirewallRuleActionDrop,
			NetworkContextProfiles: []types.ParamsFirewallGroupMember{{ID: ncpID, Name: "ncp-1"}},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-2", resp.Rules[0].ID)
	assert.Equal(t, types.FirewallRuleActionDrop, resp.Rules[0].Action)
	assert.Len(t, resp.Rules[0].NetworkContextProfiles, 1)
	assert.Equal(t, ncpID, resp.Rules[0].NetworkContextProfiles[0].ID)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWRules())
}

func TestDeleteFirewall(t *testing.T) {
	edgeGatewayID := generator.MustGenerate("{urn:edgegateway}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.SetResponse(endpoints.GetEdgeGateway(), &itypes.APIResponseEdgegateway{
		ID:       edgeGatewayID,
		Name:     "edge-1",
		OwnerRef: &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled:       true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{ID: "default-policy", Name: "Default", Enabled: &enabledTrue},
	}, nil)
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDFWDefaultPolicy(), &itypes.APIDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledFalse,
	}, nil)
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.SetResponse(endpoints.UpdateDFWPolicies(), &itypes.APIDFWPolicies{Enabled: false}, nil)

	err := client.DeleteFirewall(t.Context(), types.ParamsDeleteEdgeGatewayFirewall{EdgeGatewayID: edgeGatewayID})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetEdgeGateway())
	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
}
