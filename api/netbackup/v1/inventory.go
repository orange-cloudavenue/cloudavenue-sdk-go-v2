/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package netbackup

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

const opListInventory = "NetBackup.ListInventory"

// ListInventory returns the NetBackup inventory.
func (c *Client) ListInventory(ctx context.Context) ([]itypes.APIResponseNetbackupInventory, error) {
	resp, err := cav.Execute(ctx, c.c, cav.Operation[struct{}, []itypes.APIResponseNetbackupInventory]{
		Name:     opListInventory,
		Endpoint: endpoints.ListNetbackupInventory(),
		Extract: func(resp *cav.Response, _ struct{}) ([]itypes.APIResponseNetbackupInventory, error) {
			list, ok := resp.Result().(*itypes.APIResponseListNetbackupInventory)
			if !ok || list == nil {
				return nil, fmt.Errorf("unexpected list response type %T", resp.Result())
			}
			return *list, nil
		},
	}, struct{}{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
