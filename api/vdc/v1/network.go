/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetVDCNetwork  = "VDC.VDCNetwork.Get"
	opListVDCNetwork = "VDC.VDCNetwork.List"
)

// GetVDCNetwork returns an Org VDC Network by ID or by name within an owner
// reference. The lookup is type-agnostic: routed and isolated networks share the
// same cloudapi endpoints, and the owner is only an ownerRef.id filter, so a
// plain VDC URN is accepted without any VDC group semantics.
func (c *Client) GetVDCNetwork(ctx context.Context, params types.ParamsGetVDCNetwork) (*types.ModelGetVDCNetwork, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVDCNetwork, err)
	}

	if urn.IsNetwork(params.ID) {
		ep := endpoints.GetVDCNetwork()
		resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.ID))
		if err != nil {
			return nil, fmt.Errorf("%s: get: %w", opGetVDCNetwork, err)
		}

		network, ok := resp.Result().(*itypes.APIResponseVDCNetwork)
		if !ok || network == nil {
			return nil, fmt.Errorf("%s: unexpected response type %T", opGetVDCNetwork, resp.Result())
		}

		model := network.ToModel()
		inetworkobjects.LegacyVDCNetworkHREF(&model, c.c.GetConsole().GetAPIVCDEndpoint())
		return &model, nil
	}

	ownerID, err := c.resolveVDCNetworkOwner(ctx, params.OwnerID, params.OwnerName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVDCNetwork, err)
	}

	matches, err := inetworkobjects.ListVDCNetworkModels(ctx, c.c,
		fmt.Sprintf("name==%s;ownerRef.id==%s", params.Name, ownerID),
		c.c.GetConsole().GetAPIVCDEndpoint())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVDCNetwork, err)
	}

	switch len(matches) {
	case 0:
		return nil, errors.Newf("org vdc network %q not found", params.Name)
	case 1:
		return &matches[0], nil
	default:
		return nil, errors.Newf("multiple org vdc networks found for %q", params.Name)
	}
}

// ListVDCNetwork lists the Org VDC Networks owned by an owner reference,
// fetching every page.
func (c *Client) ListVDCNetwork(ctx context.Context, params types.ParamsListVDCNetwork) (*types.ModelListVDCNetwork, error) {
	ownerID, err := c.resolveVDCNetworkOwner(ctx, params.OwnerID, params.OwnerName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListVDCNetwork, err)
	}

	networks, err := inetworkobjects.ListVDCNetworkModels(ctx, c.c,
		fmt.Sprintf("ownerRef.id==%s", ownerID),
		c.c.GetConsole().GetAPIVCDEndpoint())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListVDCNetwork, err)
	}

	return &types.ModelListVDCNetwork{VDCNetworks: networks}, nil
}

// resolveVDCNetworkOwner normalizes an owner URN or resolves a VDC name to its
// URN through the VDC listing endpoint.
func (c *Client) resolveVDCNetworkOwner(ctx context.Context, ownerID, ownerName string) (string, error) {
	if ownerID != "" {
		return urn.Normalize(urn.VDC, ownerID).String(), nil
	}
	if ownerName == "" {
		return "", errors.Newf("owner id or owner name is required")
	}

	vdcs, err := c.ListVDC(ctx, types.ParamsListVDC{Name: ownerName})
	if err != nil {
		return "", fmt.Errorf("resolve vdc: %w", err)
	}
	if len(vdcs.VDCS) == 0 {
		return "", errors.Newf("vdc %q not found", ownerName)
	}

	// The VDC listing exposes a bare UUID; the orgVdcNetworks filter needs the
	// full ownerRef.id URN.
	return urn.Normalize(urn.VDC, vdcs.VDCS[0].ID).String(), nil
}
