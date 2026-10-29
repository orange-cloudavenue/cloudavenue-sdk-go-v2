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
	"strconv"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// toAPIEntityReferences converts public reference params to their API representation.
func toAPIEntityReferences(refs []types.ParamsEntityReference) []itypes.APIEntityReference {
	result := make([]itypes.APIEntityReference, len(refs))
	for i, r := range refs {
		result[i] = itypes.APIEntityReference{ID: r.ID, Name: r.Name}
	}
	return result
}

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

// endpointQueryParam returns the endpoint query parameter registered under name.
// It falls back to a bare parameter so callers stay correct if registration evolves.
func endpointQueryParam(ep *cav.Endpoint, name string) cav.QueryParam {
	for _, p := range ep.QueryParams {
		if p.Name == name {
			return p
		}
	}
	return cav.QueryParam{Name: name}
}

// ListOrgRoles lists all org-scoped roles, fetching every page.
func (c *Client) ListOrgRoles(ctx context.Context) ([]*types.ModelRole, error) {
	var result []*types.ModelRole

	for page := 1; ; page++ {
		ep := endpoints.ListRoles()
		resp, err := c.c.Do(
			ctx,
			ep,
			cav.WithQueryParam(endpointQueryParam(ep, "page"), strconv.Itoa(page)),
			cav.OverrideSetResult(new(itypes.APIResponseListRoles)),
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", opListOrgRoles, err)
		}

		roles, ok := resp.Result().(*itypes.APIResponseListRoles)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected response type %T", opListOrgRoles, resp.Result())
		}

		result = append(result, roles.ToModel()...)
		if len(roles.Roles) == 0 || page >= roles.PageCount {
			break
		}
	}

	return result, nil
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

// ListOrgRoleRights lists the rights of a role, fetching every page.
// CloudAPI serves role rights as entity references, so only ID and Name are set.
func (c *Client) ListOrgRoleRights(ctx context.Context, params types.ParamsListOrgRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opListOrgRoleRights, err)
	}

	var result []*types.ModelRight

	for page := 1; ; page++ {
		ep := endpoints.ListRoleRights()
		resp, err := c.c.Do(
			ctx,
			ep,
			cav.WithPathParam(ep.PathParams[0], params.ID),
			cav.WithQueryParam(endpointQueryParam(ep, "page"), strconv.Itoa(page)),
			cav.OverrideSetResult(new(itypes.APIResponseRoleRights)),
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", opListOrgRoleRights, err)
		}

		rights, ok := resp.Result().(*itypes.APIResponseRoleRights)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected response type %T", opListOrgRoleRights, resp.Result())
		}

		result = append(result, rights.ToModel()...)
		if len(rights.Rights) == 0 || page >= rights.PageCount {
			break
		}
	}

	return result, nil
}

// AddOrgRoleRights adds rights to a role.
func (c *Client) AddOrgRoleRights(ctx context.Context, params types.ParamsAddOrgRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opAddOrgRoleRights, err)
	}

	body := itypes.APIRequestRoleRights{
		Rights: toAPIEntityReferences(params.Rights),
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
		Rights: toAPIEntityReferences(params.Rights),
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
		Rights: []itypes.APIEntityReference{},
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
