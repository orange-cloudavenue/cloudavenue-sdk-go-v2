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
	opGetVDCNetworkRouted    = "VdcNetworkRouted.Get"
	opCreateVDCNetworkRouted = "VdcNetworkRouted.Create"
	opUpdateVDCNetworkRouted = "VdcNetworkRouted.Update"
	opDeleteVDCNetworkRouted = "VdcNetworkRouted.Delete"
)

// GetVDCNetworkRouted gets NSX-T routed Org VDC network by ID or name within VDC group.
func (c *Client) GetVDCNetworkRouted(ctx context.Context, params types.ParamsGetVDCNetworkRouted) (*types.ModelGetVDCNetwork, error) {
	resp, err := getVDCNetworkModel(ctx, c.c, params.ID, params.Name, params.VDCGroupID, params.VDCGroupName, types.VDCNetworkTypeRouted)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetVDCNetworkRouted, err)
	}

	return resp, nil
}

// CreateVDCNetworkRouted creates NSX-T routed Org VDC network within VDC group.
func (c *Client) CreateVDCNetworkRouted(ctx context.Context, params types.ParamsCreateVDCNetworkRouted) (*types.ModelGetVDCNetwork, error) {
	vdcGroupID, vdcGroupName, err := resolveVDCGroupRef(ctx, c.c, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateVDCNetworkRouted, err)
	}

	body := itypes.APIRequestVDCNetwork{
		Name:        params.Name,
		Description: params.Description,
		NetworkType: types.VDCNetworkTypeRouted,
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
		Connection: &itypes.APIVDCNetworkConnection{
			RouterRef:           itypes.APIObjectReference{ID: params.EdgeGatewayID, Name: params.EdgeGatewayName},
			ConnectionTypeValue: "INTERNAL",
		},
	}

	ep := endpoints.CreateVDCNetwork()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateVDCNetworkRouted, err)
	}

	network, ok := resp.Result().(*itypes.APIResponseVDCNetwork)
	if !ok || network == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateVDCNetworkRouted, resp.Result())
	}

	model := network.ToModel()
	return &model, nil
}

// UpdateVDCNetworkRouted updates NSX-T routed Org VDC network within VDC group.
func (c *Client) UpdateVDCNetworkRouted(ctx context.Context, params types.ParamsUpdateVDCNetworkRouted) (*types.ModelGetVDCNetwork, error) {
	current, err := getVDCNetworkWithRetry(ctx, c.c, params.ID, "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVDCNetworkRouted, err)
	}

	if current.NetworkType != types.VDCNetworkTypeRouted {
		return nil, fmt.Errorf("%s: %w", opUpdateVDCNetworkRouted, errors.Newf("org vdc network %q is not a routed network", params.ID))
	}

	body := mergeVDCNetworkRoutedUpdate(current, params)

	ep := endpoints.UpdateVDCNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVDCNetworkRouted, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVDCNetworkRouted deletes NSX-T routed Org VDC network from VDC group.
func (c *Client) DeleteVDCNetworkRouted(ctx context.Context, params types.ParamsDeleteVDCNetworkRouted) error {
	current, err := deleteVDCNetworkTarget(ctx, c.c, params.ID, params.Name, params.VDCGroupID, params.VDCGroupName)
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opDeleteVDCNetworkRouted, err)
	}

	if current.NetworkType != types.VDCNetworkTypeRouted {
		idOrName := params.ID
		if idOrName == "" {
			idOrName = params.Name
		}
		return fmt.Errorf("%s: %w", opDeleteVDCNetworkRouted, errors.Newf("org vdc network %q is not a routed network", idOrName))
	}

	ep := endpoints.DeleteVDCNetwork()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID)); err != nil {
		return fmt.Errorf("%s: %w", opDeleteVDCNetworkRouted, err)
	}

	return nil
}

// mergeVDCNetworkRoutedUpdate applies optional update params onto current network state.
func mergeVDCNetworkRoutedUpdate(current *itypes.APIResponseVDCNetwork, params types.ParamsUpdateVDCNetworkRouted) itypes.APIRequestVDCNetwork {
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

	connection := current.Connection
	if params.EdgeGatewayID != "" {
		connection = &itypes.APIVDCNetworkConnection{
			RouterRef:           itypes.APIObjectReference{ID: params.EdgeGatewayID, Name: params.EdgeGatewayName},
			ConnectionTypeValue: "INTERNAL",
		}
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
		Connection:              connection,
	}
}
