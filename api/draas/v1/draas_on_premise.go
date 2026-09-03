/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package draas

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// ListOnPremiseIp lists on-premise IP addresses allowed for organization DRaaS offer.
func (c *Client) ListOnPremiseIp(ctx context.Context) (*types.ModelListDraasOnPremise, error) {
	ep := endpoints.ListDraasOnPremiseIp()

	resp, err := c.c.Do(
		ctx,
		ep,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list Draas OnPremiseIP: %w", err)
	}

	return resp.Result().(*itypes.ApiResponseListDraasOnPremise).ToModel(), nil
}

// AddOnPremiseIp adds on-premise IPv4 address to organization DRaaS offer.
func (c *Client) AddOnPremiseIp(ctx context.Context, params types.ParamsAddDraasOnPremiseIP) error {
	ep := endpoints.AddDraasOnPremiseIp()

	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.IP),
	)
	if err != nil {
		return fmt.Errorf("failed to add Draas OnPremiseIP: %w", err)
	}

	return nil
}

// RemoveOnPremiseIp removes on-premise IP address from organization DRaaS offer.
func (c *Client) RemoveOnPremiseIp(ctx context.Context, params types.ParamsRemoveDraasOnPremiseIP) error {
	ep := endpoints.RemoveDraasOnPremiseIp()

	_, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], params.IP),
	)
	if err != nil {
		return fmt.Errorf("failed to remove Draas OnPremiseIP: %w", err)
	}

	return nil
}
