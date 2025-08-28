/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// List all OnPremise IPs allowed allowed for this organization's draas offer
func (c *Client) ListOnPremiseIp(ctx context.Context) (*types.ModelListDraasOnPremise, error) {
	x, err := cmds.Get("Draas", "OnPremiseIp", "List").Run(ctx, c, nil)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListDraasOnPremise), nil
}
// Add a new OnPremise IP (only IPV4) address to this organization's draas offer
func (c *Client) AddOnPremiseIp(ctx context.Context, params types.ParamsAddDraasOnPremiseIP) error {
	_, err := cmds.Get("Draas", "OnPremiseIp", "Add").Run(ctx, c, params)
	return err
}
// Remove an existing OnPremise IP address from this organization's draas offer
func (c *Client) RemoveOnPremiseIp(ctx context.Context, params types.ParamsRemoveDraasOnPremiseIP) error {
	_, err := cmds.Get("Draas", "OnPremiseIp", "Remove").Run(ctx, c, params)
	return err
}

