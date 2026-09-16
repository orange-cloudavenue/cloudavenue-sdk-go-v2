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

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	pkgerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListEdgeGatewayNetworkContextProfile          = "EdgeGateway.NetworkContextProfile.List"
	opGetEdgeGatewayNetworkContextProfile           = "EdgeGateway.NetworkContextProfile.Get"
	opCreateEdgeGatewayNetworkContextProfile        = "EdgeGateway.NetworkContextProfile.Create"
	opUpdateEdgeGatewayNetworkContextProfile        = "EdgeGateway.NetworkContextProfile.Update"
	opDeleteEdgeGatewayNetworkContextProfile        = "EdgeGateway.NetworkContextProfile.Delete"
	opGetEdgeGatewayNetworkContextProfileAttributes = "EdgeGateway.NetworkContextProfile.Attributes.Get"
)

func edgeGatewayNetworkContextProfileFilterKey(ownerID string) (string, error) {
	if err := validators.New().Var(ownerID, "urn=vdcGroup"); err == nil {
		return "vdcGroupId", nil
	}
	if err := validators.New().Var(ownerID, "urn=vdc"); err == nil {
		return "orgVdcId", nil
	}

	return "", fmt.Errorf("edge gateway owner must be a vdc or vdc group")
}

func resolveEdgeGatewayNetworkContextProfileContext(ctx context.Context, c *Client, edgeGatewayID, edgeGatewayName string) (*itypes.APIObjectReference, string, string, error) {
	ep := endpoints.GetEdgeGateway()
	identifier := edgeGatewayID
	if identifier == "" {
		identifier = edgeGatewayName
	}

	resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], identifier))
	if err != nil {
		return nil, "", "", err
	}

	edgeGateway := resp.Result().(*itypes.APIResponseEdgegateway)
	if edgeGateway.OwnerRef == nil {
		return nil, "", "", fmt.Errorf("edge gateway owner is missing")
	}

	filterKey, err := edgeGatewayNetworkContextProfileFilterKey(edgeGateway.OwnerRef.ID)
	if err != nil {
		return nil, "", "", err
	}

	cd := cav.GetExtraDataFromContext(resp.Request.Context())
	return &itypes.APIObjectReference{ID: edgeGateway.OwnerRef.ID, Name: edgeGateway.OwnerRef.Name}, cd.OrganizationID, filterKey, nil
}

func findEdgeGatewayNetworkContextProfile(ctx context.Context, c *Client, id, name, edgeGatewayID, edgeGatewayName string) (*itypes.APIResponseNetworkContextProfile, error) {
	if id != "" {
		ep := endpoints.GetNetworkContextProfile()
		resp, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], id))
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.APIResponseNetworkContextProfile), nil
	}

	ownerRef, _, filterKey, err := resolveEdgeGatewayNetworkContextProfileContext(ctx, c, edgeGatewayID, edgeGatewayName)
	if err != nil {
		return nil, err
	}

	ep := endpoints.ListNetworkContextProfile()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%s==%s", filterKey, ownerRef.ID)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.APIResponseListNetworkContextProfile)
	var found []itypes.APIResponseNetworkContextProfile
	for _, profile := range list.Values {
		if profile.Name == name {
			found = append(found, profile)
		}
	}

	if len(found) == 0 {
		return nil, &pkgerrors.APIError{Operation: "GetNetworkContextProfile", StatusCode: 404, Message: fmt.Sprintf("network context profile %q not found", name)}
	}
	if len(found) > 1 {
		return nil, pkgerrors.Newf("found multiple network context profiles with the same name %q, use the ID to disambiguate", name)
	}

	return &found[0], nil
}

func toAPIEdgeGatewayNetworkContextProfileAttributes(attrs []types.ParamsNetworkContextProfileAttribute) []itypes.APINetworkContextProfileAttribute {
	out := make([]itypes.APINetworkContextProfileAttribute, 0, len(attrs))
	for _, attr := range attrs {
		apiAttr := itypes.APINetworkContextProfileAttribute{
			Type:   attr.Type,
			Values: attr.Values,
		}

		for _, sub := range attr.SubAttributes {
			apiAttr.SubAttributes = append(apiAttr.SubAttributes, itypes.APINetworkContextProfileSubAttribute{
				Type:   sub.Type,
				Values: sub.Values,
			})
		}

		out = append(out, apiAttr)
	}

	return out
}

func createEdgeGatewayNetworkContextProfileBody(ctx context.Context, c *Client, params types.ParamsCreateEdgeGatewayNetworkContextProfile) (itypes.APIRequestNetworkContextProfile, string, string, error) {
	ownerRef, orgID, filterKey, err := resolveEdgeGatewayNetworkContextProfileContext(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return itypes.APIRequestNetworkContextProfile{}, "", "", err
	}

	return itypes.APIRequestNetworkContextProfile{
		Name:            params.Name,
		Description:     params.Description,
		Scope:           types.NetworkContextProfileScopeTenant,
		ContextEntityID: ownerRef.ID,
		OrgRef:          &itypes.APIObjectReference{ID: orgID},
		Attributes:      toAPIEdgeGatewayNetworkContextProfileAttributes(params.Attributes),
	}, ownerRef.ID, filterKey, nil
}

