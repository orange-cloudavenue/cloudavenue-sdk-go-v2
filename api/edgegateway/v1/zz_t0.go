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








func (c *Client) ListT0(ctx context.Context) (*types.ModelT0s, error) {
    x, err := cmds.Get("T0", "", "List").Run(ctx, c, nil)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelT0s), nil
}


func (c *Client) GetT0(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
    x, err := cmds.Get("T0", "", "Get").Run(ctx, c, params)
    if err != nil {
        return nil, err
    }
    return x.(*types.ModelT0), nil
}








