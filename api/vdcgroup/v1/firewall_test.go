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

func TestCreateFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledFalse := false
	enabledTrue := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwPolicies())
	ms.SetResponse(endpoints.UpdateDfwPolicies(), &itypes.ApiDfwPolicies{Enabled: true}, nil)

	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.SetResponse(endpoints.GetDfwPolicies(), &itypes.ApiDfwPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.ApiDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledFalse,
		},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDfwDefaultPolicy(), &itypes.ApiDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledTrue,
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwRules())
	ms.SetResponse(endpoints.UpdateDfwRules(), &itypes.ApiDistributedFirewallRules{
		Values: []itypes.ApiDistributedFirewallRule{{
			ID:          "rule-1",
			Name:        "allow-web",
			Description: "allow web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionInOut,
			IpProtocol:  types.FirewallRuleIPProtocolIPv4,
			ActionValue: types.FirewallRuleActionAllow,
			Logging:     true,
		}},
	}, nil)

	resp, err := client.CreateFirewall(t.Context(), types.ParamsCreateFirewall{
		VdcGroupName: vdcGroupName,
		Enabled:      true,
		Rules: []types.ParamsFirewallRule{{
			Name:        "allow-web",
			Description: "allow web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionInOut,
			IPProtocol:  types.FirewallRuleIPProtocolIPv4,
			Action:      types.FirewallRuleActionAllow,
			Logging:     true,
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-1", resp.Rules[0].ID)
	assert.Equal(t, "allow-web", resp.Rules[0].Name)
	assert.Equal(t, types.FirewallRuleActionAllow, resp.Rules[0].Action)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.UpdateDfwPolicies())
	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDfwRules())
}

func TestGetFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.SetResponse(endpoints.GetDfwPolicies(), &itypes.ApiDfwPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.ApiDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
	}, nil)

	ms.CleanResponse(endpoints.GetDfwRules())
	ms.SetResponse(endpoints.GetDfwRules(), &itypes.ApiDistributedFirewallRules{
		Values: []itypes.ApiDistributedFirewallRule{{
			ID:          "rule-1",
			Name:        "allow-web",
			Description: "allow web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionInOut,
			IpProtocol:  types.FirewallRuleIPProtocolIPv4,
			ActionValue: types.FirewallRuleActionAllow,
			Logging:     true,
		}},
	}, nil)

	resp, err := client.GetFirewall(t.Context(), types.ParamsGetFirewall{VdcGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-1", resp.Rules[0].ID)
	assert.Equal(t, types.FirewallRuleActionAllow, resp.Rules[0].Action)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.CleanResponse(endpoints.GetDfwRules())
}

func TestUpdateFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.SetResponse(endpoints.GetDfwPolicies(), &itypes.ApiDfwPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.ApiDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDfwDefaultPolicy(), &itypes.ApiDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledFalse,
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwRules())
	ms.SetResponse(endpoints.UpdateDfwRules(), &itypes.ApiDistributedFirewallRules{
		Values: []itypes.ApiDistributedFirewallRule{{
			ID:          "rule-2",
			Name:        "deny-web",
			Description: "deny web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionOut,
			IpProtocol:  types.FirewallRuleIPProtocolIPv6,
			ActionValue: types.FirewallRuleActionDrop,
			Logging:     false,
		}},
	}, nil)

	resp, err := client.UpdateFirewall(t.Context(), types.ParamsUpdateFirewall{
		VdcGroupName: vdcGroupName,
		Enabled:      false,
		Rules: []types.ParamsFirewallRule{{
			Name:        "deny-web",
			Description: "deny web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionOut,
			IPProtocol:  types.FirewallRuleIPProtocolIPv6,
			Action:      types.FirewallRuleActionDrop,
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-2", resp.Rules[0].ID)
	assert.Equal(t, "deny-web", resp.Rules[0].Name)
	assert.Equal(t, types.FirewallRuleActionDrop, resp.Rules[0].Action)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDfwRules())
}

func TestDeleteFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.SetResponse(endpoints.ListVdcGroup(), &itypes.ApiResponseListVdcGroup{
		Values: []itypes.ApiResponseListVdcGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDfwRules())
	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.SetResponse(endpoints.GetDfwPolicies(), &itypes.ApiDfwPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.ApiDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
	}, nil)
	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDfwDefaultPolicy(), &itypes.ApiDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledFalse,
	}, nil)
	ms.CleanResponse(endpoints.UpdateDfwPolicies())
	ms.SetResponse(endpoints.UpdateDfwPolicies(), &itypes.ApiDfwPolicies{Enabled: false}, nil)

	err := client.DeleteFirewall(t.Context(), types.ParamsDeleteFirewall{VdcGroupName: vdcGroupName})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.ListVdcGroup())
	ms.CleanResponse(endpoints.UpdateDfwRules())
	ms.CleanResponse(endpoints.GetDfwPolicies())
	ms.CleanResponse(endpoints.UpdateDfwDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDfwPolicies())
}
