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

	vdcapi "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const opListVDCNetwork = "VApp.VDCNetwork.List"

// ListVDCNetwork lists the Org VDC networks owned by a VDC or VDC group,
// fetching every page. The caller provides either OwnerID (full URN or bare
// UUID) or OwnerName; when OwnerName is set and OwnerID is empty, the name is
// resolved through the VDC listing endpoint.
func (c *Client) ListVDCNetwork(ctx context.Context, params types.ParamsListVDCNetwork) (*types.ModelListVDCNetwork, error) {
	vdcClient, err := vdcapi.New(c.c)
	if err != nil {
		return nil, fmt.Errorf("%s: create vdc client: %w", opListVDCNetwork, err)
	}

	result, err := vdcClient.ListVDCNetwork(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opListVDCNetwork, err)
	}

	return result, nil
}
