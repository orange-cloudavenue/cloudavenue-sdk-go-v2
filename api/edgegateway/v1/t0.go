/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// ListT0 lists T0 routers visible to organization.
func (c *Client) ListT0(ctx context.Context) (*types.ModelT0s, error) {
	ep := endpoints.ListT0()

	resp, err := c.c.Do(ctx, ep)
	if err != nil {
		return nil, fmt.Errorf("error listing T0s: %w", err)
	}

	return resp.Result().(*itypes.APIResponseT0s).ToModel(), nil
}

// GetT0 gets T0 router by name or by attached edge gateway.
func (c *Client) GetT0(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	ep := endpoints.ListT0()

	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], params.T0Name),
		cav.WithQueryParam(ep.QueryParams[1], params.EdgegatewayName),
		cav.WithQueryParam(ep.QueryParams[2], params.EdgegatewayID),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting T0: %w", err)
	}

	t0s := resp.Result().(*itypes.APIResponseT0s).ToModel()
	var t0 *types.ModelT0

	for _, t := range t0s.T0s {
		if params.T0Name != "" && t.Name == params.T0Name {
			t0 = &t
			break
		}
		if params.EdgegatewayID != "" || params.EdgegatewayName != "" {
			for _, edgeGateway := range t.EdgeGateways {
				if params.EdgegatewayID == edgeGateway.ID || params.EdgegatewayName == edgeGateway.Name {
					t0 = &t
					break
				}
			}
		}
	}

	if t0 == nil {
		return nil, &errors.APIError{
			Operation:     "GetT0",
			StatusCode:    http.StatusNotFound,
			StatusMessage: http.StatusText(http.StatusNotFound),
			Message: func() string {
				if params.T0Name != "" {
					return fmt.Sprintf("T0 with name %s not found", params.T0Name)
				}
				if params.EdgegatewayID != "" {
					return fmt.Sprintf("T0 for edge gateway with ID %s not found", params.EdgegatewayID)
				}
				return fmt.Sprintf("T0 for edge gateway with name %s not found", params.EdgegatewayName)
			}(),
			Duration: resp.Duration(),
			Endpoint: resp.Request.URL,
			Method:   resp.Request.Method,
		}
	}

	return t0, nil
}
