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

<<<<<<< HEAD
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)


func (c *Client) GetBandwidth(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayBandwidth, error) {
=======

func (c *Client) GetBandwidth(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGatewayBandwidth, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "Bandwidth", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGatewayBandwidth), nil
=======
    return x.(*ModelEdgeGatewayBandwidth), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}








