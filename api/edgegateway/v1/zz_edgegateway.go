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


func (c *Client) GetEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGateway, error) {
    x, err := cmds.Get("EdgeGateway", "", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelEdgeGateway), nil
}














func (c *Client) ListEdgeGateway(ctx context.Context) (*types.ModelEdgeGateways, error) {
    x, err := cmds.Get("EdgeGateway", "", "List").Run(ctx, c, nil)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelEdgeGateways), nil
}


func (c *Client) CreateEdgeGateway(ctx context.Context, params types.ParamsCreateEdgeGateway) (*types.ModelEdgeGateway, error) {
    x, err := cmds.Get("EdgeGateway", "", "Create").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelEdgeGateway), nil
}










func (c *Client) DeleteEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) error {
    _, err := cmds.Get("EdgeGateway", "", "Delete").Run(ctx, c, params)
    return err
}






func (c *Client) UpdateEdgeGateway(ctx context.Context, params types.ParamsUpdateEdgeGateway) (*types.ModelEdgeGateway, error) {
    x, err := cmds.Get("EdgeGateway", "", "Update").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelEdgeGateway), nil
}








