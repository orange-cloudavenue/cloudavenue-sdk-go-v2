/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package types

const (
	FirewallRuleDirectionInOut = "IN_OUT"
	FirewallRuleDirectionOut   = "OUT"
	FirewallRuleDirectionIn    = "IN"

	FirewallRuleIPProtocolIPv4     = "IPV4"
	FirewallRuleIPProtocolIPv6     = "IPV6"
	FirewallRuleIPProtocolIPv4IPv6 = "IPV4_IPV6"

	FirewallRuleActionAllow  = "ALLOW"
	FirewallRuleActionDrop   = "DROP"
	FirewallRuleActionReject = "REJECT"
)

var (
	// FirewallRuleDirections lists all valid Direction values for a Firewall rule.
	FirewallRuleDirections = []string{FirewallRuleDirectionInOut, FirewallRuleDirectionOut, FirewallRuleDirectionIn}

	// FirewallRuleIPProtocols lists all valid IPProtocol values for a Firewall rule.
	FirewallRuleIPProtocols = []string{FirewallRuleIPProtocolIPv4, FirewallRuleIPProtocolIPv6, FirewallRuleIPProtocolIPv4IPv6}

	// FirewallRuleActions lists all valid Action values for a Firewall rule.
	FirewallRuleActions = []string{FirewallRuleActionAllow, FirewallRuleActionDrop, FirewallRuleActionReject}
)

type (
	ParamsGetFirewall struct {
		// VDCGroupID is the ID of the Vdc Group owning the Distributed Firewall.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Distributed Firewall.
		VDCGroupName string
	}

	ParamsCreateFirewall struct {
		// VDCGroupID is the ID of the Vdc Group owning the Distributed Firewall.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Distributed Firewall.
		VDCGroupName string

		// Enabled activates or deactivates the Distributed Firewall.
		Enabled bool

		// Rules is the list of Distributed Firewall rules.
		Rules []ParamsFirewallRule
	}

	ParamsUpdateFirewall struct {
		// VDCGroupID is the ID of the Vdc Group owning the Distributed Firewall.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Distributed Firewall.
		VDCGroupName string

		// Enabled activates or deactivates the Distributed Firewall.
		Enabled bool

		// Rules is the list of Distributed Firewall rules.
		Rules []ParamsFirewallRule
	}

	ParamsDeleteFirewall struct {
		// VDCGroupID is the ID of the Vdc Group owning the Distributed Firewall.
		VDCGroupID string

		// VDCGroupName is the name of the Vdc Group owning the Distributed Firewall.
		VDCGroupName string
	}

	ParamsFirewallRule struct {
		// ID is the unique identifier of the rule (output-only on create).
		ID string

		// Name is the name of the rule.
		Name string

		// Description is the description of the rule.
		Description string

		// Enabled activates or deactivates the rule.
		Enabled bool

		// Direction is the traffic direction the rule applies to (IN_OUT, OUT, IN).
		Direction string

		// IPProtocol is the IP protocol the rule applies to (IPV4, IPV6, IPV4_IPV6).
		IPProtocol string

		// Action is the action to apply when the rule matches (ALLOW, DROP, REJECT).
		Action string

		// Logging enables logging of the rule matches.
		Logging bool

		// ApplicationPortProfiles is the list of Application Port Profiles this rule applies to. Empty means Any.
		ApplicationPortProfiles []ParamsFirewallGroupMember

		// SourceFirewallGroups is the list of source Firewall Groups this rule applies to. Empty means Any.
		SourceFirewallGroups []ParamsFirewallGroupMember

		// DestinationFirewallGroups is the list of destination Firewall Groups this rule applies to. Empty means Any.
		DestinationFirewallGroups []ParamsFirewallGroupMember

		// NetworkContextProfiles is the list of Network Context Profiles this rule applies to. Empty means Any.
		NetworkContextProfiles []ParamsFirewallGroupMember

		// SourceGroupsExcluded excludes the source Firewall Groups from the rule instead of including them.
		SourceGroupsExcluded *bool

		// DestinationGroupsExcluded excludes the destination Firewall Groups from the rule instead of including them.
		DestinationGroupsExcluded *bool
	}

	ModelGetFirewall struct {
		Enabled bool                `documentation:"Whether the Distributed Firewall is enabled"`
		Rules   []ModelFirewallRule `documentation:"List of Distributed Firewall rules"`
	}

	ModelFirewallRule struct {
		ID          string `documentation:"ID of the rule"`
		Name        string `documentation:"Name of the rule"`
		Description string `documentation:"Description of the rule"`
		Enabled     bool   `documentation:"Whether the rule is enabled"`
		Direction   string `documentation:"Traffic direction the rule applies to (IN_OUT, OUT, IN)"`
		IPProtocol  string `documentation:"IP protocol the rule applies to (IPV4, IPV6, IPV4_IPV6)"`
		Action      string `documentation:"Action applied when the rule matches (ALLOW, DROP, REJECT)"`
		Logging     bool   `documentation:"Whether logging is enabled for the rule"`

		ApplicationPortProfiles   []ModelGetFirewallGroupMember `documentation:"List of Application Port Profiles this rule applies to"`
		SourceFirewallGroups      []ModelGetFirewallGroupMember `documentation:"List of source Firewall Groups this rule applies to"`
		DestinationFirewallGroups []ModelGetFirewallGroupMember `documentation:"List of destination Firewall Groups this rule applies to"`
		NetworkContextProfiles    []ModelGetFirewallGroupMember `documentation:"List of Network Context Profiles this rule applies to"`

		SourceGroupsExcluded      *bool `documentation:"Whether the source Firewall Groups are excluded from the rule"`
		DestinationGroupsExcluded *bool `documentation:"Whether the destination Firewall Groups are excluded from the rule"`
	}
)
