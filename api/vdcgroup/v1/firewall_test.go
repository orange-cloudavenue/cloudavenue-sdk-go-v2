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

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.SetResponse(endpoints.UpdateDFWPolicies(), &itypes.APIDFWPolicies{Enabled: true}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledFalse,
		},
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
			ID:          "rule-1",
			Name:        "allow-web",
			Description: "allow web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionInOut,
			IPProtocol:  types.FirewallRuleIPProtocolIPv4,
			ActionValue: types.FirewallRuleActionAllow,
			Logging:     true,
		}},
	}, nil)

	resp, err := client.CreateFirewall(t.Context(), types.ParamsCreateFirewall{
		VDCGroupName: vdcGroupName,
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

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWRules())
}

func TestGetFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWRules())
	ms.SetResponse(endpoints.GetDFWRules(), &itypes.APIDistributedFirewallRules{
		Values: []itypes.APIDistributedFirewallRule{{
			ID:          "rule-1",
			Name:        "allow-web",
			Description: "allow web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionInOut,
			IPProtocol:  types.FirewallRuleIPProtocolIPv4,
			ActionValue: types.FirewallRuleActionAllow,
			Logging:     true,
		}},
	}, nil)

	resp, err := client.GetFirewall(t.Context(), types.ParamsGetFirewall{VDCGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Enabled)
	assert.Len(t, resp.Rules, 1)
	assert.Equal(t, "rule-1", resp.Rules[0].ID)
	assert.Equal(t, types.FirewallRuleActionAllow, resp.Rules[0].Action)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.GetDFWRules())
}

func TestUpdateFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
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
			ID:          "rule-2",
			Name:        "deny-web",
			Description: "deny web traffic",
			Enabled:     true,
			Direction:   types.FirewallRuleDirectionOut,
			IPProtocol:  types.FirewallRuleIPProtocolIPv6,
			ActionValue: types.FirewallRuleActionDrop,
			Logging:     false,
		}},
	}, nil)

	resp, err := client.UpdateFirewall(t.Context(), types.ParamsUpdateFirewall{
		VDCGroupName: vdcGroupName,
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

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWRules())
}

func TestDeleteFirewall(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	enabledTrue := true
	enabledFalse := false

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.SetResponse(endpoints.GetDFWPolicies(), &itypes.APIDFWPolicies{
		Enabled: true,
		DefaultPolicy: &itypes.APIDfwDefaultPolicy{
			ID:      "default-policy",
			Name:    "Default",
			Enabled: &enabledTrue,
		},
	}, nil)
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.SetResponse(endpoints.UpdateDFWDefaultPolicy(), &itypes.APIDfwDefaultPolicy{
		ID:      "default-policy",
		Name:    "Default",
		Enabled: &enabledFalse,
	}, nil)
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
	ms.SetResponse(endpoints.UpdateDFWPolicies(), &itypes.APIDFWPolicies{Enabled: false}, nil)

	err := client.DeleteFirewall(t.Context(), types.ParamsDeleteFirewall{VDCGroupName: vdcGroupName})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.UpdateDFWRules())
	ms.CleanResponse(endpoints.GetDFWPolicies())
	ms.CleanResponse(endpoints.UpdateDFWDefaultPolicy())
	ms.CleanResponse(endpoints.UpdateDFWPolicies())
}
