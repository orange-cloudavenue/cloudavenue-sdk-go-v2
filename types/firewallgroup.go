/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

type (
	// ModelListFirewallGroup is list shape shared by firewall group resources.
	ModelListFirewallGroup struct {
		FirewallGroups []ModelGetFirewallGroup `documentation:"List of Firewall Groups"`
	}

	// ModelGetFirewallGroup is shared shape for Security Group, IP Set, and Dynamic Security Group.
	ModelGetFirewallGroup struct {
		ID          string `documentation:"ID of the Firewall Group"`
		Name        string `documentation:"Name of the Firewall Group"`
		Description string `documentation:"Description of the Firewall Group"`
		TypeValue   string `documentation:"Type of the Firewall Group (SECURITY_GROUP, IP_SET, VM_CRITERIA)"`

		OwnerID   string `documentation:"ID of the owner (EdgeGateway or VdcGroup) of the Firewall Group"`
		OwnerName string `documentation:"Name of the owner (EdgeGateway or VdcGroup) of the Firewall Group"`

		IPAddresses []string                          `documentation:"List of IP addresses (only for IP_SET Firewall Groups)"`
		Members     []ModelGetFirewallGroupMember     `documentation:"List of members (only for SECURITY_GROUP Firewall Groups)"`
		Criteria    []ModelGetFirewallGroupVMCriteria `documentation:"List of criteria (only for VM_CRITERIA Firewall Groups)"`
	}

	ModelGetFirewallGroupMember struct {
		ID   string `documentation:"ID of the member"`
		Name string `documentation:"Name of the member"`
	}

	ModelGetFirewallGroupVMCriteria struct {
		Rules []ModelGetFirewallGroupVMCriteriaRule `documentation:"List of rules for this criteria"`
	}

	ModelGetFirewallGroupVMCriteriaRule struct {
		RuleType string `documentation:"Type of the rule (VM_NAME or VM_TAG)"`
		Operator string `documentation:"Operator of the rule (EQUALS, CONTAINS, STARTS_WITH, ENDS_WITH)"`
		Value    string `documentation:"Value to match against"`
	}
)
