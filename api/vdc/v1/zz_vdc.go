/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// List all Virtual Data Centers (VDCs) available in your organization. If no filters are applied, it returns all VDCs.
func (c *Client) ListVDC(ctx context.Context, params types.ParamsListVDC) (*types.ModelListVDC, error) {
	x, err := cmds.Get("VDC", "", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListVDC), nil
}

// Retrieve detailed information about a specific Virtual Data Center (VDC) by its name.
func (c *Client) GetVDC(ctx context.Context, params types.ParamsGetVDC) (*types.ModelGetVDC, error) {
	x, err := cmds.Get("VDC", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVDC), nil
}

// Create a new Virtual Data Center (VDC) with the specified parameters.
func (c *Client) CreateVDC(ctx context.Context, params types.ParamsCreateVDC) (*types.ModelGetVDC, error) {
	x, err := cmds.Get("VDC", "", "Create").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVDC), nil
}

// Update VDC performs a PUT request to update an existing VDC. Enter only the fields you want to update.
func (c *Client) UpdateVDC(ctx context.Context, params types.ParamsUpdateVDC) (*types.ModelGetVDC, error) {
	x, err := cmds.Get("VDC", "", "Update").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVDC), nil
}

// Delete VDC performs a DELETE request to delete an existing VDC.
func (c *Client) DeleteVDC(ctx context.Context, params types.ParamsDeleteVDC) error {
	_, err := cmds.Get("VDC", "", "Delete").Run(ctx, c, params)
	return err
}
