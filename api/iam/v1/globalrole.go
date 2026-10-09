/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
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
	opListGlobalRoles               = "IAM.ListGlobalRoles"
	opGetGlobalRole                 = "IAM.GetGlobalRole"
	opCreateGlobalRole              = "IAM.CreateGlobalRole"
	opUpdateGlobalRole              = "IAM.UpdateGlobalRole"
	opDeleteGlobalRole              = "IAM.DeleteGlobalRole"
	opListGlobalRoleRights          = "IAM.ListGlobalRoleRights"
	opAddGlobalRoleRights           = "IAM.AddGlobalRoleRights"
	opReplaceGlobalRoleRights       = "IAM.ReplaceGlobalRoleRights"
	opListGlobalRoleTenants         = "IAM.ListGlobalRoleTenants"
	opSetGlobalRoleTenants          = "IAM.SetGlobalRoleTenants"
	opPublishGlobalRoleTenants      = "IAM.PublishGlobalRoleTenants"
	opUnpublishGlobalRoleTenants    = "IAM.UnpublishGlobalRoleTenants"
	opPublishAllGlobalRoleTenants   = "IAM.PublishAllGlobalRoleTenants"
	opUnpublishAllGlobalRoleTenants = "IAM.UnpublishAllGlobalRoleTenants"
)

// ListGlobalRoles lists all global roles.
func (c *Client) ListGlobalRoles(ctx context.Context) ([]*types.ModelGlobalRole, error) {
	ep := endpoints.ListGlobalRoles()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.OverrideSetResult(new(itypes.APIResponseListGlobalRoles)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListGlobalRoles, err)
	}

	roles, ok := resp.Result().(*itypes.APIResponseListGlobalRoles)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListGlobalRoles, resp.Result())
	}

	return roles.ToModel(), nil
}

// GetGlobalRole retrieves a global role by ID.
func (c *Client) GetGlobalRole(ctx context.Context, params types.ParamsGetGlobalRole) (*types.ModelGlobalRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetGlobalRole, err)
	}

	ep := endpoints.GetGlobalRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetGlobalRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseGlobalRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetGlobalRole, resp.Result())
	}

	return role.ToModel(), nil
}

// CreateGlobalRole creates a new global role.
func (c *Client) CreateGlobalRole(ctx context.Context, params types.ParamsCreateGlobalRole) (*types.ModelGlobalRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opCreateGlobalRole, err)
	}

	body, err := modelGlobalRoleToAPIRequest(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateGlobalRole, err)
	}

	ep := endpoints.CreateGlobalRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateGlobalRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseGlobalRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opCreateGlobalRole, resp.Result())
	}

	return role.ToModel(), nil
}

// UpdateGlobalRole updates an existing global role.
func (c *Client) UpdateGlobalRole(ctx context.Context, params types.ParamsUpdateGlobalRole) (*types.ModelGlobalRole, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateGlobalRole, err)
	}

	body, err := modelGlobalRoleToAPIRequest(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateGlobalRole, err)
	}

	ep := endpoints.UpdateGlobalRole()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRole)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateGlobalRole, err)
	}

	role, ok := resp.Result().(*itypes.APIResponseGlobalRole)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateGlobalRole, resp.Result())
	}

	return role.ToModel(), nil
}

// DeleteGlobalRole deletes a global role by ID.
func (c *Client) DeleteGlobalRole(ctx context.Context, params types.ParamsDeleteGlobalRole) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opDeleteGlobalRole, err)
	}

	ep := endpoints.DeleteGlobalRole()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opDeleteGlobalRole, err)
	}

	return nil
}

// ListGlobalRoleRights lists the rights of a global role.
func (c *Client) ListGlobalRoleRights(ctx context.Context, params types.ParamsListGlobalRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opListGlobalRoleRights, err)
	}

	ep := endpoints.ListGlobalRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListGlobalRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseGlobalRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListGlobalRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// AddGlobalRoleRights adds rights to a global role.
func (c *Client) AddGlobalRoleRights(ctx context.Context, params types.ParamsAddGlobalRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opAddGlobalRoleRights, err)
	}

	body := itypes.APIRequestGlobalRoleRights{
		Rights: make([]itypes.APIRequestGlobalRoleRight, len(params.Rights)),
	}
	for i, r := range params.Rights {
		body.Rights[i] = itypes.APIRequestGlobalRoleRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	ep := endpoints.AddGlobalRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opAddGlobalRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseGlobalRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opAddGlobalRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// ReplaceGlobalRoleRights replaces the rights of a global role.
func (c *Client) ReplaceGlobalRoleRights(ctx context.Context, params types.ParamsReplaceGlobalRoleRights) ([]*types.ModelRight, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opReplaceGlobalRoleRights, err)
	}

	body := itypes.APIRequestGlobalRoleRights{
		Rights: make([]itypes.APIRequestGlobalRoleRight, len(params.Rights)),
	}
	for i, r := range params.Rights {
		body.Rights[i] = itypes.APIRequestGlobalRoleRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}

	ep := endpoints.ReplaceGlobalRoleRights()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRoleRights)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opReplaceGlobalRoleRights, err)
	}

	rights, ok := resp.Result().(*itypes.APIResponseGlobalRoleRights)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opReplaceGlobalRoleRights, resp.Result())
	}

	return rights.ToModel(), nil
}

