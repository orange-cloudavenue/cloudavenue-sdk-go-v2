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
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetFirewall    = "Firewall.Get"
	opCreateFirewall = "Firewall.Create"
	opUpdateFirewall = "Firewall.Update"
	opDeleteFirewall = "Firewall.Delete"
)

// toAPIConfiguredFirewallRules converts public firewall rules to wire format.
//
// Comments intentionally mirror Description for v1 compatibility, and ActionValue is always used.
func toApiDistributedFirewallRules(rules []types.ParamsFirewallRule) []itypes.ApiDistributedFirewallRule {
	out := make([]itypes.ApiDistributedFirewallRule, 0, len(rules))

	for _, r := range rules {
		out = append(out, itypes.ApiDistributedFirewallRule{
			ID:                        r.ID,
			Name:                      r.Name,
			Description:               r.Description,
			Comments:                  r.Description,
			ApplicationPortProfiles:   toApiObjectReferences(r.ApplicationPortProfiles),
			SourceFirewallGroups:      toApiObjectReferences(r.SourceFirewallGroups),
			DestinationFirewallGroups: toApiObjectReferences(r.DestinationFirewallGroups),
			NetworkContextProfiles:    toApiObjectReferences(r.NetworkContextProfiles),
			Direction:                 r.Direction,
			Enabled:                   r.Enabled,
			IpProtocol:                r.IPProtocol,
			Logging:                   r.Logging,
			ActionValue:               r.Action,
			SourceGroupsExcluded:      r.SourceGroupsExcluded,
			DestinationGroupsExcluded: r.DestinationGroupsExcluded,
		})
	}

	return out
}

// toApiObjectReferences converts firewall group members to wire references.
func toApiObjectReferences(members []types.ParamsFirewallGroupMember) []itypes.ApiObjectReference {
	out := make([]itypes.ApiObjectReference, 0, len(members))
	for _, m := range members {
		out = append(out, itypes.ApiObjectReference{ID: m.ID, Name: m.Name})
	}
	return out
}

// dfwEnableOrDisable updates default-policy enabled state only when needed.
func dfwEnableOrDisable(ctx context.Context, cc *Client, vdcGroupID string, enable bool) error {
	epGet := endpoints.GetDfwPolicies()

	resp, err := cc.c.Do(ctx, epGet, cav.WithPathParam(epGet.PathParams[0], vdcGroupID))
	if err != nil {
		return err
	}

	policies := resp.Result().(*itypes.ApiDfwPolicies)

	isEnabled := policies.DefaultPolicy != nil && policies.DefaultPolicy.Enabled != nil && *policies.DefaultPolicy.Enabled
	if isEnabled == enable {
		return nil
	}

	defaultPolicy := itypes.ApiDfwDefaultPolicy{
		Enabled: &enable,
	}
	if policies.DefaultPolicy != nil {
		defaultPolicy.Description = policies.DefaultPolicy.Description
		defaultPolicy.ID = policies.DefaultPolicy.ID
		defaultPolicy.Name = policies.DefaultPolicy.Name
		defaultPolicy.Version = policies.DefaultPolicy.Version
	}
	if defaultPolicy.Name == "" {
		defaultPolicy.Name = "Default"
	}

	epUpdate := endpoints.UpdateDfwDefaultPolicy()

	_, err = cc.c.Do(
		ctx,
		epUpdate,
		cav.WithPathParam(epUpdate.PathParams[0], vdcGroupID),
		cav.SetBody(defaultPolicy),
	)

	return err
}

func firewallRulesBody(rules []types.ParamsFirewallRule) itypes.ApiDistributedFirewallRules {
	return itypes.ApiDistributedFirewallRules{Values: toApiDistributedFirewallRules(rules)}
}

// GetFirewall returns distributed firewall state for a VDC group.
func (c *Client) GetFirewall(ctx context.Context, params types.ParamsGetFirewall) (*types.ModelGetFirewall, error) {
	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opGetFirewall, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	epPolicies := endpoints.GetDfwPolicies()
	respPolicies, err := c.c.Do(ctx, epPolicies, cav.WithPathParam(epPolicies.PathParams[0], vdcGroupID))
	if err != nil {
		return nil, fmt.Errorf("%s: get policies: %w", opGetFirewall, err)
	}

	policies := respPolicies.Result().(*itypes.ApiDfwPolicies)
	enabled := policies.DefaultPolicy != nil && policies.DefaultPolicy.Enabled != nil && *policies.DefaultPolicy.Enabled

	epRules := endpoints.GetDfwRules()
	respRules, err := c.c.Do(ctx, epRules, cav.WithPathParam(epRules.PathParams[0], vdcGroupID))
	if err != nil {
		return nil, fmt.Errorf("%s: get rules: %w", opGetFirewall, err)
	}

	rules := respRules.Result().(*itypes.ApiDistributedFirewallRules)

	return &types.ModelGetFirewall{
		Enabled: enabled,
		Rules:   rules.ToModel(),
	}, nil
}

