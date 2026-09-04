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
	"reflect"
	"testing"

	"github.com/orange-cloudavenue/common-go/generator"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav/mock"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func newFirewallGroupTestClient(t *testing.T) (*Client, *mock.Server) {
	t.Helper()

	results := reflect.ValueOf(newClient).Call([]reflect.Value{reflect.ValueOf(t)})
	client, _ := results[0].Interface().(*Client)

	if len(results) < 2 || results[1].IsNil() {
		return client, nil
	}

	ms, _ := results[1].Interface().(*mock.Server)
	return client, ms
}

func TestUpdateSecurityGroup(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          fwGroupID,
		Name:        "sg-1",
		Description: "old-desc",
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		Members: []itypes.APIObjectReference{
			{ID: generator.MustGenerate("{uuid}"), Name: "net-a"},
		},
	}, nil)

	ms.CleanResponse(endpoints.UpdateFirewallGroup())

	resp, err := client.UpdateSecurityGroup(t.Context(), types.ParamsUpdateSecurityGroup{
		ID:          fwGroupID,
		Description: "new-desc",
		Members: []types.ParamsFirewallGroupMember{
			{ID: generator.MustGenerate("{uuid}"), Name: "net-a"},
			{ID: generator.MustGenerate("{uuid}"), Name: "net-b"},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fwGroupID, resp.ID)
	assert.Equal(t, "new-desc", resp.Description)
	assert.Len(t, resp.Members, 2)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.UpdateFirewallGroup())
}

func TestDeleteSecurityGroup(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:        fwGroupID,
		Name:      "sg-1",
		TypeValue: itypes.FirewallGroupTypeSecurityGroup,
	}, nil)

	ms.CleanResponse(endpoints.DeleteFirewallGroup())

	err := client.DeleteSecurityGroup(t.Context(), types.ParamsDeleteSecurityGroup{
		ID: fwGroupID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.DeleteFirewallGroup())
}

func TestUpdateIPSet(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          fwGroupID,
		Name:        "ipset-1",
		Description: "old",
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		IPAddresses: []string{"10.0.0.1"},
	}, nil)

	ms.CleanResponse(endpoints.UpdateFirewallGroup())

	resp, err := client.UpdateIPSet(t.Context(), types.ParamsUpdateIPSet{
		ID:          fwGroupID,
		Description: "updated",
		IPAddresses: []string{"10.0.0.1", "10.0.0.2"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fwGroupID, resp.ID)
	assert.Equal(t, "updated", resp.Description)
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, resp.IPAddresses)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.UpdateFirewallGroup())
}

func TestDeleteIPSet(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:        fwGroupID,
		Name:      "ipset-1",
		TypeValue: itypes.FirewallGroupTypeIPSet,
	}, nil)

	ms.CleanResponse(endpoints.DeleteFirewallGroup())

	err := client.DeleteIPSet(t.Context(), types.ParamsDeleteIPSet{
		ID: fwGroupID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.DeleteFirewallGroup())
}

func TestUpdateDynamicSecurityGroup(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          fwGroupID,
		Name:        "dsg-1",
		Description: "old",
		TypeValue:   itypes.FirewallGroupTypeVMCriteria,
		VMCriteria: []itypes.APIFirewallGroupVMCriteria{{
			VMCriteriaRule: []itypes.APIFirewallGroupVMCriteriaRule{{
				AttributeType:  types.DynamicSecurityGroupCriteriaRuleTypeVMTag,
				Operator:       types.DynamicSecurityGroupCriteriaRuleOperatorEquals,
				AttributeValue: "old-tag",
			}},
		}},
	}, nil)

	ms.CleanResponse(endpoints.UpdateFirewallGroup())

	resp, err := client.UpdateDynamicSecurityGroup(t.Context(), types.ParamsUpdateDynamicSecurityGroup{
		ID:          fwGroupID,
		Description: "updated",
		Criteria: []types.ParamsDynamicSecurityGroupCriteria{{
			Rules: []types.ParamsDynamicSecurityGroupCriteriaRule{{
				RuleType: types.DynamicSecurityGroupCriteriaRuleTypeVMName,
				Operator: types.DynamicSecurityGroupCriteriaRuleOperatorContains,
				Value:    "web",
			}},
		}},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fwGroupID, resp.ID)
	assert.Equal(t, "updated", resp.Description)
	assert.Len(t, resp.Criteria, 1)
	assert.Equal(t, types.DynamicSecurityGroupCriteriaRuleTypeVMName, resp.Criteria[0].Rules[0].RuleType)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.UpdateFirewallGroup())
}

func TestDeleteDynamicSecurityGroup(t *testing.T) {
	fwGroupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newFirewallGroupTestClient(t)
	if ms == nil {
		t.Skip("mock server not available")
	}

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:        fwGroupID,
		Name:      "dsg-1",
		TypeValue: itypes.FirewallGroupTypeVMCriteria,
	}, nil)

	ms.CleanResponse(endpoints.DeleteFirewallGroup())

	err := client.DeleteDynamicSecurityGroup(t.Context(), types.ParamsDeleteDynamicSecurityGroup{
		ID: fwGroupID,
	})

	assert.NoError(t, err)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.CleanResponse(endpoints.DeleteFirewallGroup())
}
