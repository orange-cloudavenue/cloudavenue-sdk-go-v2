/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vapp

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetVAppLeaseSettings    = "VAppLeaseSettings.Get"
	opUpdateVAppLeaseSettings = "VAppLeaseSettings.Update"
)

// GetVAppLeaseSettings retrieves org/VDC-level lease settings.
func (c *Client) GetVAppLeaseSettings(ctx context.Context, params types.ParamsGetVAppLeaseSettings) (*types.ModelVAppLeaseSettings, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetVAppLeaseSettings, err)
	}

	if params.OrgID == "" {
		return nil, fmt.Errorf("%s: org id is required", opGetVAppLeaseSettings)
	}

	ep := endpoints.GetVAppLeaseSettings()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.OrgID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetVAppLeaseSettings, err)
	}

	lease, ok := resp.Result().(*itypes.APIResponseOrgLeaseSettings)
	if !ok || lease == nil {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetVAppLeaseSettings, resp.Result())
	}

	return &types.ModelVAppLeaseSettings{
		DeploymentLeaseInSeconds: lease.DeploymentLeaseInSeconds,
		StorageLeaseInSeconds:    lease.StorageLeaseInSeconds,
	}, nil
}

// UpdateVAppLeaseSettings updates org/VDC-level lease settings.
func (c *Client) UpdateVAppLeaseSettings(ctx context.Context, params types.ParamsUpdateVAppLeaseSettings) (*types.ModelVAppLeaseSettings, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opUpdateVAppLeaseSettings, err)
	}

	if params.OrgID == "" {
		return nil, fmt.Errorf("%s: org id is required", opUpdateVAppLeaseSettings)
	}

	body := itypes.APIRequestOrgLeaseSettings{
		DeploymentLeaseInSeconds: params.LeaseSettings.DeploymentLeaseInSeconds,
		StorageLeaseInSeconds:    params.LeaseSettings.StorageLeaseInSeconds,
	}

	ep := endpoints.UpdateVAppLeaseSettings()
	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.OrgID), cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opUpdateVAppLeaseSettings, err)
	}

	lease, ok := resp.Result().(*itypes.APIResponseOrgLeaseSettings)
	if !ok || lease == nil {
		return nil, fmt.Errorf("%s: unexpected response type %T", opUpdateVAppLeaseSettings, resp.Result())
	}

	return &types.ModelVAppLeaseSettings{
		DeploymentLeaseInSeconds: lease.DeploymentLeaseInSeconds,
		StorageLeaseInSeconds:    lease.StorageLeaseInSeconds,
	}, nil
}
