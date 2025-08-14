/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"context"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func (c *Client) ListVdcGroup(ctx context.Context, params types.ParamsListVdcGroup) (*types.ModelListVdcGroup, error) {
	x, err := cmds.Get("VdcGroup", "", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListVdcGroup), nil
}

func (c *Client) GetVdcGroup(ctx context.Context, params types.ParamsGetVdcGroup) (*types.ModelGetVdcGroup, error) {
	x, err := cmds.Get("VdcGroup", "", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVdcGroup), nil
}

func (c *Client) CreateVdcGroup(ctx context.Context, params types.ParamsCreateVdcGroup) (*types.ModelGetVdcGroup, error) {
	x, err := cmds.Get("VdcGroup", "", "Create").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVdcGroup), nil
}

func (c *Client) UpdateVdcGroup(ctx context.Context, params types.ParamsUpdateVdcGroup) (*types.ModelGetVdcGroup, error) {
	x, err := cmds.Get("VdcGroup", "", "Update").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetVdcGroup), nil
}

func (c *Client) DeleteVdcGroup(ctx context.Context, params types.ParamsDeleteVdcGroup) error {
	_, err := cmds.Get("VdcGroup", "", "Delete").Run(ctx, c, params)
	return err
}

func (c *Client) AddVdcToVdcGroup(ctx context.Context, params types.ParamsAddVdcToVdcGroup) error {
	_, err := cmds.Get("VdcGroup", "Vdc", "Add").Run(ctx, c, params)
	return err
}

func (c *Client) RemoveVdcFromVdcGroup(ctx context.Context, params types.ParamsRemoveVdcFromVdcGroup) error {
	_, err := cmds.Get("VdcGroup", "Vdc", "Remove").Run(ctx, c, params)
	return err
}
