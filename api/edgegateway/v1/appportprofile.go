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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListEdgeGatewayAppPortProfile   = "EdgeGateway.AppPortProfile.List"
	opGetEdgeGatewayAppPortProfile    = "EdgeGateway.AppPortProfile.Get"
	opCreateEdgeGatewayAppPortProfile = "EdgeGateway.AppPortProfile.Create"
	opUpdateEdgeGatewayAppPortProfile = "EdgeGateway.AppPortProfile.Update"
	opDeleteEdgeGatewayAppPortProfile = "EdgeGateway.AppPortProfile.Delete"
)

func resolveEdgeGatewayAppPortProfileContext(ctx context.Context, c cav.Client, edgeGatewayID, edgeGatewayName string) (*itypes.ApiObjectReference, string, error) {
	ep := endpoints.GetEdgeGateway()
	identifier := edgeGatewayID
	if identifier == "" {
		identifier = edgeGatewayName
	}

	resp, err := c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], identifier))
	if err != nil {
		return nil, "", err
	}

	edgeGateway := resp.Result().(*itypes.ApiResponseEdgegateway)
	if edgeGateway.OwnerRef == nil {
		return nil, "", fmt.Errorf("edge gateway owner is missing")
	}

	cd := cav.GetExtraDataFromContext(resp.Request.Context())
	return &itypes.ApiObjectReference{ID: edgeGateway.OwnerRef.ID, Name: edgeGateway.OwnerRef.Name}, cd.OrganizationID, nil
}

func createEdgeGatewayAppPortProfileBody(ctx context.Context, c cav.Client, params types.ParamsCreateEdgeGatewayAppPortProfile) (itypes.ApiRequestAppPortProfile, error) {
	if err := inetworkobjects.ValidateAppPortProfileApplicationPorts(params.ApplicationPorts); err != nil {
		return itypes.ApiRequestAppPortProfile{}, err
	}

	ownerRef, orgID, err := resolveEdgeGatewayAppPortProfileContext(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return itypes.ApiRequestAppPortProfile{}, err
	}

	return itypes.ApiRequestAppPortProfile{
		Name:             params.Name,
		Description:      params.Description,
		ApplicationPorts: inetworkobjects.ToApiAppPortProfilePorts(params.ApplicationPorts),
		OrgRef:           &itypes.ApiObjectReference{ID: orgID},
		ContextEntityId:  ownerRef.ID,
		Scope:            types.AppPortProfileScopeTenant,
	}, nil
}

// ListAppPortProfile lists application port profiles attached to an edge gateway.
func (c *Client) ListAppPortProfile(ctx context.Context, params types.ParamsListEdgeGatewayAppPortProfile) (*types.ModelListAppPortProfile, error) {
	ownerRef, _, err := resolveEdgeGatewayAppPortProfileContext(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opListEdgeGatewayAppPortProfile, err)
	}

	ep := endpoints.ListAppPortProfile()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("_context==%s", ownerRef.ID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListEdgeGatewayAppPortProfile, err)
	}

	return resp.Result().(*itypes.ApiResponseListAppPortProfile).ToModel(), nil
}

// GetAppPortProfile returns an application port profile by ID or name for an edge gateway.
func (c *Client) GetAppPortProfile(ctx context.Context, params types.ParamsGetEdgeGatewayAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	contextEntityID := ""
	if params.Name != "" {
		ownerRef, _, err := resolveEdgeGatewayAppPortProfileContext(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve edge gateway: %w", opGetEdgeGatewayAppPortProfile, err)
		}
		contextEntityID = ownerRef.ID
	}

	profile, err := inetworkobjects.FindAppPortProfile(ctx, c.c, idOrName, contextEntityID)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetEdgeGatewayAppPortProfile, err)
	}

	model := profile.ToModel()
	return &model, nil
}

// CreateAppPortProfile creates a tenant-scoped application port profile for an edge gateway.
func (c *Client) CreateAppPortProfile(ctx context.Context, params types.ParamsCreateEdgeGatewayAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	body, err := createEdgeGatewayAppPortProfileBody(ctx, c.c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateEdgeGatewayAppPortProfile, err)
	}

	ep := endpoints.CreateAppPortProfile()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGatewayAppPortProfile, err)
	}

	profile, ok := resp.Result().(*itypes.ApiResponseAppPortProfile)
	if !ok || profile == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateEdgeGatewayAppPortProfile, resp.Result())
	}

	model := profile.ToModel()
	return &model, nil
}

// UpdateAppPortProfile replaces application port profile fields for an edge gateway.
func (c *Client) UpdateAppPortProfile(ctx context.Context, params types.ParamsUpdateEdgeGatewayAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	if len(params.ApplicationPorts) != 0 {
		if err := inetworkobjects.ValidateAppPortProfileApplicationPorts(params.ApplicationPorts); err != nil {
			return nil, fmt.Errorf("%s: validate: %w", opUpdateEdgeGatewayAppPortProfile, err)
		}
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	contextEntityID := ""
	if params.ID == "" {
		ownerRef, _, err := resolveEdgeGatewayAppPortProfileContext(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve edge gateway: %w", opUpdateEdgeGatewayAppPortProfile, err)
		}
		contextEntityID = ownerRef.ID
	}

	current, err := inetworkobjects.FindAppPortProfile(ctx, c.c, idOrName, contextEntityID)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateEdgeGatewayAppPortProfile, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	applicationPorts := current.ApplicationPorts
	if len(params.ApplicationPorts) != 0 {
		applicationPorts = inetworkobjects.ToApiAppPortProfilePorts(params.ApplicationPorts)
	}

	body := itypes.ApiRequestAppPortProfile{
		ID:               current.ID,
		Name:             current.Name,
		Description:      description,
		ApplicationPorts: applicationPorts,
		OrgRef:           current.OrgRef,
		ContextEntityId:  current.ContextEntityId,
		Scope:            current.Scope,
	}

	ep := endpoints.UpdateAppPortProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
		cav.SetBody(body),
	); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateEdgeGatewayAppPortProfile, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteAppPortProfile deletes an application port profile from an edge gateway.
func (c *Client) DeleteAppPortProfile(ctx context.Context, params types.ParamsDeleteEdgeGatewayAppPortProfile) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	contextEntityID := ""
	if params.ID == "" {
		ownerRef, _, err := resolveEdgeGatewayAppPortProfileContext(ctx, c.c, params.EdgeGatewayID, params.EdgeGatewayName)
		if err != nil {
			return fmt.Errorf("%s: resolve edge gateway: %w", opDeleteEdgeGatewayAppPortProfile, err)
		}
		contextEntityID = ownerRef.ID
	}

	current, err := inetworkobjects.FindAppPortProfile(ctx, c.c, idOrName, contextEntityID)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteEdgeGatewayAppPortProfile, err)
	}

	ep := endpoints.DeleteAppPortProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteEdgeGatewayAppPortProfile, err)
	}

	return nil
}
