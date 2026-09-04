/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetEdgeGatewayFirewall    = "EdgeGateway.Firewall.Get"
	opCreateEdgeGatewayFirewall = "EdgeGateway.Firewall.Create"
	opUpdateEdgeGatewayFirewall = "EdgeGateway.Firewall.Update"
	opDeleteEdgeGatewayFirewall = "EdgeGateway.Firewall.Delete"
)

func resolveEdgeGatewayVDCGroupOwnerRef(ctx context.Context, c *Client, edgeGatewayID, edgeGatewayName string) (*itypes.APIObjectReference, error) {
	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c.c, edgeGatewayID, edgeGatewayName)
	if err != nil {
		return nil, err
	}
	if err := validators.New().Var(ownerRef.ID, "urn=vdcGroup"); err != nil {
		return nil, fmt.Errorf("edge gateway owner must be a vdc group")
	}

	return ownerRef, nil
}

func edgeGatewayToAPIObjectReferences(members []types.ParamsFirewallGroupMember) []itypes.APIObjectReference {
	refs := make([]itypes.APIObjectReference, 0, len(members))
	for _, member := range members {
		refs = append(refs, itypes.APIObjectReference{ID: member.ID, Name: member.Name})
	}
	return refs
}

func edgeGatewayToAPIDistributedFirewallRules(rules []types.ParamsFirewallRule) []itypes.APIDistributedFirewallRule {
	out := make([]itypes.APIDistributedFirewallRule, 0, len(rules))
	for _, rule := range rules {
		out = append(out, itypes.APIDistributedFirewallRule{
			ID:                        rule.ID,
			Name:                      rule.Name,
			Description:               rule.Description,
			Comments:                  rule.Description,
			ApplicationPortProfiles:   edgeGatewayToAPIObjectReferences(rule.ApplicationPortProfiles),
			SourceFirewallGroups:      edgeGatewayToAPIObjectReferences(rule.SourceFirewallGroups),
			DestinationFirewallGroups: edgeGatewayToAPIObjectReferences(rule.DestinationFirewallGroups),
			NetworkContextProfiles:    edgeGatewayToAPIObjectReferences(rule.NetworkContextProfiles),
			Direction:                 rule.Direction,
			Enabled:                   rule.Enabled,
			IPProtocol:                rule.IPProtocol,
			Logging:                   rule.Logging,
			ActionValue:               rule.Action,
			SourceGroupsExcluded:      rule.SourceGroupsExcluded,
			DestinationGroupsExcluded: rule.DestinationGroupsExcluded,
		})
	}
	return out
}

func edgeGatewayDFWEnableOrDisable(ctx context.Context, c *Client, vdcGroupID string, enable bool) error {
	epGet := endpoints.GetDFWPolicies()
	resp, err := c.c.Do(ctx, epGet, cav.WithPathParam(epGet.PathParams[0], vdcGroupID))
	if err != nil {
		return err
	}

	policies := resp.Result().(*itypes.APIDFWPolicies)
	isEnabled := policies.DefaultPolicy != nil && policies.DefaultPolicy.Enabled != nil && *policies.DefaultPolicy.Enabled
	if isEnabled == enable {
		return nil
	}

	defaultPolicy := itypes.APIDfwDefaultPolicy{Enabled: &enable}
	if policies.DefaultPolicy != nil {
		defaultPolicy.Description = policies.DefaultPolicy.Description
		defaultPolicy.ID = policies.DefaultPolicy.ID
		defaultPolicy.Name = policies.DefaultPolicy.Name
		defaultPolicy.Version = policies.DefaultPolicy.Version
	}
	if defaultPolicy.Name == "" {
		defaultPolicy.Name = "Default"
	}

	epUpdate := endpoints.UpdateDFWDefaultPolicy()
	_, err = c.c.Do(
		ctx,
		epUpdate,
		cav.WithPathParam(epUpdate.PathParams[0], vdcGroupID),
		cav.SetBody(defaultPolicy),
	)
	return err
}

func edgeGatewayFirewallRulesBody(rules []types.ParamsFirewallRule) itypes.APIDistributedFirewallRules {
	return itypes.APIDistributedFirewallRules{Values: edgeGatewayToAPIDistributedFirewallRules(rules)}
}

// GetFirewall returns distributed firewall state for an edge gateway owner VDC group.
func (c *Client) GetFirewall(ctx context.Context, params types.ParamsGetEdgeGatewayFirewall) (*types.ModelGetFirewall, error) {
	ownerRef, err := resolveEdgeGatewayVDCGroupOwnerRef(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opGetEdgeGatewayFirewall, err)
	}

	epPolicies := endpoints.GetDFWPolicies()
	respPolicies, err := c.c.Do(ctx, epPolicies, cav.WithPathParam(epPolicies.PathParams[0], ownerRef.ID))
	if err != nil {
		return nil, fmt.Errorf("%s: get policies: %w", opGetEdgeGatewayFirewall, err)
	}

	policies := respPolicies.Result().(*itypes.APIDFWPolicies)
	enabled := policies.DefaultPolicy != nil && policies.DefaultPolicy.Enabled != nil && *policies.DefaultPolicy.Enabled

	epRules := endpoints.GetDFWRules()
	respRules, err := c.c.Do(ctx, epRules, cav.WithPathParam(epRules.PathParams[0], ownerRef.ID))
	if err != nil {
		return nil, fmt.Errorf("%s: get rules: %w", opGetEdgeGatewayFirewall, err)
	}

	rules := respRules.Result().(*itypes.APIDistributedFirewallRules)
	return &types.ModelGetFirewall{Enabled: enabled, Rules: rules.ToModel()}, nil
}

