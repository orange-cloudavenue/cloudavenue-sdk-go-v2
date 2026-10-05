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

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

// nodeTypeEdgeGateway is the network hierarchy node type of an edge gateway.
const nodeTypeEdgeGateway = "edge-gateway"

// ListT0 lists T0 routers visible to organization.
func (c *Client) ListT0(ctx context.Context) (*types.ModelT0s, error) {
	ep := endpoints.ListT0()

	resp, err := c.c.Do(ctx, ep)
	if err != nil {
		return nil, fmt.Errorf("error listing T0s: %w", err)
	}

	return resp.Result().(*itypes.APIResponseT0Names).ToModel(), nil
}

// GetT0 gets T0 router by name or by attached edge gateway.
func (c *Client) GetT0(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	if params.T0Name != "" {
		return c.getT0ByName(ctx, params)
	}

	return c.getT0ByEdgeGateway(ctx, params)
}

func (c *Client) getT0ByName(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	ep := cav.MustGetEndpoint("GetT0")

	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.T0Name))
	if err != nil {
		return nil, fmt.Errorf("error getting T0: %w", err)
	}

	result := resp.Result().(*itypes.APIResponseT0).ToModel()
	if len(result.T0s) == 0 || result.T0s[0].Name == "" {
		return nil, newT0NotFoundError(params)
	}

	return &result.T0s[0], nil
}

func (c *Client) getT0ByEdgeGateway(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	resp, err := c.c.Do(ctx, endpoints.GetEdgeGatewayServices())
	if err != nil {
		return nil, fmt.Errorf("error getting T0: %w", err)
	}

	t0Name := findT0NameForEdgeGateway(resp.Result().(*itypes.APIResponseNetworkServices), params)
	if t0Name == "" {
		return nil, newT0NotFoundError(params)
	}

	return c.getT0ByName(ctx, types.ParamsGetT0{T0Name: t0Name})
}

func findT0NameForEdgeGateway(t0s *itypes.APIResponseNetworkServices, params types.ParamsGetT0) string {
	edgeID := urn.Normalize(urn.EdgeGateway, params.EdgegatewayID).String()

	for _, t0 := range *t0s {
		if t0.Name == "" {
			continue
		}

		for _, edgeGateway := range t0.Children {
			if edgeGateway.Type != nodeTypeEdgeGateway {
				continue
			}

			if edgeID != "" && edgeID == urn.Normalize(urn.EdgeGateway, edgeGateway.Properties.EdgeUUID).String() {
				return t0.Name
			}

			if params.EdgegatewayName != "" && params.EdgegatewayName == edgeGateway.Name {
				return t0.Name
			}
		}
	}

	return ""
}

func newT0NotFoundError(params types.ParamsGetT0) error {
	return &errors.APIError{
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
	}
}
