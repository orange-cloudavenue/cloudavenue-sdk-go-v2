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

// Get the bandwidth of an edge gateway. This command retrieves the bandwidth information for a specific edge gateway.
func (c *Client) GetBandwidth(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayBandwidth, error) {
	x, err := cmds.Get("EdgeGateway", "Bandwidth", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelEdgeGatewayBandwidth), nil
}

