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

func TestCreateDynamicSecurityGroup(t *testing.T) {
	createdID := generator.MustGenerate("{urn:firewallGroup}")
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	criteria := []types.ParamsDynamicSecurityGroupCriteria{{
		Rules: []types.ParamsDynamicSecurityGroupCriteriaRule{{
			RuleType: types.DynamicSecurityGroupCriteriaRuleTypeVMTag,
			Operator: types.DynamicSecurityGroupCriteriaRuleOperatorEquals,
			Value:    "web",
		}},
	}}

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.SetResponse(endpoints.CreateFirewallGroup(), &itypes.APIResponseFirewallGroup{ID: createdID}, nil)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{
			ID:   vdcGroupID,
			Name: vdcGroupName,
		}},
	}, nil)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:          createdID,
		Name:        "dsg-1",
		Description: "desc",
		TypeValue:   itypes.FirewallGroupTypeVMCriteria,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		VMCriteria: []itypes.APIFirewallGroupVMCriteria{{
			VMCriteriaRule: []itypes.APIFirewallGroupVMCriteriaRule{{
				AttributeType:  types.DynamicSecurityGroupCriteriaRuleTypeVMTag,
				Operator:       types.DynamicSecurityGroupCriteriaRuleOperatorEquals,
				AttributeValue: "web",
			}},
		}},
	}, nil)

	resp, err := client.CreateDynamicSecurityGroup(t.Context(), types.ParamsCreateDynamicSecurityGroup{
		Name:         "dsg-1",
		Description:  "desc",
		VDCGroupName: vdcGroupName,
		Criteria:     criteria,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, createdID, resp.ID)
	assert.Equal(t, vdcGroupID, resp.OwnerID)
	assert.Equal(t, vdcGroupName, resp.OwnerName)
	assert.Len(t, resp.Criteria, 1)
	assert.Len(t, resp.Criteria[0].Rules, 1)
	assert.Equal(t, types.DynamicSecurityGroupCriteriaRuleTypeVMTag, resp.Criteria[0].Rules[0].RuleType)
	assert.Equal(t, types.DynamicSecurityGroupCriteriaRuleOperatorEquals, resp.Criteria[0].Rules[0].Operator)
	assert.Equal(t, "web", resp.Criteria[0].Rules[0].Value)

	ms.CleanResponse(endpoints.CreateFirewallGroup())
	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.CleanResponse(endpoints.GetFirewallGroup())
}

func TestListDynamicSecurityGroup(t *testing.T) {
	vdcGroupID := generator.MustGenerate("{urn:vdcGroup}")
	vdcGroupName := generator.MustGenerate("{word}")
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.ListVDCGroup())
	ms.SetResponse(endpoints.ListVDCGroup(), &itypes.APIResponseListVDCGroup{
		Values: []itypes.APIResponseListVDCGroupDetails{{ID: vdcGroupID, Name: vdcGroupName}},
	}, nil)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.SetResponse(endpoints.ListFirewallGroup(), &itypes.APIResponseListFirewallGroup{
		Values: []itypes.APIResponseFirewallGroup{{
			ID:        groupID,
			Name:      "dsg-1",
			TypeValue: itypes.FirewallGroupTypeVMCriteria,
			OwnerRef:  &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
			VMCriteria: []itypes.APIFirewallGroupVMCriteria{{
				VMCriteriaRule: []itypes.APIFirewallGroupVMCriteriaRule{{
					AttributeType:  types.DynamicSecurityGroupCriteriaRuleTypeVMTag,
					Operator:       types.DynamicSecurityGroupCriteriaRuleOperatorEquals,
					AttributeValue: "web",
				}},
			}},
		}},
	}, nil)

	resp, err := client.ListDynamicSecurityGroup(t.Context(), types.ParamsListDynamicSecurityGroup{VDCGroupName: vdcGroupName})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.FirewallGroups, 1)
	assert.Equal(t, groupID, resp.FirewallGroups[0].ID)
	assert.Len(t, resp.FirewallGroups[0].Criteria, 1)

	ms.CleanResponse(endpoints.ListFirewallGroup())
	ms.CleanResponse(endpoints.ListVDCGroup())
}

func TestGetDynamicSecurityGroup(t *testing.T) {
	groupID := generator.MustGenerate("{urn:firewallGroup}")

	client, ms := newClient(t)

	ms.CleanResponse(endpoints.GetFirewallGroup())
	ms.SetResponse(endpoints.GetFirewallGroup(), &itypes.APIResponseFirewallGroup{
		ID:        groupID,
		Name:      "dsg-1",
		TypeValue: itypes.FirewallGroupTypeVMCriteria,
		VMCriteria: []itypes.APIFirewallGroupVMCriteria{{
			VMCriteriaRule: []itypes.APIFirewallGroupVMCriteriaRule{{
				AttributeType:  types.DynamicSecurityGroupCriteriaRuleTypeVMName,
				Operator:       types.DynamicSecurityGroupCriteriaRuleOperatorContains,
				AttributeValue: "web",
			}},
		}},
	}, nil)

	resp, err := client.GetDynamicSecurityGroup(t.Context(), types.ParamsGetDynamicSecurityGroup{ID: groupID})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, resp.ID)
	assert.Len(t, resp.Criteria, 1)
	assert.Equal(t, types.DynamicSecurityGroupCriteriaRuleTypeVMName, resp.Criteria[0].Rules[0].RuleType)

	ms.CleanResponse(endpoints.GetFirewallGroup())
}
