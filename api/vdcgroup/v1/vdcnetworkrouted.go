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
	opGetVdcNetworkRouted    = "VdcNetworkRouted.Get"
	opCreateVdcNetworkRouted = "VdcNetworkRouted.Create"
	opUpdateVdcNetworkRouted = "VdcNetworkRouted.Update"
	opDeleteVdcNetworkRouted = "VdcNetworkRouted.Delete"
)

// GetVdcNetworkRouted gets NSX-T routed Org VDC network by ID or name within VDC group.
func (c *Client) GetVdcNetworkRouted(ctx context.Context, params types.ParamsGetVdcNetworkRouted) (*types.ModelGetVdcNetwork, error) {
	resp, err := getVdcNetworkModel(ctx, c.c, params.ID, params.Name, params.VdcGroupID, params.VdcGroupName, types.VdcNetworkTypeRouted)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetVdcNetworkRouted, err)
	}

	return resp, nil
}

// CreateVdcNetworkRouted creates NSX-T routed Org VDC network within VDC group.
func (c *Client) CreateVdcNetworkRouted(ctx context.Context, params types.ParamsCreateVdcNetworkRouted) (*types.ModelGetVdcNetwork, error) {
	vdcGroupID, vdcGroupName, err := resolveVdcGroupRef(ctx, c.c, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateVdcNetworkRouted, err)
	}

	body := itypes.ApiRequestVdcNetwork{
		Name:        params.Name,
		Description: params.Description,
		NetworkType: types.VdcNetworkTypeRouted,
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
		Connection: &itypes.ApiVdcNetworkConnection{
			RouterRef:           itypes.ApiObjectReference{ID: params.EdgeGatewayID, Name: params.EdgeGatewayName},
			ConnectionTypeValue: "INTERNAL",
		},
	}

	ep := endpoints.CreateVdcNetwork()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVdcNetworkRouted, err)
	}

	network, ok := resp.Result().(*itypes.ApiResponseVdcNetwork)
	if !ok || network == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateVdcNetworkRouted, resp.Result())
	}

	model := network.ToModel()
	return &model, nil
}

// UpdateVdcNetworkRouted updates NSX-T routed Org VDC network within VDC group.
func (c *Client) UpdateVdcNetworkRouted(ctx context.Context, params types.ParamsUpdateVdcNetworkRouted) (*types.ModelGetVdcNetwork, error) {
	current, err := getVdcNetworkWithRetry(ctx, c.c, params.ID, "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVdcNetworkRouted, err)
	}

	if current.NetworkType != types.VdcNetworkTypeRouted {
		return nil, fmt.Errorf("%s: %w", opUpdateVdcNetworkRouted, errors.Newf("org vdc network %q is not a routed network", params.ID))
	}

	body := mergeVdcNetworkRoutedUpdate(current, params)

	ep := endpoints.UpdateVdcNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVdcNetworkRouted, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVdcNetworkRouted deletes NSX-T routed Org VDC network from VDC group.
func (c *Client) DeleteVdcNetworkRouted(ctx context.Context, params types.ParamsDeleteVdcNetworkRouted) error {
	current, err := deleteVdcNetworkTarget(ctx, c.c, params.ID, params.Name, params.VdcGroupID, params.VdcGroupName)
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opDeleteVdcNetworkRouted, err)
	}

	if current.NetworkType != types.VdcNetworkTypeRouted {
		idOrName := params.ID
		if idOrName == "" {
			idOrName = params.Name
		}
		return fmt.Errorf("%s: %w", opDeleteVdcNetworkRouted, errors.Newf("org vdc network %q is not a routed network", idOrName))
	}

	ep := endpoints.DeleteVdcNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: %w", opDeleteVdcNetworkRouted, err)
	}

	return nil
}

// mergeVdcNetworkRoutedUpdate applies optional update params onto current network state.
func mergeVdcNetworkRoutedUpdate(current *itypes.ApiResponseVdcNetwork, params types.ParamsUpdateVdcNetworkRouted) itypes.ApiRequestVdcNetwork {
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

	connection := current.Connection
	if params.EdgeGatewayID != "" {
		connection = &itypes.ApiVdcNetworkConnection{
			RouterRef:           itypes.ApiObjectReference{ID: params.EdgeGatewayID, Name: params.EdgeGatewayName},
			ConnectionTypeValue: "INTERNAL",
		}
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
		Connection:              connection,
	}
}
