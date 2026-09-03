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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVdcNetwork           = "VdcNetwork.List"
	opGetVdcNetworkIsolated    = "VdcNetworkIsolated.Get"
	opCreateVdcNetworkIsolated = "VdcNetworkIsolated.Create"
	opUpdateVdcNetworkIsolated = "VdcNetworkIsolated.Update"
	opDeleteVdcNetworkIsolated = "VdcNetworkIsolated.Delete"
)

// ListVdcNetwork lists routed and isolated VDC networks owned by a VDC group.
func (c *Client) ListVdcNetwork(ctx context.Context, params types.ParamsListVdcNetwork) (*types.ModelListVdcNetwork, error) {
	resp, err := listVdcNetworkModel(ctx, c.c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListVdcNetwork, err)
	}

	return resp, nil
}

// GetVdcNetworkIsolated returns an isolated VDC network by ID or name.
func (c *Client) GetVdcNetworkIsolated(ctx context.Context, params types.ParamsGetVdcNetworkIsolated) (*types.ModelGetVdcNetwork, error) {
	resp, err := getVdcNetworkModel(ctx, c.c, params.ID, params.Name, params.VdcGroupID, params.VdcGroupName, types.VdcNetworkTypeIsolated)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetVdcNetworkIsolated, err)
	}

	return resp, nil
}

func createVdcNetworkIsolatedBody(ctx context.Context, c cav.Client, params types.ParamsCreateVdcNetworkIsolated) (itypes.ApiRequestVdcNetwork, error) {
	vdcGroupID, vdcGroupName, err := resolveVdcGroupRef(ctx, c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return itypes.ApiRequestVdcNetwork{}, err
	}

	return itypes.ApiRequestVdcNetwork{
		Name:        params.Name,
		Description: params.Description,
		NetworkType: types.VdcNetworkTypeIsolated,
		OwnerRef:    &itypes.ApiObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		Subnets: itypes.ApiVdcNetworkSubnets{
			Values: []itypes.ApiVdcNetworkSubnetValue{{
				Gateway:      params.Subnet.Gateway,
				PrefixLength: params.Subnet.PrefixLength,
				DNSServer1:   params.Subnet.DNSServer1,
				DNSServer2:   params.Subnet.DNSServer2,
				DNSSuffix:    params.Subnet.DNSSuffix,
				IPRanges: itypes.ApiVdcNetworkIPRanges{
					Values: createVdcNetworkIPRanges(params.Subnet.IPRanges),
				},
			}},
		},
		GuestVlanTaggingAllowed: params.GuestVlanTaggingAllowed,
		Shared:                  new(true),
	}, nil
}

// CreateVdcNetworkIsolated creates an isolated VDC network in a VDC group.
func (c *Client) CreateVdcNetworkIsolated(ctx context.Context, params types.ParamsCreateVdcNetworkIsolated) (*types.ModelGetVdcNetwork, error) {
	body, err := createVdcNetworkIsolatedBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateVdcNetworkIsolated, err)
	}

	ep := endpoints.CreateVdcNetwork()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVdcNetworkIsolated, err)
	}

	network, ok := resp.Result().(*itypes.ApiResponseVdcNetwork)
	if !ok || network == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateVdcNetworkIsolated, resp.Result())
	}

	model := network.ToModel()
	return &model, nil
}

// mergeVdcNetworkIsolatedUpdate applies optional update params onto the current network state.
func mergeVdcNetworkIsolatedUpdate(current *itypes.ApiResponseVdcNetwork, params types.ParamsUpdateVdcNetworkIsolated) itypes.ApiRequestVdcNetwork {
	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	subnets := current.Subnets
	if params.Subnet != nil {
		subnets = itypes.ApiVdcNetworkSubnets{
			Values: []itypes.ApiVdcNetworkSubnetValue{{
				Gateway:      params.Subnet.Gateway,
				PrefixLength: params.Subnet.PrefixLength,
				DNSServer1:   params.Subnet.DNSServer1,
				DNSServer2:   params.Subnet.DNSServer2,
				DNSSuffix:    params.Subnet.DNSSuffix,
				IPRanges:     itypes.ApiVdcNetworkIPRanges{Values: createVdcNetworkIPRanges(params.Subnet.IPRanges)},
			}},
		}
	}

	guestVlanTaggingAllowed := current.GuestVlanTaggingAllowed
	if params.GuestVlanTaggingAllowed != nil {
		guestVlanTaggingAllowed = params.GuestVlanTaggingAllowed
	}

	return itypes.ApiRequestVdcNetwork{
		ID:                      current.ID,
		Name:                    current.Name,
		Description:             description,
		NetworkType:             current.NetworkType,
		OwnerRef:                current.OwnerRef,
		Subnets:                 subnets,
		GuestVlanTaggingAllowed: guestVlanTaggingAllowed,
		Shared:                  current.Shared,
	}
}

// UpdateVdcNetworkIsolated updates an isolated VDC network in a VDC group.
func (c *Client) UpdateVdcNetworkIsolated(ctx context.Context, params types.ParamsUpdateVdcNetworkIsolated) (*types.ModelGetVdcNetwork, error) {
	current, err := getVdcNetworkWithRetry(ctx, c.c, params.ID, "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVdcNetworkIsolated, err)
	}

	if current.NetworkType != types.VdcNetworkTypeIsolated {
		return nil, fmt.Errorf("%s: %w", opUpdateVdcNetworkIsolated, errors.Newf("org vdc network %q is not an isolated network", params.ID))
	}

	body := mergeVdcNetworkIsolatedUpdate(current, params)

	ep := endpoints.UpdateVdcNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVdcNetworkIsolated, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVdcNetworkIsolated deletes an isolated VDC network from a VDC group.
func (c *Client) DeleteVdcNetworkIsolated(ctx context.Context, params types.ParamsDeleteVdcNetworkIsolated) error {
	current, err := deleteVdcNetworkTarget(ctx, c.c, params.ID, params.Name, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opDeleteVdcNetworkIsolated, err)
	}

	if current.NetworkType != types.VdcNetworkTypeIsolated {
		idOrName := params.ID
		if idOrName == "" {
			idOrName = params.Name
		}
		return fmt.Errorf("%s: %w", opDeleteVdcNetworkIsolated, errors.Newf("org vdc network %q is not an isolated network", idOrName))
	}

	ep := endpoints.DeleteVdcNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: %w", opDeleteVdcNetworkIsolated, err)
	}

	return nil
}