// CreateFirewall activates distributed firewalling and sets initial rules.
func (c *Client) CreateFirewall(ctx context.Context, params types.ParamsCreateFirewall) (*types.ModelGetFirewall, error) {
	vdcGroupID, _, err := resolveVdcGroupRef(ctx, c.c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdc group: %w", opCreateFirewall, err)
	}

	epPolicies := endpoints.UpdateDfwPolicies()
	if _, err := c.c.Do(
		ctx,
		epPolicies,
		cav.WithPathParam(epPolicies.PathParams[0], vdcGroupID),
		cav.SetBody(itypes.ApiDfwPolicies{Enabled: true}),
	); err != nil {
		return nil, fmt.Errorf("%s: activate: %w", opCreateFirewall, err)
	}

	if err := dfwEnableOrDisable(ctx, c, vdcGroupID, params.Enabled); err != nil {
		return nil, fmt.Errorf("%s: set enabled state: %w", opCreateFirewall, err)
	}

	epRules := endpoints.UpdateDfwRules()
	respRules, err := c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], vdcGroupID),
		cav.SetBody(firewallRulesBody(params.Rules)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: update rules: %w", opCreateFirewall, err)
	}

	rules, ok := respRules.Result().(*itypes.ApiDistributedFirewallRules)
	if !ok || rules == nil {
		return nil, fmt.Errorf("%s: unexpected rules response type %T", opCreateFirewall, respRules.Result())
	}

	return &types.ModelGetFirewall{
		Enabled: params.Enabled,
		Rules:   rules.ToModel(),
	}, nil
}

// UpdateFirewall replaces distributed firewall rules and enabled state.
func (c *Client) UpdateFirewall(ctx context.Context, params types.ParamsUpdateFirewall) (*types.ModelGetFirewall, error) {
	vdcGroupID, _, err := resolveVdcGroupRef(ctx, c.c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdc group: %w", opUpdateFirewall, err)
	}

	if err := dfwEnableOrDisable(ctx, c, vdcGroupID, params.Enabled); err != nil {
		return nil, fmt.Errorf("%s: set enabled state: %w", opUpdateFirewall, err)
	}

	epRules := endpoints.UpdateDfwRules()
	respRules, err := c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], vdcGroupID),
		cav.SetBody(firewallRulesBody(params.Rules)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: update rules: %w", opUpdateFirewall, err)
	}

	rules, ok := respRules.Result().(*itypes.ApiDistributedFirewallRules)
	if !ok || rules == nil {
		return nil, fmt.Errorf("%s: unexpected rules response type %T", opUpdateFirewall, respRules.Result())
	}

	return &types.ModelGetFirewall{
		Enabled: params.Enabled,
		Rules:   rules.ToModel(),
	}, nil
}

// DeleteFirewall clears distributed firewall rules and deactivates the feature.
func (c *Client) DeleteFirewall(ctx context.Context, params types.ParamsDeleteFirewall) error {
	vdcGroupID, _, err := resolveVdcGroupRef(ctx, c.c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return fmt.Errorf("%s: resolve vdc group: %w", opDeleteFirewall, err)
	}

	epRules := endpoints.UpdateDfwRules()
	if _, err = c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], vdcGroupID),
		cav.SetBody(itypes.ApiDistributedFirewallRules{}),
	); err != nil {
		return fmt.Errorf("%s: clear rules: %w", opDeleteFirewall, err)
	}

	if err := dfwEnableOrDisable(ctx, c, vdcGroupID, false); err != nil {
		return fmt.Errorf("%s: disable default policy: %w", opDeleteFirewall, err)
	}

	epPolicies := endpoints.UpdateDfwPolicies()
	if _, err = c.c.Do(
		ctx,
		epPolicies,
		cav.WithPathParam(epPolicies.PathParams[0], vdcGroupID),
		cav.SetBody(itypes.ApiDfwPolicies{Enabled: false}),
	); err != nil {
		return fmt.Errorf("%s: deactivate: %w", opDeleteFirewall, err)
	}

	return nil
}
