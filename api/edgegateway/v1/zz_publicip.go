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

// This command allows you to create a new Public IP in the Edge Gateway.
func (c *Client) CreatePublicIP(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayPublicIP, error) {
	x, err := cmds.Get("EdgeGateway", "PublicIP", "Create").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelEdgeGatewayPublicIP), nil
}
// This command allows you to list all Public IPs in the Edge Gateway.
func (c *Client) ListPublicIP(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayPublicIPs, error) {
	x, err := cmds.Get("EdgeGateway", "PublicIP", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelEdgeGatewayPublicIPs), nil
}
// This command allows you to retrieve information about a Public IP in the Edge Gateway.
func (c *Client) GetPublicIP(ctx context.Context, params types.ParamsGetEdgeGatewayPublicIP) (*types.ModelEdgeGatewayPublicIP, error) {
	x, err := cmds.Get("EdgeGateway", "PublicIP", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelEdgeGatewayPublicIP), nil
}
// This command allows you to delete a Public IP in the Edge Gateway.
func (c *Client) DeletePublicIP(ctx context.Context, params types.ParamsDeleteEdgeGatewayPublicIP) error {
	_, err := cmds.Get("EdgeGateway", "PublicIP", "Delete").Run(ctx, c, params)
	return err
}

