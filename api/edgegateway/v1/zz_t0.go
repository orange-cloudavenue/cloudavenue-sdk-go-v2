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

// List all T0s available in the organization. This command retrieves a list of all T0s, which are the top-level network services in the Edge Gateway architecture.
func (c *Client) ListT0(ctx context.Context) (*types.ModelT0s, error) {
	x, err := cmds.Get("T0", "", "List").Run(ctx, c, nil)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelT0s), nil
}

// Retrieve a specific T0 directly by its name or by the edge gateway it is associated with. This command allows you to fetch detailed information about a specific T0.
func (c *Client) GetT0(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	x, err := cmds.Get("T0", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelT0), nil
}
