/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"context"
	"fmt"

	vdcv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const opListVDCNetwork = "VApp.VDCNetwork.List"

// ListVDCNetwork lists the Org VDC networks owned by a VDC.
// The caller provides either VDCID (URN) or VDCName; when VDCName is set and
// VDCID is empty the name is resolved through the VDC listing endpoint.
func (c *Client) ListVDCNetwork(ctx context.Context, params types.ParamsListVDCNetwork) (*types.ModelListVDCNetwork, error) {
	vdcID := params.OwnerID
	if vdcID == "" {
		if params.OwnerName == "" {
			return nil, fmt.Errorf("%s: validate: vdc id or vdc name is required", opListVDCNetwork)
		}

		vdcClient, err := vdcv1.New(c.c)
		if err != nil {
			return nil, fmt.Errorf("%s: create vdc client: %w", opListVDCNetwork, err)
		}

		vdcs, err := vdcClient.ListVDC(ctx, types.ParamsListVDC{Name: params.OwnerName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc: %w", opListVDCNetwork, err)
		}
		if len(vdcs.VDCS) == 0 {
			return nil, fmt.Errorf("%s: vdc %q not found", opListVDCNetwork, params.OwnerName)
		}
		vdcID = vdcs.VDCS[0].ID
	}

	ep := endpoints.ListVDCNetwork()
	resp, err := c.c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("ownerRef.id==%s", vdcID)))
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListVDCNetwork, err)
	}

	list, ok := resp.Result().(*itypes.APIResponseListVDCNetwork)
	if !ok || list == nil {
		return nil, fmt.Errorf("%s: unexpected list response type %T", opListVDCNetwork, resp.Result())
	}

	return list.ToModel(), nil
}
