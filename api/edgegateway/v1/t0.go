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

	"resty.dev/v3"

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
	resp, err := c.c.Do(ctx, endpoints.ListT0())
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

	if params.EdgegatewayID == "" && params.EdgegatewayName == "" {
		return nil, fmt.Errorf("edge gateway id or name is required when t0 name is not provided")
	}

	if err := validateEdgeGatewayRef(params.EdgegatewayName); err != nil {
		return nil, err
	}

	return c.getT0ByEdgeGateway(ctx, params)
}

func (c *Client) getT0ByName(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	return c.getT0ByNameWithHierarchy(ctx, params, nil)
}

func (c *Client) getT0ByNameWithHierarchy(ctx context.Context, params types.ParamsGetT0, hierarchy *itypes.APIResponseNetworkServices) (*types.ModelT0, error) {
	ep := cav.MustGetEndpoint("GetT0")

	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], params.T0Name))
	if err != nil {
		return nil, wrapT0APIError("GetT0", err, resp)
	}

	details := resp.Result().(*itypes.APIResponseT0)
	if details == nil || details.Tier0VRF == "" {
		return nil, newT0NotFoundError(params, resp)
	}

	if hierarchy == nil {
		hierarchyResp, hErr := c.c.Do(ctx, endpoints.GetEdgeGatewayServices())
		if hErr != nil {
			return nil, wrapT0APIError("GetT0", hErr, hierarchyResp)
		}
		hierarchy = hierarchyResp.Result().(*itypes.APIResponseNetworkServices)
	}

	return details.ToModel(hierarchy), nil
}

func (c *Client) getT0ByEdgeGateway(ctx context.Context, params types.ParamsGetT0) (*types.ModelT0, error) {
	resp, err := c.c.Do(ctx, endpoints.GetEdgeGatewayServices())
	if err != nil {
		return nil, wrapT0APIError("GetT0", err, resp)
	}

	hierarchy := resp.Result().(*itypes.APIResponseNetworkServices)

	t0Name := findT0NameForEdgeGateway(hierarchy, params)
	if t0Name == "" {
		return nil, newT0NotFoundError(params, resp)
	}

	return c.getT0ByNameWithHierarchy(ctx, types.ParamsGetT0{T0Name: t0Name}, hierarchy)
}

func findT0NameForEdgeGateway(t0s *itypes.APIResponseNetworkServices, params types.ParamsGetT0) string {
	edgeID := params.EdgegatewayID

	for _, t0 := range *t0s {
		if t0.Name == "" {
			continue
		}

		for _, edgeGateway := range t0.Children {
			if edgeGateway.Type != nodeTypeEdgeGateway {
				continue
			}

			if edgeID != "" && edgeID == edgeGateway.Properties.EdgeUUID {
				return t0.Name
			}

			if params.EdgegatewayName != "" && params.EdgegatewayName == edgeGateway.Name {
				return t0.Name
			}
		}
	}

	return ""
}

func newT0NotFoundError(params types.ParamsGetT0, resp *resty.Response) error {
	apiErr := &errors.APIError{
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
	if resp != nil {
		apiErr.Duration = resp.Duration()
		if resp.Request != nil {
			apiErr.Endpoint = resp.Request.URL
			apiErr.Method = resp.Request.Method
		}
	}

	return apiErr
}

func wrapT0APIError(operation string, err error, resp *resty.Response) error {
	if err == nil {
		return nil
	}

	var apiErr *errors.APIError
	if !errors.As(err, &apiErr) {
		return fmt.Errorf("error getting T0: %w", err)
	}

	if resp != nil {
		apiErr.Duration = resp.Duration()
		if resp.Request != nil {
			apiErr.Endpoint = resp.Request.URL
			apiErr.Method = resp.Request.Method
		}
	}

	if apiErr.Operation == "" {
		apiErr.Operation = operation
	}

	return apiErr
}