// ListGlobalRoleTenants lists the tenants of a global role.
func (c *Client) ListGlobalRoleTenants(ctx context.Context, params types.ParamsListGlobalRoleTenants) ([]*types.ModelTenant, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opListGlobalRoleTenants, err)
	}

	ep := endpoints.ListGlobalRoleTenants()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRoleTenants)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListGlobalRoleTenants, err)
	}

	tenants, ok := resp.Result().(*itypes.APIResponseGlobalRoleTenants)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opListGlobalRoleTenants, resp.Result())
	}

	return tenants.ToModel(), nil
}

// SetGlobalRoleTenants sets the tenants of a global role.
func (c *Client) SetGlobalRoleTenants(ctx context.Context, params types.ParamsSetGlobalRoleTenants) ([]*types.ModelTenant, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opSetGlobalRoleTenants, err)
	}

	body := globalRoleTenantsToAPIRequest(params.Tenants)

	ep := endpoints.SetGlobalRoleTenants()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
		cav.OverrideSetResult(new(itypes.APIResponseGlobalRoleTenants)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opSetGlobalRoleTenants, err)
	}

	tenants, ok := resp.Result().(*itypes.APIResponseGlobalRoleTenants)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opSetGlobalRoleTenants, resp.Result())
	}

	return tenants.ToModel(), nil
}

// PublishGlobalRoleTenants publishes the tenants of a global role.
func (c *Client) PublishGlobalRoleTenants(ctx context.Context, params types.ParamsPublishGlobalRoleTenants) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opPublishGlobalRoleTenants, err)
	}

	body := globalRoleTenantsToAPIRequest(params.Tenants)

	ep := endpoints.PublishGlobalRoleTenants()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opPublishGlobalRoleTenants, err)
	}

	return nil
}

// UnpublishGlobalRoleTenants unpublishes the tenants of a global role.
func (c *Client) UnpublishGlobalRoleTenants(ctx context.Context, params types.ParamsUnpublishGlobalRoleTenants) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opUnpublishGlobalRoleTenants, err)
	}

	body := globalRoleTenantsToAPIRequest(params.Tenants)

	ep := endpoints.UnpublishGlobalRoleTenants()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
		cav.SetBody(body),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opUnpublishGlobalRoleTenants, err)
	}

	return nil
}

// PublishAllGlobalRoleTenants publishes all tenants of a global role.
func (c *Client) PublishAllGlobalRoleTenants(ctx context.Context, params types.ParamsPublishAllGlobalRoleTenants) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opPublishAllGlobalRoleTenants, err)
	}

	ep := endpoints.PublishAllGlobalRoleTenants()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opPublishAllGlobalRoleTenants, err)
	}

	return nil
}

// UnpublishAllGlobalRoleTenants unpublishes all tenants of a global role.
func (c *Client) UnpublishAllGlobalRoleTenants(ctx context.Context, params types.ParamsUnpublishAllGlobalRoleTenants) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("%s: validate: %w", opUnpublishAllGlobalRoleTenants, err)
	}

	ep := endpoints.UnpublishAllGlobalRoleTenants()
	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.ID),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", opUnpublishAllGlobalRoleTenants, err)
	}

	return nil
}

// modelGlobalRoleToAPIRequest converts public params to internal API request.
func modelGlobalRoleToAPIRequest(params any) (itypes.APIRequestGlobalRole, error) {
	switch p := params.(type) {
	case types.ParamsCreateGlobalRole:
		return itypes.APIRequestGlobalRole{
			Name:        p.Name,
			Description: p.Description,
			Rights:      entityReferencesToAPI(p.Rights),
		}, nil
	case types.ParamsUpdateGlobalRole:
		return itypes.APIRequestGlobalRole{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Rights:      entityReferencesToAPI(p.Rights),
		}, nil
	default:
		return itypes.APIRequestGlobalRole{}, fmt.Errorf("unsupported params type %T", params)
	}
}

func entityReferencesToAPI(refs []types.ParamsEntityReference) []itypes.APIRequestGlobalRoleRight {
	out := make([]itypes.APIRequestGlobalRoleRight, len(refs))
	for i, r := range refs {
		out[i] = itypes.APIRequestGlobalRoleRight{
			ID:   r.ID,
			Name: r.Name,
		}
	}
	return out
}

func globalRoleTenantsToAPIRequest(tenants []types.ParamsEntityReference) itypes.APIRequestGlobalRoleTenants {
	out := make([]itypes.APIRequestGlobalRoleTenant, len(tenants))
	for i, t := range tenants {
		out[i] = itypes.APIRequestGlobalRoleTenant{
			ID:   t.ID,
			Name: t.Name,
		}
	}
	return itypes.APIRequestGlobalRoleTenants{Tenants: out}
}
