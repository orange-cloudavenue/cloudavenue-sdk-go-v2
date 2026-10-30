/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
)

const opGetUserQuotas = "IAM.GetUserQuotas"

// GetUserQuotas reads the deployed/stored VM quotas of a user.
// CloudAPI accepts quotas on user create/update but never returns them, so
// the values are read from the legacy admin XML user record.
func (c *Client) GetUserQuotas(ctx context.Context, params ParamsGetUserQuotas) (*ModelUserQuotas, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: validate: %w", opGetUserQuotas, err)
	}

	ep := endpoints.GetAdminUser()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], extractUUID(params.ID)),
		cav.OverrideSetResult(new(itypes.AdminUser)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetUserQuotas, err)
	}

	user, ok := resp.Result().(*itypes.AdminUser)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected response type %T", opGetUserQuotas, resp.Result())
	}

	return &ModelUserQuotas{
		ID:              user.ID,
		Name:            user.Name,
		DeployedVMQuota: user.DeployedVMQuota,
		StoredVMQuota:   user.StoredVMQuota,
	}, nil
}

// extractUUID reduces an entity URN (urn:vcloud:<type>:<uuid>) to its UUID
// part, which the legacy admin API expects as path identifier. Plain values
// are returned unchanged.
func extractUUID(id string) string {
	if idx := strings.LastIndex(id, ":"); idx >= 0 && strings.HasPrefix(id, "urn:") {
		return id[idx+1:]
	}
	return id
}
