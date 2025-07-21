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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/urn"
	"github.com/orange-cloudavenue/common-go/validators"
)

func (c *Client) GetNetworkServices(ctx context.Context, params ParamsEdgeGateway) (*ModelNetworkServicesSvcs, error) {
	if err := validators.New().Struct(&params); err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Get the endpoint for the edge gateway services
	// Error is ignored here because the endpoint is registered at package init time.
	ep, _ := cav.GetEndpoint("NetworkServices", cav.MethodGET)

	// Get network services
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], params.ID),   // Only for filtering the response
		cav.WithQueryParam(ep.QueryParams[1], params.Name), // Only for filtering the response
	)
	if err != nil {
		return nil, fmt.Errorf("error getting network services: %w", err)
	}

	originalResponse := resp.Result().(*apiResponseNetworkServices)
	if len(*originalResponse) == 0 {
		return nil, fmt.Errorf("no network services found")
	}

	return originalResponse.toModel(params), nil
}

func (c *Client) EnableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
	if err := validators.New().Struct(&params); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// ID is required to request the API.
	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return err
		}
	}

	// Get the endpoint for the edge gateway services
	// Error is ignored here because the endpoint is registered at package init time.
	ep, _ := cav.GetEndpoint("NetworkServices", cav.MethodPOST)

	// Prepare the request body
	requestBody := &apiRequestNetworkServicesCavSvc{
		NetworkType:   "cav-services",
		EdgeGatewayID: urn.ExtractUUID(params.ID),
	}

	// Validate apiRequestNetworkServicesCavSvc.
	// This will ensure default values are set.
	if err := validators.New().Struct(requestBody); err != nil {
		return err
	}

	// Enable network services
	_, err := c.c.Do(
		ctx,
		ep,
		cav.SetBody(requestBody),
	)
	if err != nil {
		return fmt.Errorf("error enabling network services: %w", err)
	}

	return nil
}

func (c *Client) DisableCloudavenueServices(ctx context.Context, params ParamsEdgeGateway) error {
	// Ensure the edge gateway exists and retrieve its services
	nSvc, err := c.GetNetworkServices(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get network services: %w", err)
	}

	// Get the endpoint for the edge gateway services
	// Error is ignored here because the endpoint is registered at package init time.
	ep, _ := cav.GetEndpoint("NetworkServices", cav.MethodDELETE)

	// Disable network services
	_, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], nSvc.Service.ID),
	)
	if err != nil {
		return fmt.Errorf("error disabling network services: %w", err)
	}

	return nil
}
