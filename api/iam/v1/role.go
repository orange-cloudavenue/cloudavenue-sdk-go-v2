/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListOrgRoles           = "IAM.ListOrgRoles"
	opGetOrgRole             = "IAM.GetOrgRole"
	opCreateOrgRole          = "IAM.CreateOrgRole"
	opUpdateOrgRole          = "IAM.UpdateOrgRole"
	opDeleteOrgRole          = "IAM.DeleteOrgRole"
	opListOrgRoleRights      = "IAM.ListOrgRoleRights"
	opAddOrgRoleRights       = "IAM.AddOrgRoleRights"
	opReplaceOrgRoleRights   = "IAM.ReplaceOrgRoleRights"
	opRemoveAllOrgRoleRights = "IAM.RemoveAllOrgRoleRights"
)

// ListOrgRoles lists all org-scoped roles.
func (c *Client) ListOrgRoles(ctx context.Context) ([]*types.ModelRole, error) {
	ep := endpoints.ListRoles()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.OverrideSetResult(new(itypes.APIResponseListRoles)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListOrgRoles, err)
	}

	roles, ok := resp.Result().(*itypes.APIResponseListRoles)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListOrgRoles, resp.Result())
	}

	return roles.ToModel(), nil
}

// GetOrgRole retrieves a role by ID.
func (c *Client) GetOrgRole(ctx context.Context, params types.ParamsGetOrgRole) (*types.ModelRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetOrgRole, err)
	}

	ep := endpoints.GetRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetOrgRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetOrgRole, resp.Result())
	}

	return role.ToModel(), nil
}

// CreateOrgRole creates a new role.
func (c *Client) CreateOrgRole(ctx context.Context, params types.ParamsCreateOrgRole) (*types.ModelRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateOrgRole, err)
	}

	body := &itypes.APIRequestRole{
		Name:        params.Name,
		Description: params.Description,
	}

	ep := endpoints.CreateRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateOrgRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateOrgRole, resp.Result())
	}

	return role.ToModel(), nil
}

// UpdateOrgRole updates an existing role.
func (c *Client) UpdateOrgRole(ctx context.Context, params types.ParamsUpdateOrgRole) (*types.ModelRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateOrgRole, err)
	}

	body := &itypes.APIRequestRole{
		ID:          params.ID,
		Name:        params.Name,
		Description: params.Description,
	}

	ep := endpoints.UpdateRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateOrgRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateOrgRole, resp.Result())
	}

	return role.ToModel(), nil
}

// DeleteOrgRole deletes a role by ID.
func (c *Client) DeleteOrgRole(ctx context.Context, params types.ParamsDeleteOrgRole) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteOrgRole, err)
	}

	ep := endpoints.DeleteRole()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteOrgRole, err)
	}

	return nil
}

// ListOrgRoleRights lists the rights of a role.
func (c *Client) ListOrgRoleRights(ctx context.Context, params types.ParamsListOrgRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opListOrgRoleRights, err)
	}

	ep := endpoints.ListRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListOrgRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListOrgRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// AddOrgRoleRights adds rights to a role.
func (c *Client) AddOrgRoleRights(ctx context.Context, params types.ParamsAddOrgRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opAddOrgRoleRights, err)
	}

	body := itypes.APIRequestRoleRights{
		Rights: make([]itypes.APIRequestRoleRight, len(params.Rights)),
	}
	for i, r := range params.Rights {
		body.Rights[i] = itypes.APIRequestRoleRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	ep := endpoints.AddRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opAddOrgRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opAddOrgRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// ReplaceOrgRoleRights replaces the rights of a role.
func (c *Client) ReplaceOrgRoleRights(ctx context.Context, params types.ParamsReplaceOrgRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opReplaceOrgRoleRights, err)
	}

	body := itypes.APIRequestRoleRights{
		Rights: make([]itypes.APIRequestRoleRight, len(params.Rights)),
	}
	for i, r := range params.Rights {
		body.Rights[i] = itypes.APIRequestRoleRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	ep := endpoints.ReplaceRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opReplaceOrgRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opReplaceOrgRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// RemoveAllOrgRoleRights removes all rights from a role.
// Uses PUT with empty rights list since CloudAPI has no DELETE verb for role rights.
func (c *Client) RemoveAllOrgRoleRights(ctx context.Context, params types.ParamsDeleteOrgRole) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opRemoveAllOrgRoleRights, err)
	}

	body := itypes.APIRequestRoleRights{
		Rights: []itypes.APIRequestRoleRight{},
	}

	ep := endpoints.ReplaceRoleRights()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opRemoveAllOrgRoleRights, err)
	}

	return nil
}
