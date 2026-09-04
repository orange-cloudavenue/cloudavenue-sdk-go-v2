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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opGetEdgeGateway    = "EdgeGateway.Get"
	opListEdgeGateway   = "EdgeGateway.List"
	opCreateEdgeGateway = "EdgeGateway.Create"
	opDeleteEdgeGateway = "EdgeGateway.Delete"
	opUpdateEdgeGateway = "EdgeGateway.Update"
)

type getEdgeGatewayByIDParams struct {
	ID string
}

type deleteEdgeGatewayByIDParams struct {
	ID string
}

type updateEdgeGatewayByIDParams struct {
	ID        string
	Bandwidth int
}

var (
	listEdgeGatewayOp = cav.Operation[struct{}, *types.ModelEdgeGateways]{
		Name:     opListEdgeGateway,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.ListEdgeGateway(),
		Extract: func(resp *cav.Response, _ struct{}) (*types.ModelEdgeGateways, error) {
			list, ok := resp.Result().(*itypes.APIResponseEdgegateways)
			if !ok || list == nil {
				return nil, fmt.Errorf("unexpected list response type %T", resp.Result())
			}

			return list.ToModel(), nil
		},
	}
	getEdgeGatewayByIDOp = cav.Operation[getEdgeGatewayByIDParams, *types.ModelEdgeGateway]{
		Name:     opGetEdgeGateway,
		Backend:  cav.BackendVMware,
		Endpoint: endpoints.GetEdgeGateway(),
		Validate: func(p getEdgeGatewayByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p getEdgeGatewayByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.GetEdgeGateway()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(resp *cav.Response, _ getEdgeGatewayByIDParams) (*types.ModelEdgeGateway, error) {
			edgeGateway, ok := resp.Result().(*itypes.APIResponseEdgegateway)
			if !ok || edgeGateway == nil {
				return nil, fmt.Errorf("unexpected get response type %T", resp.Result())
			}

			return edgeGateway.ToModel(), nil
		},
	}
	deleteEdgeGatewayByIDOp = cav.Operation[deleteEdgeGatewayByIDParams, struct{}]{
		Name:     opDeleteEdgeGateway,
		Backend:  cav.BackendInfrapi,
		Endpoint: endpoints.DeleteEdgeGateway(),
		Validate: func(p deleteEdgeGatewayByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}

			return nil
		},
		RequestOptions: func(p deleteEdgeGatewayByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.DeleteEdgeGateway()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Extract: func(_ *cav.Response, _ deleteEdgeGatewayByIDParams) (struct{}, error) {
			return struct{}{}, nil
		},
	}
	updateEdgeGatewayByIDOp = cav.Operation[updateEdgeGatewayByIDParams, struct{}]{
		Name:     opUpdateEdgeGateway,
		Backend:  cav.BackendInfrapi,
		Endpoint: endpoints.UpdateEdgeGatewayBandwidth(),
		Validate: func(p updateEdgeGatewayByIDParams) error {
			if p.ID == "" {
				return fmt.Errorf("id is required")
			}
			if p.Bandwidth <= 0 {
				return fmt.Errorf("bandwidth must be greater than 0")
			}

			return nil
		},
		RequestOptions: func(p updateEdgeGatewayByIDParams) ([]cav.EndpointRequestOption, error) {
			ep := endpoints.UpdateEdgeGatewayBandwidth()
			return []cav.EndpointRequestOption{cav.WithPathParam(ep.PathParams[0], p.ID)}, nil
		},
		Transform: func(p updateEdgeGatewayByIDParams) (any, error) {
			return itypes.APIRequestBandwidth{Bandwidth: p.Bandwidth}, nil
		},
		Extract: func(_ *cav.Response, _ updateEdgeGatewayByIDParams) (struct{}, error) {
			return struct{}{}, nil
		},
	}
)

// GetEdgeGateway returns an edge gateway by ID or name.
func (c *Client) GetEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) (*types.ModelEdgeGateway, error) {
	if params.ID == "" && params.Name == "" {
		return nil, fmt.Errorf("%s: validate: id or name is required", opGetEdgeGateway)
	}

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve id: %w", opGetEdgeGateway, err)
		}
	}

	return cav.Execute(ctx, c.c, getEdgeGatewayByIDOp, getEdgeGatewayByIDParams{ID: params.ID})
}

// ListEdgeGateway lists edge gateways visible to current organization.
func (c *Client) ListEdgeGateway(ctx context.Context) (*types.ModelEdgeGateways, error) {
	return cav.Execute(ctx, c.c, listEdgeGatewayOp, struct{}{})
}

// DeleteEdgeGateway deletes an edge gateway by ID or name.
func (c *Client) DeleteEdgeGateway(ctx context.Context, params types.ParamsEdgeGateway) error {
	if params.ID == "" && params.Name == "" {
		return fmt.Errorf("%s: validate: id or name is required", opDeleteEdgeGateway)
	}

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return fmt.Errorf("%s: resolve id: %w", opDeleteEdgeGateway, err)
		}
	}

	if _, err := cav.Execute(ctx, c.c, deleteEdgeGatewayByIDOp, deleteEdgeGatewayByIDParams{ID: params.ID}); err != nil {
		return err
	}

	return nil
}

// UpdateEdgeGateway updates edge gateway bandwidth and returns refreshed state.
func (c *Client) UpdateEdgeGateway(ctx context.Context, params types.ParamsUpdateEdgeGateway) (*types.ModelEdgeGateway, error) {
	if params.ID == "" && params.Name == "" {
		return nil, fmt.Errorf("%s: validate: id or name is required", opUpdateEdgeGateway)
	}

	if params.ID == "" {
		var err error
		params.ID, err = c.retrieveEdgeGatewayIDByName(ctx, params.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve id: %w", opUpdateEdgeGateway, err)
		}
	}

	if _, err := cav.Execute(ctx, c.c, updateEdgeGatewayByIDOp, updateEdgeGatewayByIDParams{ID: params.ID, Bandwidth: params.Bandwidth}); err != nil {
		return nil, err
	}

	return c.GetEdgeGateway(ctx, types.ParamsEdgeGateway{ID: params.ID})
}
