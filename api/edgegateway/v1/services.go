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
	"strings"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// GetServices returns network service state for an edge gateway.
func (c *Client) GetServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGatewayServices, error) {
	if params.ID == "" && params.Name == "" {
		return nil, fmt.Errorf("id or name is required")
	}

	ep := endpoints.GetEdgeGatewayServices()

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, err
		}
	}

	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], params.ID),
		cav.WithQueryParam(ep.QueryParams[1], params.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving network services for edge gateway %s: %w", params.ID, err)
	}

	data := resp.Result().(*itypes.APIResponseNetworkServices).ToModel(params)
	if data == nil {
		return nil, fmt.Errorf("no network services found for edge gateway %s", params.ID)
	}

	return data, nil
}

// GetCloudavenueServices returns CloudAvenue-specific service state for an edge gateway.
func (c *Client) GetCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelCloudavenueServices, error) {
	svcs, err := c.GetServices(ctx, params)
	if err != nil {
		return nil, err
	}

	return svcs.Services, nil
}

// EnableCloudavenueServices enables CloudAvenue services for an edge gateway.
func (c *Client) EnableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
	if params.ID == "" && params.Name == "" {
		return fmt.Errorf("id or name is required")
	}

	ep := endpoints.EnableCloudavenueServices()

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return err
		}
	}

	requestBody := &itypes.APIRequestNetworkServicesCavSvc{
		NetworkType:   "cav-services",
		EdgeGatewayID: urn.ExtractUUID(params.ID),
		Properties: struct {
			PrefixLength int "json:\"prefixLength,omitempty\" validate:\"omitempty,min=25,max=28\" default:\"27\""
		}{
			PrefixLength: 27,
		},
	}

	_, err := c.c.Do(
		ctx,
		ep,
		cav.SetBody(requestBody),
	)
	if err != nil {
		if !strings.Contains(err.Error(), "subnet not fully consumed") {
			return fmt.Errorf("error enabling network services: %w", err)
		}
	}
	return nil
}

// DisableCloudavenueServices disables CloudAvenue services for an edge gateway.
func (c *Client) DisableCloudavenueServices(ctx context.Context, params types.ParamsEdgeGateway) error {
	ep := endpoints.DisableCloudavenueServices()

	nSvc, err := c.GetServices(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get network services: %w", err)
	}

	_, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], nSvc.Services.ID),
	)
	if err != nil {
		return fmt.Errorf("error disabling network services: %w", err)
	}
	return nil
}
