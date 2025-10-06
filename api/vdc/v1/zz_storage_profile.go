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

// Retrieves a comprehensive list of storage profiles for a specific VDC.
func (c *Client) ListStorageProfile(ctx context.Context, params types.ParamsListStorageProfile) (*types.ModelListStorageProfiles, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListStorageProfiles), nil
}
// Retrieves details of a specific storage profile within a VDC, identified by its unique ID or class name.
func (c *Client) GetStorageProfile(ctx context.Context, params types.ParamsGetStorageProfile) (*types.ModelGetStorageProfile, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "Get").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelGetStorageProfile), nil
}
// Creates one or more storage profiles within a specified VDC. Each profile requires a storage class and capacity limit, with an optional default designation.
func (c *Client) AddStorageProfile(ctx context.Context, params types.ParamsAddStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Add").Run(ctx, c, params)
	return err
}
// Removes a storage profile from the specified VDC. Deletion is restricted for default profiles, the last remaining profile, or profiles currently in use.
func (c *Client) DeleteStorageProfile(ctx context.Context, params types.ParamsDeleteStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Delete").Run(ctx, c, params)
	return err
}
// Modifies one or more storage profiles within a VDC. Supported updates include capacity limits and default profile designation. Storage class names cannot be modified.
func (c *Client) UpdateStorageProfile(ctx context.Context, params types.ParamsUpdateStorageProfile) (*types.ModelListStorageProfilesVDC, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "Update").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListStorageProfilesVDC), nil
}

