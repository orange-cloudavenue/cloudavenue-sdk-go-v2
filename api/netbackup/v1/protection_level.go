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

const (
	opGetProtectionLevel   = "NetBackup.GetProtectionLevel"
	opListProtectionLevels = "NetBackup.ListProtectionLevels"
)

// GetProtectionLevel returns a NetBackup protection level by ID.
func (c *Client) GetProtectionLevel(ctx context.Context, id string) (*itypes.APIResponseNetbackupProtectionLevel, error) {
	if id == "" {
		return nil, fmt.Errorf("%s: id is required", opGetProtectionLevel)
	}

	ep := endpoints.GetNetbackupProtectionLevel()

	resp, err := cav.Execute(ctx, c.c, cav.Operation[struct{}, *itypes.APIResponseNetbackupProtectionLevel]{
		Name:     opGetProtectionLevel,
		Endpoint: ep,
		RequestOptions: func(_ struct{}) ([]cav.EndpointRequestOption, error) {
			return []cav.EndpointRequestOption{
				cav.WithPathParam(ep.PathParams[0], id),
			}, nil
		},
		Extract: func(resp *cav.Response, _ struct{}) (*itypes.APIResponseNetbackupProtectionLevel, error) {
			result, ok := resp.Result().(*itypes.APIResponseNetbackupProtectionLevel)
			if !ok || result == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}
			return result, nil
		},
	}, struct{}{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListProtectionLevels returns all NetBackup protection levels.
func (c *Client) ListProtectionLevels(ctx context.Context) ([]itypes.APIResponseNetbackupProtectionLevel, error) {
	resp, err := cav.Execute(ctx, c.c, cav.Operation[struct{}, []itypes.APIResponseNetbackupProtectionLevel]{
		Name:     opListProtectionLevels,
		Endpoint: endpoints.ListNetbackupProtectionLevels(),
		Extract: func(resp *cav.Response, _ struct{}) ([]itypes.APIResponseNetbackupProtectionLevel, error) {
			list, ok := resp.Result().(*itypes.APIResponseListNetbackupProtectionLevels)
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