// ListNetworkContextProfile lists network context profiles visible from an edge gateway owner.
func (c *Client) ListNetworkContextProfile(ctx context.Context, params types.ParamsListEdgeGatewayNetworkContextProfile) (*types.ModelListNetworkContextProfile, error) {
	ownerRef, _, filterKey, err := resolveEdgeGatewayNetworkContextProfileContext(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opListEdgeGatewayNetworkContextProfile, err)
	}

	ep := endpoints.ListNetworkContextProfile()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%s==%s", filterKey, ownerRef.ID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListEdgeGatewayNetworkContextProfile, err)
	}

	return resp.Result().(*itypes.APIResponseListNetworkContextProfile).ToModel(), nil
}

// GetNetworkContextProfile returns a network context profile by ID or name for an edge gateway owner.
func (c *Client) GetNetworkContextProfile(ctx context.Context, params types.ParamsGetEdgeGatewayNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	profile, err := findEdgeGatewayNetworkContextProfile(ctx, c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetEdgeGatewayNetworkContextProfile, err)
	}

	model := profile.ToModel()
	return &model, nil
}

// CreateNetworkContextProfile creates a tenant-scoped network context profile for an edge gateway owner.
func (c *Client) CreateNetworkContextProfile(ctx context.Context, params types.ParamsCreateEdgeGatewayNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	body, ownerID, filterKey, err := createEdgeGatewayNetworkContextProfileBody(ctx, c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateEdgeGatewayNetworkContextProfile, err)
	}

	if _, err := c.c.Do(
		ctx,
		endpoints.CreateNetworkContextProfile(),
		cav.SetBody(body),
	); err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateEdgeGatewayNetworkContextProfile, err)
	}

	created, err := findEdgeGatewayNetworkContextProfile(ctx, c, "", params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		if filterKey != "" && ownerID != "" {
			model := body.ToModel()
			return &model, nil
		}
		return nil, fmt.Errorf("%s: extract: %w", opCreateEdgeGatewayNetworkContextProfile, err)
	}

	model := created.ToModel()
	return &model, nil
}

// UpdateNetworkContextProfile replaces network context profile fields for an edge gateway owner.
func (c *Client) UpdateNetworkContextProfile(ctx context.Context, params types.ParamsUpdateEdgeGatewayNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	current, err := findEdgeGatewayNetworkContextProfile(ctx, c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateEdgeGatewayNetworkContextProfile, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	attributes := current.Attributes
	if len(params.Attributes) != 0 {
		attributes = toAPIEdgeGatewayNetworkContextProfileAttributes(params.Attributes)
	}

	body := itypes.APIRequestNetworkContextProfile{
		ID:              current.ID,
		Name:            current.Name,
		Description:     description,
		Scope:           current.Scope,
		ContextEntityID: current.ContextEntityID,
		OrgRef:          current.OrgRef,
		Attributes:      attributes,
	}

	ep := endpoints.UpdateNetworkContextProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
		cav.SetBody(body),
	); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateEdgeGatewayNetworkContextProfile, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteNetworkContextProfile deletes a network context profile from an edge gateway owner.
func (c *Client) DeleteNetworkContextProfile(ctx context.Context, params types.ParamsDeleteEdgeGatewayNetworkContextProfile) error {
	current, err := findEdgeGatewayNetworkContextProfile(ctx, c, params.ID, params.Name, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteEdgeGatewayNetworkContextProfile, err)
	}

	ep := endpoints.DeleteNetworkContextProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteEdgeGatewayNetworkContextProfile, err)
	}

	return nil
}

// GetNetworkContextProfileAttributes returns live APP_ID and DOMAIN_NAME values for an edge gateway owner.
func (c *Client) GetNetworkContextProfileAttributes(ctx context.Context, params types.ParamsGetEdgeGatewayNetworkContextProfileAttributes) (*types.ModelGetNetworkContextProfileAttributesCatalog, error) {
	ownerRef, _, filterKey, err := resolveEdgeGatewayNetworkContextProfileContext(ctx, c, params.EdgeGatewayID, params.EdgeGatewayName)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve edge gateway: %w", opGetEdgeGatewayNetworkContextProfileAttributes, err)
	}

	ep := endpoints.GetNetworkContextProfileAttributes()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%s==%s", filterKey, ownerRef.ID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: get attributes: %w", opGetEdgeGatewayNetworkContextProfileAttributes, err)
	}

	catalog := &types.ModelGetNetworkContextProfileAttributesCatalog{}
	for _, attr := range resp.Result().(*itypes.APINetworkContextProfileAttributesResponse).Attributes {
		switch attr.Type {
		case types.NetworkContextProfileAttributeTypeAppID:
			catalog.AppIDValues = append(catalog.AppIDValues, attr.Values...)
		case types.NetworkContextProfileAttributeTypeDomainName:
			catalog.DomainNameValues = append(catalog.DomainNameValues, attr.Values...)
		}
	}

	return catalog, nil
}
