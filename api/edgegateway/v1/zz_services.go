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

// Retrieve services information about a specific EdgeGateway. This command retrieves the network services available on the EdgeGateway, such as load balancers, public IPs, and Cloud Avenue services.
func (c *Client) GetServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayServices, error) {
	x, err := cmds.Get("EdgeGateway", "Services", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelEdgeGatewayServices), nil
}

// Retrieve Cloud Avenue services on an EdgeGateway. This command returns the Cloud Avenue services available on the EdgeGateway, such as DNS, DHCP, and others.
func (c *Client) GetCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelCloudavenueServices, error) {
	x, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelCloudavenueServices), nil
}

// Enable Cloud Avenue services on an EdgeGateway.
func (c *Client) EnableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
	_, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Enable").Run(ctx, c, params)
	return err
}

// Disable Cloud Avenue services on an EdgeGateway. Cloudavenue services is a network setting that allows the EdgeGateway to connect to the mutualized Cloud Avenue services (DNS, DHCP, etc.).
func (c *Client) DisableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
	_, err := cmds.Get("EdgeGateway", "CloudavenueServices", "Disable").Run(ctx, c, params)
	return err
}
