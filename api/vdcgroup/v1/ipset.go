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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListIPSet   = "IPSet.List"
	opGetIPSet    = "IPSet.Get"
	opCreateIPSet = "IPSet.Create"
	opUpdateIPSet = "IPSet.Update"
	opDeleteIPSet = "IPSet.Delete"
)

// ListIPSet lists NSX-T IP sets owned by VDC group.
func (c *Client) ListIPSet(ctx context.Context, params types.ParamsListIPSet) (*types.ModelListFirewallGroup, error) {
	model, err := inetworkobjects.ListFirewallGroupsByType(ctx, c.c, params.VDCGroupID, params.VDCGroupName, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, id, name string) (inetworkobjects.VDCGroupRef, error) {
		vdcGroup, err := c.GetVDCGroup(ctx, types.ParamsGetVDCGroup{ID: id, Name: name})
		if err != nil {
			return inetworkobjects.VDCGroupRef{}, err
		}

		return inetworkobjects.VDCGroupRef{ID: vdcGroup.ID, Name: vdcGroup.Name}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListIPSet, err)
	}

	return model, nil
}

// GetIPSet gets NSX-T IP set by ID or name within VDC group.
func (c *Client) GetIPSet(ctx context.Context, params types.ParamsGetIPSet) (*types.ModelGetFirewallGroup, error) {
	model, err := inetworkobjects.GetFirewallGroupModel(ctx, params.ID, params.Name, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, idOrName, typeValue string) (*itypes.APIResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetIPSet, err)
	}

	return model, nil
}

// CreateIPSet creates NSX-T IP set within VDC group.
func (c *Client) CreateIPSet(ctx context.Context, params types.ParamsCreateIPSet) (*types.ModelGetFirewallGroup, error) {
	body, err := createIPSetBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateIPSet, err)
	}

	ep := endpoints.CreateFirewallGroup()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateIPSet, err)
	}

	created, ok := resp.Result().(*itypes.APIResponseFirewallGroup)
	if !ok || created == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateIPSet, resp.Result())
	}

	fwGroup, err := getFirewallGroupWithRetry(ctx, c.c, created.ID, itypes.FirewallGroupTypeIPSet)
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateIPSet, err)
	}

	model := fwGroup.ToModel()
	return &model, nil
}

// UpdateIPSet updates NSX-T IP set within VDC group.
// Update is full-replace for IP addresses when IPAddresses provided.
func (c *Client) UpdateIPSet(ctx context.Context, params types.ParamsUpdateIPSet) (*types.ModelGetFirewallGroup, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, idOrName, typeValue string) (*itypes.APIResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateIPSet, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	ipAddresses := current.IPAddresses
	if len(params.IPAddresses) != 0 {
		ipAddresses = params.IPAddresses
	}

	body := itypes.APIRequestFirewallGroup{
		ID:          current.ID,
		Name:        current.Name,
		Description: description,
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		IPAddresses: ipAddresses,
		OwnerRef:    current.OwnerRef,
	}

	if err := inetworkobjects.PutFirewallGroup(ctx, c.c, body); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateIPSet, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteIPSet deletes NSX-T IP set from VDC group.
func (c *Client) DeleteIPSet(ctx context.Context, params types.ParamsDeleteIPSet) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := inetworkobjects.ResolveFirewallGroupTarget(ctx, idOrName, itypes.FirewallGroupTypeIPSet, func(ctx context.Context, idOrName, typeValue string) (*itypes.APIResponseFirewallGroup, error) {
		return getFirewallGroupWithRetry(ctx, c.c, idOrName, typeValue)
	})
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteIPSet, err)
	}

	if err := inetworkobjects.DeleteFirewallGroup(ctx, c.c, current.ID); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteIPSet, err)
	}

	return nil
}

func createIPSetBody(ctx context.Context, c cav.Client, params types.ParamsCreateIPSet) (itypes.APIRequestFirewallGroup, error) {
	vdcGroupID, vdcGroupName, err := resolveVDCGroupRef(ctx, c, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return itypes.APIRequestFirewallGroup{}, err
	}

	return itypes.APIRequestFirewallGroup{
		Name:        params.Name,
		Description: params.Description,
		TypeValue:   itypes.FirewallGroupTypeIPSet,
		IPAddresses: params.IPAddresses,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
	}, nil
}
