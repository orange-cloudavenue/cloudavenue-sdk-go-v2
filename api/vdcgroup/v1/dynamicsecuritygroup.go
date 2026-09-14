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
	"slices"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListDynamicSecurityGroup   = "DynamicSecurityGroup.List"
	opGetDynamicSecurityGroup    = "DynamicSecurityGroup.Get"
	opCreateDynamicSecurityGroup = "DynamicSecurityGroup.Create"
	opUpdateDynamicSecurityGroup = "DynamicSecurityGroup.Update"
	opDeleteDynamicSecurityGroup = "DynamicSecurityGroup.Delete"
)

// ListDynamicSecurityGroup lists NSX-T dynamic security groups owned by VDC group.
func (c *Client) ListDynamicSecurityGroup(ctx context.Context, params types.ParamsListDynamicSecurityGroup) (*types.ModelListFirewallGroup, error) {
	model, err := inetworkobjects.ListFirewallGroupsByType(ctx, c.c, params.VdcGroupID, params.VdcGroupName, itypes.FirewallGroupTypeVMCriteria, func(ctx context.Context, id, name string) (inetworkobjects.VdcGroupRef, error) {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{ID: id, Name: name})
		if err != nil {
			return inetworkobjects.VdcGroupRef{}, err
		}

		return inetworkobjects.VdcGroupRef{ID: vdcGroup.ID, Name: vdcGroup.Name}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListDynamicSecurityGroup, err)
	}

	return model, nil
}

// GetDynamicSecurityGroup gets NSX-T dynamic security group by ID or name within VDC group.
func (c *Client) GetDynamicSecurityGroup(ctx context.Context, params types.ParamsGetDynamicSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	model, err := inetworkobjects.GetFirewallGroupModel(ctx, params.ID, params.Name, itypes.FirewallGroupTypeVMCriteria, func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetDynamicSecurityGroup, err)
	}

	return model, nil
}

// CreateDynamicSecurityGroup creates NSX-T dynamic security group within VDC group.
func (c *Client) CreateDynamicSecurityGroup(ctx context.Context, params types.ParamsCreateDynamicSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	body, err := createDynamicSecurityGroupBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateDynamicSecurityGroup, err)
	}

	ep := endpoints.CreateFirewallGroup()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateDynamicSecurityGroup, err)
	}

	created, ok := resp.Result().(*itypes.ApiResponseFirewallGroup)
	if !ok || created == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateDynamicSecurityGroup, resp.Result())
	}

	fwGroup, err := getFirewallGroupWithRetry(ctx, c.c, created.ID, itypes.FirewallGroupTypeVMCriteria)
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateDynamicSecurityGroup, err)
	}

	model := fwGroup.ToModel()
	return &model, nil
}

// UpdateDynamicSecurityGroup updates NSX-T dynamic security group within VDC group.
// Update is full-replace for criteria when Criteria provided.
func (c *Client) UpdateDynamicSecurityGroup(ctx context.Context, params types.ParamsUpdateDynamicSecurityGroup) (*types.ModelGetFirewallGroup, error) {
	if err := validateDynamicSecurityGroupCriteria(params.Criteria); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateDynamicSecurityGroup, err)
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeVMCriteria, func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateDynamicSecurityGroup, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	vmCriteria := current.VMCriteria
	if len(params.Criteria) != 0 {
		vmCriteria = toApiFirewallGroupVMCriteria(params.Criteria)
	}

	body := itypes.ApiRequestFirewallGroup{
		ID:          current.ID,
		Name:        current.Name,
		Description: description,
		TypeValue:   itypes.FirewallGroupTypeVMCriteria,
		VMCriteria:  vmCriteria,
		OwnerRef:    current.OwnerRef,
	}

	if err := inetworkobjects.PutFirewallGroup(ctx, c.c, body); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateDynamicSecurityGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteDynamicSecurityGroup deletes NSX-T dynamic security group from VDC group.
func (c *Client) DeleteDynamicSecurityGroup(ctx context.Context, params types.ParamsDeleteDynamicSecurityGroup) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeVMCriteria, func(ctx context.Context, idOrName, typeValue string) (*itypes.ApiResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteDynamicSecurityGroup, err)
	}

	if err := inetworkobjects.DeleteFirewallGroup(ctx, c.c, current.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteDynamicSecurityGroup, err)
	}

	return nil
}

func createDynamicSecurityGroupBody(ctx context.Context, c cav.Client, params types.ParamsCreateDynamicSecurityGroup) (itypes.ApiRequestFirewallGroup, error) {
	if err := validateDynamicSecurityGroupCriteria(params.Criteria); err != nil {
		return itypes.ApiRequestFirewallGroup{}, err
	}

	vdcGroupID, vdcGroupName, err := resolveVdcGroupRef(ctx, c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return itypes.ApiRequestFirewallGroup{}, err
	}

	return itypes.ApiRequestFirewallGroup{
		Name:        params.Name,
		Description: params.Description,
		TypeValue:   itypes.FirewallGroupTypeVMCriteria,
		VMCriteria:  toApiFirewallGroupVMCriteria(params.Criteria),
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil
}

// validateDynamicSecurityGroupCriteria validates dynamic security group criteria against
// supported limits and rule type/operator combinations.
func validateDynamicSecurityGroupCriteria(criteria []types.ParamsDynamicSecurityGroupCriteria) error {
	if len(criteria) > types.DynamicSecurityGroupMaxCriteria {
		return errors.Newf("dynamic security group criteria: allowed max length of criteria is %d", types.DynamicSecurityGroupMaxCriteria)
	}

	for _, c := range criteria {
		if len(c.Rules) > types.DynamicSecurityGroupMaxRulesPerCriteria {
			return errors.Newf("dynamic security group criteria: allowed max length of rules is %d", types.DynamicSecurityGroupMaxRulesPerCriteria)
		}

		for _, r := range c.Rules {
			if r.Value == "" {
				return errors.Newf("dynamic security group criteria rule: value must not be empty")
			}

			allowedOperators, ok := types.DynamicSecurityGroupAllowedOperators[r.RuleType]
			if !ok {
				return errors.Newf(
					"dynamic security group criteria rule: rule type must be one of %s or %s",
					types.DynamicSecurityGroupCriteriaRuleTypeVMName,
					types.DynamicSecurityGroupCriteriaRuleTypeVMTag,
				)
			}

			if !slices.Contains(allowedOperators, r.Operator) {
				return errors.Newf(
					"dynamic security group criteria rule: rule type and operator mismatch. If rule type is %s, then operator must be one of %v",
					r.RuleType, allowedOperators,
				)
			}
		}
	}

	return nil
}

// toApiFirewallGroupVMCriteria converts criteria params into firewall group VM criteria payload.
func toApiFirewallGroupVMCriteria(criteria []types.ParamsDynamicSecurityGroupCriteria) []itypes.ApiFirewallGroupVMCriteria {
	vmCriteria := make([]itypes.ApiFirewallGroupVMCriteria, 0, len(criteria))

	for _, c := range criteria {
		rules := make([]itypes.ApiFirewallGroupVMCriteriaRule, 0, len(c.Rules))
		for _, r := range c.Rules {
			rules = append(rules, itypes.ApiFirewallGroupVMCriteriaRule{
				AttributeType:  r.RuleType,
				Operator:       r.Operator,
				AttributeValue: r.Value,
			})
		}

		vmCriteria = append(vmCriteria, itypes.ApiFirewallGroupVMCriteria{VMCriteriaRule: rules})
	}

	return vmCriteria
}
