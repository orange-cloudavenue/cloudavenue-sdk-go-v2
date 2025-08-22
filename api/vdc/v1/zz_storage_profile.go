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

// List of storage profiles. If no filters are applied, all Storage Profiles for all VDCs are displayed. You can filter by either ID or Name of the storage profile, and you can also filter by either VDC ID or VDC Name. You can combine a storage profile filter (ID or Name) with a VDC filter (ID or Name). If both ID and Name are provided, they must refer to the same object; otherwise, the result will be empty.
func (c *Client) ListStorageProfile(ctx context.Context, params types.ParamsListStorageProfile) (*types.ModelListStorageProfiles, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*types.ModelListStorageProfiles), nil
}
// Add one or more storage profiles to a specific VDC.
func (c *Client) AddStorageProfile(ctx context.Context, params types.ParamsAddStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Add").Run(ctx, c, params)
	return err
}
// Delete a storage profile from a given VDC. You cannot delete the default storage profile, the last remaining profile, or a non-empty profile.
func (c *Client) DeleteStorageProfile(ctx context.Context, params types.ParamsDeleteStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Delete").Run(ctx, c, params)
	return err
}

