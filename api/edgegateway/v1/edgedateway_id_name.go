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
	"fmt"
)

func (c *Client) retrieveEdgeGatewayIDByName(ctx context.Context, name string) (string, error) {
	resp, err := c.ListEdgeGateway(ctx)
	if err != nil {
		return "", err
	}
	for _, gateway := range resp.EdgeGateways {
		if gateway.Name == name {
			return gateway.ID, nil
		}
	}
	return "", fmt.Errorf("edge gateway %q not found", name)
}
