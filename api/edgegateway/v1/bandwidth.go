/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// GetBandwidth gets configured bandwidth for edge gateway.
func (c *Client) GetBandwidth(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayBandwidth, error) {
	pT0 := types.ParamsGetT0{
		EdgegatewayID:   params.ID,
		EdgegatewayName: params.Name,
	}

	t0, err := c.GetT0(ctx, pT0)
	if err != nil {
		return nil, err
	}

	var edgeGateway *types.ModelT0EdgeGateway
	for _, eg := range t0.EdgeGateways {
		if eg.ID == params.ID || eg.Name == params.Name {
			edgeGateway = &eg
			break
		}
	}

	bandwidth := &types.ModelEdgeGatewayBandwidth{
		ID:                     edgeGateway.ID,
		Name:                   edgeGateway.Name,
		Bandwidth:              edgeGateway.Bandwidth,
		AllowedBandwidthValues: edgeGateway.AllowedBandwidthValues,
	}

	return bandwidth, nil
}
