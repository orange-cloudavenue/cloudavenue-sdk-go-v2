/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package organization

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// Retrieve detailed information about your organization.
func (c *Client) GetOrganization(ctx context.Context) (*types.ModelGetOrganization, error) {
	x, err := cmds.Get("Organization", "", "Get").Run(ctx, c, nil)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetOrganization), nil
}

// Update the details of an existing organization.
func (c *Client) UpdateOrganization(ctx context.Context, params types.ParamsUpdateOrganization) (*types.ModelGetOrganization, error) {
	x, err := cmds.Get("Organization", "", "Update").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetOrganization), nil
}
