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
	// APIDFWPolicies represents the Distributed Firewall policies of a VDC Group.
	APIDFWPolicies struct {
		Enabled       bool                 `json:"enabled"`
		DefaultPolicy *APIDfwDefaultPolicy `json:"defaultPolicy,omitempty"`
	}

	// APIDfwDefaultPolicy represents the default Distributed Firewall policy of a VDC Group.
	APIDfwDefaultPolicy struct {
		Description string         `json:"description,omitempty"`
		Enabled     *bool          `json:"enabled,omitempty"`
		ID          string         `json:"id,omitempty"`
		Name        string         `json:"name"`
		Version     *APIDfwVersion `json:"version,omitempty"`
	}

	// APIDfwVersion is the optimistic-concurrency version marker used by the DFW default policy.
	APIDfwVersion struct {
		Version int `json:"version"`
	}

	// APIDistributedFirewallRules is the bulk container of Distributed Firewall rules.
	APIDistributedFirewallRules struct {
		Values []APIDistributedFirewallRule `json:"values"`
	}

	// APIDistributedFirewallRule represents a single Distributed Firewall rule.
	APIDistributedFirewallRule struct {
		ID                        string               `json:"id,omitempty"`
		Name                      string               `json:"name"`
		Description               string               `json:"description,omitempty"`
		Comments                  string               `json:"comments,omitempty"`
		ApplicationPortProfiles   []APIObjectReference `json:"applicationPortProfiles,omitempty"`
		SourceFirewallGroups      []APIObjectReference `json:"sourceFirewallGroups,omitempty"`
		DestinationFirewallGroups []APIObjectReference `json:"destinationFirewallGroups,omitempty"`
		NetworkContextProfiles    []APIObjectReference `json:"networkContextProfiles,omitempty"`
		Direction                 string               `json:"direction"`
		Enabled                   bool                 `json:"enabled"`
		IPProtocol                string               `json:"ipProtocol"`
		Logging                   bool                 `json:"logging"`
		ActionValue               string               `json:"actionValue,omitempty"`
		SourceGroupsExcluded      *bool                `json:"sourceGroupsExcluded,omitempty"`
		DestinationGroupsExcluded *bool                `json:"destinationGroupsExcluded,omitempty"`
	}
)

// ToModel converts the wire DFW rule into its public model representation.
func (r *APIDistributedFirewallRule) ToModel() types.ModelFirewallRule {
	m := types.ModelFirewallRule{
		ID:                        r.ID,
		Name:                      r.Name,
		Description:               r.Description,
		Direction:                 r.Direction,
		Enabled:                   r.Enabled,
		IPProtocol:                r.IPProtocol,
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
func (r *APIDistributedFirewallRules) ToModel() []types.ModelFirewallRule {
	rules := make([]types.ModelFirewallRule, 0, len(r.Values))
	for _, v := range r.Values {
		rules = append(rules, v.ToModel())
	}
	return rules
}
