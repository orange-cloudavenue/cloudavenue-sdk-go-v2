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
	opListVDCNetwork           = "VDCNetwork.List"
	opGetVDCNetworkIsolated    = "VdcNetworkIsolated.Get"
	opCreateVDCNetworkIsolated = "VdcNetworkIsolated.Create"
	opUpdateVDCNetworkIsolated = "VdcNetworkIsolated.Update"
	opDeleteVDCNetworkIsolated = "VdcNetworkIsolated.Delete"
)

// ListVDCNetwork lists routed and isolated VDC networks owned by a VDC group.
func (c *Client) ListVDCNetwork(ctx context.Context, params types.ParamsListVDCNetwork) (*types.ModelListVDCNetwork, error) {
	resp, err := listVDCNetworkModel(ctx, c.c, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListVDCNetwork, err)
	}

	return resp, nil
}

// GetVDCNetworkIsolated returns an isolated VDC network by ID or name.
func (c *Client) GetVDCNetworkIsolated(ctx context.Context, params types.ParamsGetVDCNetworkIsolated) (*types.ModelGetVDCNetwork, error) {
	resp, err := getVDCNetworkModel(ctx, c.c, params.ID, params.Name, params.VDCGroupID, params.VDCGroupName, types.VDCNetworkTypeIsolated)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetVDCNetworkIsolated, err)
	}

	return resp, nil
}

func createVDCNetworkIsolatedBody(ctx context.Context, c cav.Client, params types.ParamsCreateVDCNetworkIsolated) (itypes.APIRequestVDCNetwork, error) {
	vdcGroupID, vdcGroupName, err := resolveVDCGroupRef(ctx, c, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return itypes.APIRequestVDCNetwork{}, err
	}

	return itypes.APIRequestVDCNetwork{
		Name:        params.Name,
		Description: params.Description,
		NetworkType: types.VDCNetworkTypeIsolated,
		OwnerRef:    &itypes.APIObjectReference{ID: vdcGroupID, Name: vdcGroupName},
		Subnets: itypes.APIVDCNetworkSubnets{
			Values: []itypes.APIVDCNetworkSubnetValue{{
				Gateway:      params.Subnet.Gateway,
				PrefixLength: params.Subnet.PrefixLength,
				DNSServer1:   params.Subnet.DNSServer1,
				DNSServer2:   params.Subnet.DNSServer2,
				DNSSuffix:    params.Subnet.DNSSuffix,
				IPRanges: itypes.APIVDCNetworkIPRanges{
					Values: createVDCNetworkIPRanges(params.Subnet.IPRanges),
				},
			}},
		},
		GuestVLANTaggingAllowed: params.GuestVLANTaggingAllowed,
		Shared:                  new(true),
	}, nil
}

// CreateVDCNetworkIsolated creates an isolated VDC network in a VDC group.
func (c *Client) CreateVDCNetworkIsolated(ctx context.Context, params types.ParamsCreateVDCNetworkIsolated) (*types.ModelGetVDCNetwork, error) {
	body, err := createVDCNetworkIsolatedBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateVDCNetworkIsolated, err)
	}

	ep := endpoints.CreateVDCNetwork()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVDCNetworkIsolated, err)
	}

	network, ok := resp.Result().(*itypes.APIResponseVDCNetwork)
	if !ok || network == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateVDCNetworkIsolated, resp.Result())
	}

	model := network.ToModel()
	return &model, nil
}

// mergeVDCNetworkIsolatedUpdate applies optional update params onto the current network state.
func mergeVDCNetworkIsolatedUpdate(current *itypes.APIResponseVDCNetwork, params types.ParamsUpdateVDCNetworkIsolated) itypes.APIRequestVDCNetwork {
	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	subnets := current.Subnets
	if params.Subnet != nil {
		subnets = itypes.APIVDCNetworkSubnets{
			Values: []itypes.APIVDCNetworkSubnetValue{{
				Gateway:      params.Subnet.Gateway,
				PrefixLength: params.Subnet.PrefixLength,
				DNSServer1:   params.Subnet.DNSServer1,
				DNSServer2:   params.Subnet.DNSServer2,
				DNSSuffix:    params.Subnet.DNSSuffix,
				IPRanges:     itypes.APIVDCNetworkIPRanges{Values: createVDCNetworkIPRanges(params.Subnet.IPRanges)},
			}},
		}
	}

	guestVLANTaggingAllowed := current.GuestVLANTaggingAllowed
	if params.GuestVLANTaggingAllowed != nil {
		guestVLANTaggingAllowed = params.GuestVLANTaggingAllowed
	}

	return itypes.APIRequestVDCNetwork{
		ID:                      current.ID,
		Name:                    current.Name,
		Description:             description,
		NetworkType:             current.NetworkType,
		OwnerRef:                current.OwnerRef,
		Subnets:                 subnets,
		GuestVLANTaggingAllowed: guestVLANTaggingAllowed,
		Shared:                  current.Shared,
	}
}

// UpdateVDCNetworkIsolated updates an isolated VDC network in a VDC group.
func (c *Client) UpdateVDCNetworkIsolated(ctx context.Context, params types.ParamsUpdateVDCNetworkIsolated) (*types.ModelGetVDCNetwork, error) {
	current, err := getVDCNetworkWithRetry(ctx, c.c, params.ID, "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVDCNetworkIsolated, err)
	}

	if current.NetworkType != types.VDCNetworkTypeIsolated {
		return nil, fmt.Errorf("%s: %w", opUpdateVDCNetworkIsolated, errors.Newf("org vdc network %q is not an isolated network", params.ID))
	}

	body := mergeVDCNetworkIsolatedUpdate(current, params)

	ep := endpoints.UpdateVDCNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVDCNetworkIsolated, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVDCNetworkIsolated deletes an isolated VDC network from a VDC group.
func (c *Client) DeleteVDCNetworkIsolated(ctx context.Context, params types.ParamsDeleteVDCNetworkIsolated) error {
	current, err := deleteVDCNetworkTarget(ctx, c.c, params.ID, params.Name, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opDeleteVDCNetworkIsolated, err)
	}

	if current.NetworkType != types.VDCNetworkTypeIsolated {
		idOrName := params.ID
		if idOrName == "" {
			idOrName = params.Name
		}
		return fmt.Errorf("%s: %w", opDeleteVDCNetworkIsolated, errors.Newf("org vdc network %q is not an isolated network", idOrName))
	}

	ep := endpoints.DeleteVDCNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: %w", opDeleteVDCNetworkIsolated, err)
	}

	return nil
}
