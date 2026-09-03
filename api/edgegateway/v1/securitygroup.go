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
	stderrors "errors"
	"fmt"
	"time"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListEdgeGatewaySecurityGroup   = "EdgeGateway.SecurityGroup.List"
	opGetEdgeGatewaySecurityGroup    = "EdgeGateway.SecurityGroup.Get"
	opCreateEdgeGatewaySecurityGroup = "EdgeGateway.SecurityGroup.Create"
	opUpdateEdgeGatewaySecurityGroup = "EdgeGateway.SecurityGroup.Update"
	opDeleteEdgeGatewaySecurityGroup = "EdgeGateway.SecurityGroup.Delete"
)

func resolveEdgeGatewayOwnerRef(ctx context.Context, c cav.Client, edgeGatewayID, edgeGatewayName string) (*itypes.ApiObjectReference, error) {
	ep := endpoints.GetEdgeGateway()
	identifier := edgeGatewayID
	if identifier == "" {
		identifier = edgeGatewayName
	}

	resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], identifier))
	if err != nil {
		return nil, err
	}

	edgeGateway := resp.Result().(*itypes.ApiResponseEdgegateway)
	if edgeGateway.OwnerRef == nil {
		return nil, fmt.Errorf("edge gateway owner is missing")
	}

	return &itypes.ApiObjectReference{ID: edgeGateway.OwnerRef.ID, Name: edgeGateway.OwnerRef.Name}, nil
}

func getEdgeGatewaySecurityGroup(ctx context.Context, c cav.Client, id, name, edgeGatewayID, edgeGatewayName string) (*itypes.ApiResponseFirewallGroup, error) {
	if id != "" {
		ep := endpoints.GetFirewallGroup()
		resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.ApiResponseFirewallGroup), nil
	}

	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c, edgeGatewayID, edgeGatewayName)
	if err != nil {
		return nil, err
	}

	ep := endpoints.ListFirewallGroup()
	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;ownerRef.id==%s;name==%s", itypes.FirewallGroupTypeSecurityGroup, ownerRef.ID, name)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListFirewallGroup)
	if len(list.Values) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetFirewallGroup", StatusCode: 404, Message: fmt.Sprintf("security group %q not found", name)}
	}
	if len(list.Values) > 1 {
		return nil, pkgerrors.Newf("multiple security groups found for %q", name)
	}

	return &list.Values[0], nil
}

func getEdgeGatewaySecurityGroupWithRetry(ctx context.Context, c cav.Client, id, name, edgeGatewayID, edgeGatewayName string) (*itypes.ApiResponseFirewallGroup, error) {
	const maxAttempts = 5

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		group, err := getEdgeGatewaySecurityGroup(ctx, c, id, name, edgeGatewayID, edgeGatewayName)
		if err == nil {
			return group, nil
		}

		var apiErr *pkgerrors.APIError
		if !stderrors.As(err, &apiErr) || !apiErr.IsNotFound() {
			return nil, err
		}

		lastErr = err
	}

	return nil, lastErr
}

func securityGroupMembersToRefs(members []types.ParamsFirewallGroupMember) []itypes.ApiObjectReference {
	refs := make([]itypes.ApiObjectReference, 0, len(members))
	for _, member := range members {
		refs = append(refs, itypes.ApiObjectReference{ID: member.ID, Name: member.Name})
	}
	return refs
}

// ListSecurityGroup lists security groups owned by an edge gateway.
func (c *Client) ListSecurityGroup(ctx context.Context, params types.ParamsListEdgeGatewaySecurityGroup) (*types.ModelListFirewallGroup, error) {
	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opListEdgeGatewaySecurityGroup, err)
	}

	ep := endpoints.ListFirewallGroup()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;ownerRef.id==%s", itypes.FirewallGroupTypeSecurityGroup, ownerRef.ID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListEdgeGatewaySecurityGroup, err)
	}

	return resp.Result().(*itypes.ApiResponseListFirewallGroup).ToModel(), nil
}

// GetSecurityGroup returns a security group by ID or name for an edge gateway.
func (c *Client) GetSecurityGroup(ctx context.Context, params types.ParamsGetEdgeGatewaySecurityGroup) (*types.ModelGetFirewallGroup, error) {
	group, err := getEdgeGatewaySecurityGroup(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetEdgeGatewaySecurityGroup, err)
	}

	model := group.ToModel()
	return &model, nil
}

// CreateSecurityGroup creates a security group for a VDC group-backed edge gateway.
func (c *Client) CreateSecurityGroup(ctx context.Context, params types.ParamsCreateEdgeGatewaySecurityGroup) (*types.ModelGetFirewallGroup, error) {
	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opCreateEdgeGatewaySecurityGroup, err)
	}
	if err := validators.New().Var(ownerRef.ID, "urn=vdcGroup"); err != nil {
		return nil, fmt.Errorf("%s: edge gateway owner must be a vdc group", opCreateEdgeGatewaySecurityGroup)
	}

	body := itypes.ApiRequestFirewallGroup{
		Name:        params.Name,
		Description: params.Description,
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupMembersToRefs(params.Members),
		OwnerRef:    ownerRef,
	}

	ep := endpoints.CreateFirewallGroup()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGatewaySecurityGroup, err)
	}

	created, ok := resp.Result().(*itypes.ApiResponseFirewallGroup)
	if !ok || created == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateEdgeGatewaySecurityGroup, resp.Result())
	}

	group, err := getEdgeGatewaySecurityGroupWithRetry(ctx, c.c, created.ID, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateEdgeGatewaySecurityGroup, err)
	}

	model := group.ToModel()
	return &model, nil
}

// UpdateSecurityGroup replaces security group fields for an edge gateway.
func (c *Client) UpdateSecurityGroup(ctx context.Context, params types.ParamsUpdateEdgeGatewaySecurityGroup) (*types.ModelGetFirewallGroup, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeSecurityGroup, func(ctx context.Context, idOrName, _ string) (*itypes.ApiResponseFirewallGroup, error) {
		return getEdgeGatewaySecurityGroup(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateEdgeGatewaySecurityGroup, err)
	}

	body := itypes.ApiRequestFirewallGroup{
		ID:          current.ID,
		Name:        current.Name,
		Description: current.Description,
		Members:     current.Members,
		OwnerRef:    current.OwnerRef,
		TypeValue:   itypes.FirewallGroupTypeSecurityGroup,
	}
	if params.Description != "" {
		body.Description = params.Description
	}
	if len(params.Members) > 0 {
		body.Members = securityGroupMembersToRefs(params.Members)
	}

	if err := inetworkobjects.PutFirewallGroup(ctx, c.c, body); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateEdgeGatewaySecurityGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteSecurityGroup deletes a security group from an edge gateway.
func (c *Client) DeleteSecurityGroup(ctx context.Context, params types.ParamsDeleteEdgeGatewaySecurityGroup) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeSecurityGroup, func(ctx context.Context, idOrName, _ string) (*itypes.ApiResponseFirewallGroup, error) {
		return getEdgeGatewaySecurityGroup(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	})
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteEdgeGatewaySecurityGroup, err)
	}

	if err := inetworkobjects.DeleteFirewallGroup(ctx, c.c, current.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteEdgeGatewaySecurityGroup, err)
	}

	return nil
}