// CreateFirewall activates distributed firewalling and sets initial rules.
func (c *Client) CreateFirewall(ctx context.Context, params types.ParamsCreateEdgeGatewayFirewall) (*types.ModelGetFirewall, error) {
	ownerRef, err := resolveEdgeGatewayVDCGroupOwnerRef(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opCreateEdgeGatewayFirewall, err)
	}

	epPolicies := endpoints.UpdateDFWPolicies()
	if _, err := c.c.Do(
		ctx,
		epPolicies,
		cav.WithPathParam(epPolicies.PathParams[0], ownerRef.ID),
		cav.SetBody(itypes.APIDFWPolicies{Enabled: true}),
	); err != nil {
		return nil, fmt.Errorf("%s: activate: %w", opCreateEdgeGatewayFirewall, err)
	}

	if err := edgeGatewayDFWEnableOrDisable(ctx, c, ownerRef.ID, params.Enabled); err != nil {
		return nil, fmt.Errorf("%s: set enabled state: %w", opCreateEdgeGatewayFirewall, err)
	}

	epRules := endpoints.UpdateDFWRules()
	respRules, err := c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], ownerRef.ID),
		cav.SetBody(edgeGatewayFirewallRulesBody(params.Rules)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: update rules: %w", opCreateEdgeGatewayFirewall, err)
	}

	rules, ok := respRules.Result().(*itypes.APIDistributedFirewallRules)
	if !ok || rules == nil {
		return nil, fmt.Errorf("%s: unexpected rules response type %T", opCreateEdgeGatewayFirewall, respRules.Result())
	}

	return &types.ModelGetFirewall{Enabled: params.Enabled, Rules: rules.ToModel()}, nil
}

// UpdateFirewall replaces distributed firewall rules and enabled state.
func (c *Client) UpdateFirewall(ctx context.Context, params types.ParamsUpdateEdgeGatewayFirewall) (*types.ModelGetFirewall, error) {
	ownerRef, err := resolveEdgeGatewayVDCGroupOwnerRef(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opUpdateEdgeGatewayFirewall, err)
	}

	if err := edgeGatewayDFWEnableOrDisable(ctx, c, ownerRef.ID, params.Enabled); err != nil {
		return nil, fmt.Errorf("%s: set enabled state: %w", opUpdateEdgeGatewayFirewall, err)
	}

	epRules := endpoints.UpdateDFWRules()
	respRules, err := c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], ownerRef.ID),
		cav.SetBody(edgeGatewayFirewallRulesBody(params.Rules)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: update rules: %w", opUpdateEdgeGatewayFirewall, err)
	}

	rules, ok := respRules.Result().(*itypes.APIDistributedFirewallRules)
	if !ok || rules == nil {
		return nil, fmt.Errorf("%s: unexpected rules response type %T", opUpdateEdgeGatewayFirewall, respRules.Result())
	}

	return &types.ModelGetFirewall{Enabled: params.Enabled, Rules: rules.ToModel()}, nil
}

// DeleteFirewall clears distributed firewall rules and deactivates the feature.
func (c *Client) DeleteFirewall(ctx context.Context, params types.ParamsDeleteEdgeGatewayFirewall) error {
	ownerRef, err := resolveEdgeGatewayVDCGroupOwnerRef(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return fmt.Errorf("%s: resolve edge gateway: %w", opDeleteEdgeGatewayFirewall, err)
	}

	epRules := endpoints.UpdateDFWRules()
	if _, err = c.c.Do(
		ctx,
		epRules,
		cav.WithPathParam(epRules.PathParams[0], ownerRef.ID),
		cav.SetBody(itypes.APIDistributedFirewallRules{}),
	); err != nil {
		return fmt.Errorf("%s: clear rules: %w", opDeleteEdgeGatewayFirewall, err)
	}

	if err := edgeGatewayDFWEnableOrDisable(ctx, c, ownerRef.ID, false); err != nil {
		return fmt.Errorf("%s: disable default policy: %w", opDeleteEdgeGatewayFirewall, err)
	}

	epPolicies := endpoints.UpdateDFWPolicies()
	if _, err = c.c.Do(
		ctx,
		epPolicies,
		cav.WithPathParam(epPolicies.PathParams[0], ownerRef.ID),
		cav.SetBody(itypes.APIDFWPolicies{Enabled: false}),
	); err != nil {
		return fmt.Errorf("%s: deactivate: %w", opDeleteEdgeGatewayFirewall, err)
	}

	return nil
}
