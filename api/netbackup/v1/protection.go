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

const opProtectMachine = "NetBackup.ProtectMachine"

// ProtectMachine protects a NetBackup machine with the given protection level.
func (c *Client) ProtectMachine(ctx context.Context, machineID, protectionLevelID string) (*itypes.APIResponseNetbackupProtectMachine, error) {
	if machineID == "" {
		return nil, fmt.Errorf("%s: machine id is required", opProtectMachine)
	}
	if protectionLevelID == "" {
		return nil, fmt.Errorf("%s: protection level id is required", opProtectMachine)
	}

	ep := endpoints.ProtectNetbackupMachine()

	resp, err := cav.Execute(ctx, c.c, cav.Operation[struct{}, *itypes.APIResponseNetbackupProtectMachine]{
		Name:     opProtectMachine,
		Endpoint: ep,
		RequestOptions: func(_ struct{}) ([]cav.EndpointRequestOption, error) {
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], machineID),
			}, nil
		},
		Transform: func(_ struct{}) (any, error) {
			return itypes.APIRequestNetbackupProtectMachine{
				ProtectionLevelID: protectionLevelID,
			}, nil
		},
		Extract: func(resp *cav.Response, _ struct{}) (*itypes.APIResponseNetbackupProtectMachine, error) {
			result, ok := resp.Result().(*itypes.APIResponseNetbackupProtectMachine)
			if !ok || result == nil {
				return nil, fmt.Errorf("unexpected protect response type %T", resp.Result())
			}
			return result, nil
		},
	}, struct{}{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
