/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import "context"

func (c *Client) GetServices(ctx context.Context, params ParamsEdgeGateway) (*ModelEdgeGatewayServices, error) {
	x, err := cmds.Get("EdgeGateway", "Services", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelEdgeGatewayServices), nil
}

func (c *Client) GetCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) (*ModelCloudavenueServices, error) {
	x, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelCloudavenueServices), nil
}

func (c *Client) EnableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
	_, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Enable").Run(ctx, c, params)
	return err
}

func (c *Client) DisableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
	_, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Disable").Run(ctx, c, params)
	return err
}
