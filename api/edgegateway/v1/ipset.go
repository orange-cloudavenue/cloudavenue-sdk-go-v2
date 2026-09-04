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
	opListEdgeGatewayIPSet   = "EdgeGateway.IPSet.List"
	opGetEdgeGatewayIPSet    = "EdgeGateway.IPSet.Get"
	opCreateEdgeGatewayIPSet = "EdgeGateway.IPSet.Create"
	opUpdateEdgeGatewayIPSet = "EdgeGateway.IPSet.Update"
	opDeleteEdgeGatewayIPSet = "EdgeGateway.IPSet.Delete"
)

// ListIPSet lists NSX-T IP sets owned by edge gateway.
func (c *Client) ListIPSet(ctx context.Context, params types.ParamsListEdgeGatewayIPSet) (*types.ModelListFirewallGroup, error) {
	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opListEdgeGatewayIPSet, err)
	}

	ep := endpoints.ListFirewallGroup()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;ownerRef.id==%s", itypes.FirewallGroupTypeIPSet, ownerRef.ID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListEdgeGatewayIPSet, err)
	}

	return resp.Result().(*itypes.APIResponseListFirewallGroup).ToModel(), nil
}

// GetIPSet gets NSX-T IP set by ID or name within edge gateway.
func (c *Client) GetIPSet(ctx context.Context, params types.ParamsGetEdgeGatewayIPSet) (*types.ModelGetFirewallGroup, error) {
	group, err := getEdgeGatewayIPSet(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetEdgeGatewayIPSet, err)
	}

	model := group.ToModel()
	return &model, nil
}

// CreateIPSet creates NSX-T IP set within edge gateway.
func (c *Client) CreateIPSet(ctx context.Context, params types.ParamsCreateEdgeGatewayIPSet) (*types.ModelGetFirewallGroup, error) {
	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opCreateEdgeGatewayIPSet, err)
	}
	if err := validators.New().Var(ownerRef.ID, "urn=vdcGroup"); err != nil {
		return nil, fmt.Errorf("%s: edge gateway owner must be a vdc group", opCreateEdgeGatewayIPSet)
	}

	body := itypes.APIRequestFirewallGroup{
		Name:        params.Name,
		Description: params.Description,
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		IPAddresses: params.IPAddresses,
		OwnerRef:    ownerRef,
	}

	ep := endpoints.CreateFirewallGroup()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGatewayIPSet, err)
	}

	created, ok := resp.Result().(*itypes.APIResponseFirewallGroup)
	if !ok || created == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateEdgeGatewayIPSet, resp.Result())
	}

	group, err := getEdgeGatewayIPSetWithRetry(ctx, c.c, created.ID, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateEdgeGatewayIPSet, err)
	}

	model := group.ToModel()
	return &model, nil
}

// UpdateIPSet updates NSX-T IP set within edge gateway.
func (c *Client) UpdateIPSet(ctx context.Context, params types.ParamsUpdateEdgeGatewayIPSet) (*types.ModelGetFirewallGroup, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, idOrName, _ string) (*itypes.APIResponseFirewallGroup, error) {
		return getEdgeGatewayIPSet(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateEdgeGatewayIPSet, err)
	}

	body := itypes.APIRequestFirewallGroup{
		ID:          current.ID,
		Name:        current.Name,
		Description: current.Description,
		IPAddresses: current.IPAddresses,
		OwnerRef:    current.OwnerRef,
		TypeValue:   itypes.FirewallGroupTypeIPSet,
	}
	if params.Description != "" {
		body.Description = params.Description
	}
	if len(params.IPAddresses) > 0 {
		body.IPAddresses = params.IPAddresses
	}

	if err := inetworkobjects.PutFirewallGroup(ctx, c.c, body); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateEdgeGatewayIPSet, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteIPSet deletes NSX-T IP set from edge gateway.
func (c *Client) DeleteIPSet(ctx context.Context, params types.ParamsDeleteEdgeGatewayIPSet) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, idOrName, _ string) (*itypes.APIResponseFirewallGroup, error) {
		return getEdgeGatewayIPSet(ctx, c.c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	})
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteEdgeGatewayIPSet, err)
	}

	if err := inetworkobjects.DeleteFirewallGroup(ctx, c.c, current.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteEdgeGatewayIPSet, err)
	}

	return nil
}

func getEdgeGatewayIPSet(ctx context.Context, c cav.Client, id, name, edgeGatewayID, edgeGatewayName string) (*itypes.APIResponseFirewallGroup, error) {
	if id != "" {
		ep := endpoints.GetFirewallGroup()
		resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.APIResponseFirewallGroup), nil
	}

	ownerRef, err := resolveEdgeGatewayOwnerRef(ctx, c, edgeGatewayID, edgeGatewayName)
	if err != nil {
		return nil, err
	}

	ep := endpoints.ListFirewallGroup()
	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("typeValue==%s;ownerRef.id==%s;name==%s", itypes.FirewallGroupTypeIPSet, ownerRef.ID, name)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.APIResponseListFirewallGroup)
	if len(list.Values) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetFirewallGroup", StatusCode: 404, Message: fmt.Sprintf("ip set %q not found", name)}
	}
	if len(list.Values) > 1 {
		return nil, pkgerrors.Newf("multiple ip sets found for %q", name)
	}

	return &list.Values[0], nil
}

func getEdgeGatewayIPSetWithRetry(ctx context.Context, c cav.Client, id, name, edgeGatewayID, edgeGatewayName string) (*itypes.APIResponseFirewallGroup, error) {
	const maxAttempts = 5

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		group, err := getEdgeGatewayIPSet(ctx, c, id, name, edgeGatewayID, edgeGatewayName)
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
