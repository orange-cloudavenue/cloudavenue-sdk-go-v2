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


func (c *Client) GetEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGateway, error) {
=======

func (c *Client) GetEdgeGateway(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGateway, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGateway), nil
=======
    return x.(*ModelEdgeGateway), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}














<<<<<<< HEAD
func (c *Client) ListEdgeGateway(ctx context.Context) (*types.ModelEdgeGateways, error) {
=======
func (c *Client) ListEdgeGateway(ctx context.Context) (*ModelEdgeGateways, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "", "List").Run(ctx, c, nil)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGateways), nil
}


func (c *Client) CreateEdgeGateway(ctx context.Context, params types.ParamsCreateEdgeGateway) (*types.ModelEdgeGateway, error) {
=======
    return x.(*ModelEdgeGateways), nil
}


func (c *Client) CreateEdgeGateway(ctx context.Context, params ParamsCreateEdgeGateway) (*ModelEdgeGateway, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "", "Create").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGateway), nil
=======
    return x.(*ModelEdgeGateway), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}










<<<<<<< HEAD
func (c *Client) DeleteEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) error {
=======
func (c *Client) DeleteEdgeGateway(ctx context.Context, params ParamsEdgeGateway) error {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    _, err := cmds.Get("EdgeGateway", "", "Delete").Run(ctx, c, params)
    return err
}






<<<<<<< HEAD
func (c *Client) UpdateEdgeGateway(ctx context.Context, params types.ParamsUpdateEdgeGateway) (*types.ModelEdgeGateway, error) {
=======
func (c *Client) UpdateEdgeGateway(ctx context.Context, params ParamsUpdateEdgeGateway) (*ModelEdgeGateway, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "", "Update").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGateway), nil
=======
    return x.(*ModelEdgeGateway), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}








