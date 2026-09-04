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

func isVdcNetworkURN(idOrName string) bool {
	return urn.IsNetwork(idOrName)
}

func getVdcNetworkWithRetry(ctx context.Context, c cav.Client, idOrName, vdcGroupID string) (*itypes.ApiResponseVdcNetwork, error) {
	const maxAttempts = 5

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 {
			time.Sleep(2 * time.Second)
		}

		network, err := getVdcNetwork(ctx, c, idOrName, vdcGroupID)
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

func getVdcNetwork(ctx context.Context, c cav.Client, idOrName, vdcGroupID string) (*itypes.ApiResponseVdcNetwork, error) {
	if isVdcNetworkURN(idOrName) {
		ep := endpoints.GetVdcNetwork()
		resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], idOrName))
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.ApiResponseVdcNetwork), nil
	}

	ep := endpoints.ListVdcNetwork()
	resp, err := c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("name==%s;ownerRef.id==%s", idOrName, vdcGroupID)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListVdcNetwork)
	if len(list.Values) == 0 {
		return nil, &errors.APIError{Operation: "GetVdcNetwork", StatusCode: 404, Message: fmt.Sprintf("vdc network %q not found", idOrName)}
	}
	if len(list.Values) > 1 {
		return nil, errors.Newf("multiple vdc networks found for %q", idOrName)
	}

	return &list.Values[0], nil
}

func resolveVdcNetworkLookup(ctx context.Context, c cav.Client, idOrName, vdcGroupID, vdcGroupName string) (string, error) {
	if isVdcNetworkURN(idOrName) {
		return "", nil
	}

	if vdcGroupID == "" && vdcGroupName != "" {
		ep := endpoints.ListVdcGroup()
		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], "name=="+vdcGroupName))
		if err != nil {
			return "", err
		}

		list := resp.Result().(*itypes.ApiResponseListVdcGroup)
		if len(list.Values) > 0 {
			vdcGroupID = list.Values[0].ID
		}
	}

	if vdcGroupID == "" {
		return "", errors.Newf("vdc_group_id or vdc_group_name is required when looking up an Org VDC Network by name")
	}

	return vdcGroupID, nil
}

func listVdcNetworkModel(ctx context.Context, c cav.Client, vdcGroupID, vdcGroupName string) (*types.ModelListVdcNetwork, error) {
	if vdcGroupID == "" {
		ep := endpoints.ListVdcGroup()
		resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], "name=="+vdcGroupName))
		if err != nil {
			return nil, err
		}

		list := resp.Result().(*itypes.ApiResponseListVdcGroup)
		if len(list.Values) > 0 {
			vdcGroupID = list.Values[0].ID
		}
	}

	ep := endpoints.ListVdcNetwork()
	resp, err := c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("ownerRef.id==%s", vdcGroupID)))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*itypes.ApiResponseListVdcNetwork).ToModel(), nil
}

func getVdcNetworkModel(ctx context.Context, c cav.Client, id, name, vdcGroupID, vdcGroupName, expectedType string) (*types.ModelGetVdcNetwork, error) {
	idOrName := id
	if idOrName == "" {
		idOrName = name
	}

	resolvedVdcGroupID, err := resolveVdcNetworkLookup(ctx, c, idOrName, vdcGroupID, vdcGroupName)
	if err != nil {
		return nil, err
	}

	network, err := getVdcNetworkWithRetry(ctx, c, idOrName, resolvedVdcGroupID)
	if err != nil {
		return nil, err
	}

	if network.NetworkType != expectedType {
		return nil, errors.Newf("org vdc network %q is not a %s network", idOrName, map[string]string{types.VdcNetworkTypeIsolated: "isolated", types.VdcNetworkTypeRouted: "routed"}[expectedType])
	}

	model := network.ToModel()
	return &model, nil
}

func createVdcNetworkIPRanges(ipRanges []types.ParamsVdcNetworkIPRange) []itypes.ApiVdcNetworkIPRangeValue {
	out := make([]itypes.ApiVdcNetworkIPRangeValue, 0, len(ipRanges))
	for _, ipRange := range ipRanges {
		out = append(out, itypes.ApiVdcNetworkIPRangeValue{StartAddress: ipRange.StartAddress, EndAddress: ipRange.EndAddress})
	}

	return out
}

func deleteVdcNetworkTarget(ctx context.Context, c cav.Client, id, name, vdcGroupID, vdcGroupName string) (*itypes.ApiResponseVdcNetwork, error) {
	idOrName := id
	if idOrName == "" {
		idOrName = name
	}

	resolvedVdcGroupID, err := resolveVdcNetworkLookup(ctx, c, idOrName, vdcGroupID, vdcGroupName)
	if err != nil && !isVdcNetworkURN(idOrName) {
		return nil, err
	}

	return getVdcNetworkWithRetry(ctx, c, idOrName, resolvedVdcGroupID)
}
