/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"

const (
	FirewallGroupTypeSecurityGroup = "SECURITY_GROUP"
	FirewallGroupTypeIPSet         = "IP_SET"
	FirewallGroupTypeVMCriteria    = "VM_CRITERIA"
)

type (
	// APIResponseListFirewallGroup is shared list payload for firewall group resources.
	APIResponseListFirewallGroup struct {
		Values []APIResponseFirewallGroup `json:"values" fakesize:"3"`
	}

	// APIResponseFirewallGroup is shared get/create payload for firewall group resources.
	APIResponseFirewallGroup struct {
		ID          string                       `json:"id,omitempty" fake:"{urn:firewallGroup}"`
		Name        string                       `json:"name" fake:"mockfirewallgroup-{word}"`
		Description string                       `json:"description,omitempty" fake:"{sentence}"`
		IPAddresses []string                     `json:"ipAddresses,omitempty"`
		Members     []APIObjectReference         `json:"members,omitempty"`
		VMCriteria  []APIFirewallGroupVMCriteria `json:"vmCriteria,omitempty"`
		OwnerRef    *APIObjectReference          `json:"ownerRef,omitempty"`
		TypeValue   string                       `json:"typeValue,omitempty"`
	}

	APIFirewallGroupVMCriteria struct {
		VMCriteriaRule []APIFirewallGroupVMCriteriaRule `json:"rules,omitempty"`
	}

	APIFirewallGroupVMCriteriaRule struct {
		AttributeType  string `json:"attributeType,omitempty"`
		AttributeValue string `json:"attributeValue,omitempty"`
		Operator       string `json:"operator,omitempty"`
	}

	// APIRequestFirewallGroup is shared create/update request payload for firewall group resources.
	APIRequestFirewallGroup struct {
		ID          string                       `json:"id,omitempty" fake:"{urn:firewallGroup}"`
		Name        string                       `json:"name" fake:"mockfirewallgroup-{word}"`
		Description string                       `json:"description,omitempty" fake:"{sentence}"`
		IPAddresses []string                     `json:"ipAddresses,omitempty"`
		Members     []APIObjectReference         `json:"members,omitempty"`
		VMCriteria  []APIFirewallGroupVMCriteria `json:"vmCriteria,omitempty"`
		OwnerRef    *APIObjectReference          `json:"ownerRef,omitempty"`
		TypeValue   string                       `json:"typeValue,omitempty"`
	}
)

func (r *APIRequestFirewallGroup) ToModel() types.ModelGetFirewallGroup {
	m := types.ModelGetFirewallGroup{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		IPAddresses: r.IPAddresses,
		TypeValue:   r.TypeValue,
	}

	for _, member := range r.Members {
		m.Members = append(m.Members, types.ModelGetFirewallGroupMember{
			ID:   member.ID,
			Name: member.Name,
		})
	}

	if r.OwnerRef != nil {
		m.OwnerID = r.OwnerRef.ID
		m.OwnerName = r.OwnerRef.Name
	}

	for _, criteria := range r.VMCriteria {
		c := types.ModelGetFirewallGroupVMCriteria{}
		for _, rule := range criteria.VMCriteriaRule {
			c.Rules = append(c.Rules, types.ModelGetFirewallGroupVMCriteriaRule{
				RuleType: rule.AttributeType,
				Operator: rule.Operator,
				Value:    rule.AttributeValue,
			})
		}
		m.Criteria = append(m.Criteria, c)
	}

	return m
}

func (r *APIResponseListFirewallGroup) ToModel() *types.ModelListFirewallGroup {
	model := &types.ModelListFirewallGroup{
		FirewallGroups: make([]types.ModelGetFirewallGroup, 0),
	}

	for _, fwGroup := range r.Values {
		model.FirewallGroups = append(model.FirewallGroups, fwGroup.ToModel())
	}

	return model
}

func (r *APIResponseFirewallGroup) ToModel() types.ModelGetFirewallGroup {
	m := types.ModelGetFirewallGroup{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		IPAddresses: r.IPAddresses,
		TypeValue:   r.TypeValue,
	}

	for _, member := range r.Members {
		m.Members = append(m.Members, types.ModelGetFirewallGroupMember{
			ID:   member.ID,
			Name: member.Name,
		})
	}

	if r.OwnerRef != nil {
		m.OwnerID = r.OwnerRef.ID
		m.OwnerName = r.OwnerRef.Name
	}

	for _, criteria := range r.VMCriteria {
		c := types.ModelGetFirewallGroupVMCriteria{}
		for _, rule := range criteria.VMCriteriaRule {
			c.Rules = append(c.Rules, types.ModelGetFirewallGroupVMCriteriaRule{
				RuleType: rule.AttributeType,
				Operator: rule.Operator,
				Value:    rule.AttributeValue,
			})
		}
		m.Criteria = append(m.Criteria, c)
	}

	return m
}
