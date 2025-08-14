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


func (c *Client) GetServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayServices, error) {
=======

func (c *Client) GetServices(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGatewayServices, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "Services", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelEdgeGatewayServices), nil
=======
    return x.(*ModelEdgeGatewayServices), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}








<<<<<<< HEAD
func (c *Client) GetCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelCloudavenueServices, error) {
=======
func (c *Client) GetCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) (*ModelCloudavenueServices, error) {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    x, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
<<<<<<< HEAD
    return x.(*types.ModelCloudavenueServices), nil
=======
    return x.(*ModelCloudavenueServices), nil
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
}










<<<<<<< HEAD
func (c *Client) EnableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
=======
func (c *Client) EnableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    _, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Enable").Run(ctx, c, params)
    return err
}








<<<<<<< HEAD
func (c *Client) DisableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
=======
func (c *Client) DisableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
>>>>>>> 5f13790 (feat: Add List Storage Profile query)
    _, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Disable").Run(ctx, c, params)
    return err
}






