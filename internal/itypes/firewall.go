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

type (
	// ApiDfwPolicies represents the Distributed Firewall policies of a VDC Group.
	ApiDfwPolicies struct {
		Enabled       bool                 `json:"enabled"`
		DefaultPolicy *ApiDfwDefaultPolicy `json:"defaultPolicy,omitempty"`
	}

	// ApiDfwDefaultPolicy represents the default Distributed Firewall policy of a VDC Group.
	ApiDfwDefaultPolicy struct {
		Description string         `json:"description,omitempty"`
		Enabled     *bool          `json:"enabled,omitempty"`
		ID          string         `json:"id,omitempty"`
		Name        string         `json:"name"`
		Version     *ApiDfwVersion `json:"version,omitempty"`
	}

	// ApiDfwVersion is the optimistic-concurrency version marker used by the DFW default policy.
	ApiDfwVersion struct {
		Version int `json:"version"`
	}

	// ApiDistributedFirewallRules is the bulk container of Distributed Firewall rules.
	ApiDistributedFirewallRules struct {
		Values []ApiDistributedFirewallRule `json:"values"`
	}

	// ApiDistributedFirewallRule represents a single Distributed Firewall rule.
	ApiDistributedFirewallRule struct {
		ID                        string               `json:"id,omitempty"`
		Name                      string               `json:"name"`
		Description               string               `json:"description,omitempty"`
		Comments                  string               `json:"comments,omitempty"`
		ApplicationPortProfiles   []ApiObjectReference `json:"applicationPortProfiles,omitempty"`
		SourceFirewallGroups      []ApiObjectReference `json:"sourceFirewallGroups,omitempty"`
		DestinationFirewallGroups []ApiObjectReference `json:"destinationFirewallGroups,omitempty"`
		NetworkContextProfiles    []ApiObjectReference `json:"networkContextProfiles,omitempty"`
		Direction                 string               `json:"direction"`
		Enabled                   bool                 `json:"enabled"`
		IpProtocol                string               `json:"ipProtocol"`
		Logging                   bool                 `json:"logging"`
		ActionValue               string               `json:"actionValue,omitempty"`
		SourceGroupsExcluded      *bool                `json:"sourceGroupsExcluded,omitempty"`
		DestinationGroupsExcluded *bool                `json:"destinationGroupsExcluded,omitempty"`
	}
)

// ToModel converts the wire DFW rule into its public model representation.
func (r *ApiDistributedFirewallRule) ToModel() types.ModelFirewallRule {
	m := types.ModelFirewallRule{
		ID:                        r.ID,
		Name:                      r.Name,
		Description:               r.Description,
		Direction:                 r.Direction,
		Enabled:                   r.Enabled,
		IPProtocol:                r.IpProtocol,
		Action:                    r.ActionValue,
		Logging:                   r.Logging,
		SourceGroupsExcluded:      r.SourceGroupsExcluded,
		DestinationGroupsExcluded: r.DestinationGroupsExcluded,
	}

	for _, ref := range r.ApplicationPortProfiles {
		m.ApplicationPortProfiles = append(m.ApplicationPortProfiles, types.ModelGetFirewallGroupMember{ID: ref.ID, Name: ref.Name})
	}
	for _, ref := range r.SourceFirewallGroups {
		m.SourceFirewallGroups = append(m.SourceFirewallGroups, types.ModelGetFirewallGroupMember{ID: ref.ID, Name: ref.Name})
	}
	for _, ref := range r.DestinationFirewallGroups {
		m.DestinationFirewallGroups = append(m.DestinationFirewallGroups, types.ModelGetFirewallGroupMember{ID: ref.ID, Name: ref.Name})
	}
	for _, ref := range r.NetworkContextProfiles {
		m.NetworkContextProfiles = append(m.NetworkContextProfiles, types.ModelGetFirewallGroupMember{ID: ref.ID, Name: ref.Name})
	}

	return m
}

// ToModel converts the wire DFW rules container into its public model representation.
func (r *ApiDistributedFirewallRules) ToModel() []types.ModelFirewallRule {
	rules := make([]types.ModelFirewallRule, 0, len(r.Values))
	for _, v := range r.Values {
		rules = append(rules, v.ToModel())
	}
	return rules
}
