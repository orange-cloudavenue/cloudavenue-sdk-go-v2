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
	stderrors "errors"
	"fmt"
	"time"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func isVDCNetworkURN(idOrName string) bool {
	return urn.IsNetwork(idOrName)
}

func getVDCNetworkWithRetry(ctx context.Context, c cav.Client, idOrName, vdcGroupID string) (*itypes.APIResponseVDCNetwork, error) {
	const maxAttempts = 5

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		network, err := getVDCNetwork(ctx, c, idOrName, vdcGroupID)
		if err == nil {
			return network, nil
		}

		var apiErr *errors.APIError
		if !stderrors.As(err, &apiErr) || !apiErr.IsNotFound() {
			return nil, err
		}

		lastErr = err
	}

	return nil, lastErr
}

func getVDCNetwork(ctx context.Context, c cav.Client, idOrName, vdcGroupID string) (*itypes.APIResponseVDCNetwork, error) {
	if isVDCNetworkURN(idOrName) {
		ep := endpoints.GetVDCNetwork()
		resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], idOrName))
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.APIResponseVDCNetwork), nil
	}

	ep := endpoints.ListVDCNetwork()
	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("name==%s;ownerRef.id==%s", idOrName, vdcGroupID)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.APIResponseListVDCNetwork)
	if len(list.Values) == 0 {
		return nil, &errors.APIError{Operation: "GetVDCNetwork", StatusCode: 404, Message: fmt.Sprintf("vdc network %q not found", idOrName)}
	}
	if len(list.Values) > 1 {
		return nil, errors.Newf("multiple vdc networks found for %q", idOrName)
	}

	return &list.Values[0], nil
}

func resolveVDCNetworkLookup(ctx context.Context, c cav.Client, idOrName, vdcGroupID, vdcGroupName string) (string, error) {
	if isVDCNetworkURN(idOrName) {
		return "", nil
	}

	if vdcGroupID == "" && vdcGroupName != "" {
		ep := endpoints.ListVDCGroup()
		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], "name=="+vdcGroupName))
		if err != nil {
			return "", err
		}

		list := resp.Result().(*itypes.APIResponseListVDCGroup)
		if len(list.Values) > 0 {
			vdcGroupID = list.Values[0].ID
		}
	}

	if vdcGroupID == "" {
		return "", errors.Newf("vdc_group_id or vdc_group_name is required when looking up an Org VDC Network by name")
	}

	return vdcGroupID, nil
}

func listVDCNetworkModel(ctx context.Context, c cav.Client, vdcGroupID, vdcGroupName string) (*types.ModelListVDCNetwork, error) {
	if vdcGroupID == "" {
		ep := endpoints.ListVDCGroup()
		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], "name=="+vdcGroupName))
		if err != nil {
			return nil, err
		}

		list := resp.Result().(*itypes.APIResponseListVDCGroup)
		if len(list.Values) > 0 {
			vdcGroupID = list.Values[0].ID
		}
	}

	ep := endpoints.ListVDCNetwork()
	resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("ownerRef.id==%s", vdcGroupID)))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*itypes.APIResponseListVDCNetwork).ToModel(), nil
}

func getVDCNetworkModel(ctx context.Context, c cav.Client, id, name, vdcGroupID, vdcGroupName, expectedType string) (*types.ModelGetVDCNetwork, error) {
	idOrName := id
	if idOrName == "" {
		idOrName = name
	}

	resolvedVDCGroupID, err := resolveVDCNetworkLookup(ctx, c, idOrName, vdcGroupID, vdcGroupName)
	if err != nil {
		return nil, err
	}

	network, err := getVDCNetworkWithRetry(ctx, c, idOrName, resolvedVDCGroupID)
	if err != nil {
		return nil, err
	}

	if network.NetworkType != expectedType {
		return nil, errors.Newf("org vdc network %q is not a %s network", idOrName, map[string]string{types.VDCNetworkTypeIsolated: "isolated", types.VDCNetworkTypeRouted: "routed"}[expectedType])
	}

	model := network.ToModel()
	return &model, nil
}

func createVDCNetworkIPRanges(ipRanges []types.ParamsVDCNetworkIPRange) []itypes.APIVDCNetworkIPRangeValue {
	out := make([]itypes.APIVDCNetworkIPRangeValue, 0, len(ipRanges))
	for _, ipRange := range ipRanges {
		out = append(out, itypes.APIVDCNetworkIPRangeValue{StartAddress: ipRange.StartAddress, EndAddress: ipRange.EndAddress})
	}

	return out
}

func deleteVDCNetworkTarget(ctx context.Context, c cav.Client, id, name, vdcGroupID, vdcGroupName string) (*itypes.APIResponseVDCNetwork, error) {
	idOrName := id
	if idOrName == "" {
		idOrName = name
	}

	resolvedVDCGroupID, err := resolveVDCNetworkLookup(ctx, c, idOrName, vdcGroupID, vdcGroupName)
	if err != nil && !isVDCNetworkURN(idOrName) {
		return nil, err
	}

	return getVDCNetworkWithRetry(ctx, c, idOrName, resolvedVDCGroupID)
}
