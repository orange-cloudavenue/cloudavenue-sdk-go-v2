/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import "context"

func (c *Client) ListStorageProfile(ctx context.Context, params ParamsListStorageProfiles) (*ModelListStorageProfiles, error) {
	x, err := cmds.Get("VDC", "StorageProfile", "List").Run(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return x.(*ModelListStorageProfiles), nil
}

func (c *Client) AddStorageProfile(ctx context.Context, params ParamsAddStorageProfile) error {
	_, err := cmds.Get("VDC", "StorageProfile", "Add").Run(ctx, c, params)
	return err
}
