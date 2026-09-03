/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcgroup

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListNetworkContextProfile          = "NetworkContextProfile.List"
	opGetNetworkContextProfile           = "NetworkContextProfile.Get"
	opCreateNetworkContextProfile        = "NetworkContextProfile.Create"
	opUpdateNetworkContextProfile        = "NetworkContextProfile.Update"
	opDeleteNetworkContextProfile        = "NetworkContextProfile.Delete"
	opGetNetworkContextProfileAttributes = "NetworkContextProfile.Attributes.Get"
)

// findNetworkContextProfile resolves a Network Context Profile by ID (direct
// GET, scope-agnostic) or by name (List filtered by vdcGroupId, then a
// client-side exact-name match; the List call already spans all scopes -
// SYSTEM, PROVIDER, TENANT - so no per-scope search is required, unlike
// Application Port Profile). Note: there is no registered URN format for
// this resource, so ID vs. name resolution is driven purely by which of the
// two params is non-empty, not by URN pattern-matching.
func findNetworkContextProfile(ctx context.Context, cc *Client, id, name, vdcGroupID string) (*itypes.ApiResponseNetworkContextProfile, error) {
	if id != "" {
		ep := endpoints.GetNetworkContextProfile()

		resp, err := cc.c.Do(
			ctx,
			ep,
			cav.WithPathParam(ep.PathParams[0], id),
		)
		if err != nil {
			return nil, err
		}

		return resp.Result().(*itypes.ApiResponseNetworkContextProfile), nil
	}

	ep := endpoints.ListNetworkContextProfile()

	resp, err := cc.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("vdcGroupId==%s", vdcGroupID)),
	)
	if err != nil {
		return nil, err
	}

	list := resp.Result().(*itypes.ApiResponseListNetworkContextProfile)

	var found []itypes.ApiResponseNetworkContextProfile
	for _, profile := range list.Values {
		if profile.Name == name {
			found = append(found, profile)
		}
	}

	if len(found) == 0 {
		return nil, &errors.APIError{Operation: "GetNetworkContextProfile", StatusCode: 404, Message: fmt.Sprintf("network context profile %q not found", name)}
	}
	if len(found) > 1 {
		return nil, errors.Newf("found multiple network context profiles with the same name %q, use the ID to disambiguate", name)
	}

	return &found[0], nil
}

func toApiNetworkContextProfileAttributes(attrs []types.ParamsNetworkContextProfileAttribute) []itypes.ApiNetworkContextProfileAttribute {
	out := make([]itypes.ApiNetworkContextProfileAttribute, 0, len(attrs))
	for _, attr := range attrs {
		apiAttr := itypes.ApiNetworkContextProfileAttribute{
			Type:   attr.Type,
			Values: attr.Values,
		}

		for _, sub := range attr.SubAttributes {
			apiAttr.SubAttributes = append(apiAttr.SubAttributes, itypes.ApiNetworkContextProfileSubAttribute{
				Type:   sub.Type,
				Values: sub.Values,
			})
		}

		out = append(out, apiAttr)
	}

	return out
}

func createNetworkContextProfileBody(ctx context.Context, cc *Client, params types.ParamsCreateNetworkContextProfile) (itypes.ApiRequestNetworkContextProfile, string, error) {
	vdcGroupID := params.VdcGroupID
	epList := endpoints.ListVdcGroup()
	filter := "id==" + vdcGroupID
	if vdcGroupID == "" {
		filter = "name==" + params.VdcGroupName
	}

	respList, err := cc.c.Do(
		ctx,
		epList,
		cav.WithQueryParam(epList.QueryParams[0], filter),
	)
	if err != nil {
		return itypes.ApiRequestNetworkContextProfile{}, "", err
	}

	rL := respList.Result().(*itypes.ApiResponseListVdcGroup)
	if len(rL.Values) == 0 {
		return itypes.ApiRequestNetworkContextProfile{}, "", errors.Newf("vdc group not found")
	}
	if vdcGroupID == "" {
		vdcGroupID = rL.Values[0].ID
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())

	return itypes.ApiRequestNetworkContextProfile{
		Name:            params.Name,
		Description:     params.Description,
		Scope:           types.NetworkContextProfileScopeTenant,
		ContextEntityId: vdcGroupID,
		OrgRef:          &itypes.ApiObjectReference{ID: cd.OrganizationID},
		Attributes:      toApiNetworkContextProfileAttributes(params.Attributes),
	}, vdcGroupID, nil
}

// ListNetworkContextProfile lists network context profiles visible from a VDC group.
func (c *Client) ListNetworkContextProfile(ctx context.Context, params types.ParamsListNetworkContextProfile) (*types.ModelListNetworkContextProfile, error) {
	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opListNetworkContextProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	ep := endpoints.ListNetworkContextProfile()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("vdcGroupId==%s", vdcGroupID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListNetworkContextProfile, err)
	}

	return resp.Result().(*itypes.ApiResponseListNetworkContextProfile).ToModel(), nil
}

// GetNetworkContextProfile returns a network context profile by ID or name for a VDC group.
func (c *Client) GetNetworkContextProfile(ctx context.Context, params types.ParamsGetNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	vdcGroupID := params.VdcGroupID
	if params.ID == "" && vdcGroupID == "" && params.VdcGroupName != "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opGetNetworkContextProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	profile, err := findNetworkContextProfile(ctx, c, params.ID, params.Name, vdcGroupID)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetNetworkContextProfile, err)
	}

	model := profile.ToModel()
	return &model, nil
}

// CreateNetworkContextProfile creates a tenant-scoped network context profile for a VDC group.
func (c *Client) CreateNetworkContextProfile(ctx context.Context, params types.ParamsCreateNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	body, vdcGroupID, err := createNetworkContextProfileBody(ctx, c, params)
	if err != nil {
		return nil, fmt.Errorf("%s: transform: %w", opCreateNetworkContextProfile, err)
	}

	// Create is asynchronous; underlying client blocks until job completion.
	if _, err := c.c.Do(
		ctx,
		endpoints.CreateNetworkContextProfile(),
		cav.SetBody(body),
	); err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateNetworkContextProfile, err)
	}

	created, err := findNetworkContextProfile(ctx, c, "", params.Name, vdcGroupID)
	if err != nil {
		return nil, fmt.Errorf("%s: extract: %w", opCreateNetworkContextProfile, err)
	}

	model := created.ToModel()
	return &model, nil
}

// UpdateNetworkContextProfile replaces network context profile fields for a VDC group.
func (c *Client) UpdateNetworkContextProfile(ctx context.Context, params types.ParamsUpdateNetworkContextProfile) (*types.ModelGetNetworkContextProfile, error) {
	current, err := findNetworkContextProfile(ctx, c, params.ID, "", "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateNetworkContextProfile, err)
	}

	description := current.Description
	if params.Description != "" {
		description = params.Description
	}

	attributes := current.Attributes
	if len(params.Attributes) != 0 {
		attributes = toApiNetworkContextProfileAttributes(params.Attributes)
	}

	body := itypes.ApiRequestNetworkContextProfile{
		ID:              current.ID,
		Name:            current.Name,
		Description:     description,
		Scope:           current.Scope,
		ContextEntityId: current.ContextEntityId,
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
		return nil, fmt.Errorf("%s: update: %w", opUpdateNetworkContextProfile, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteNetworkContextProfile deletes a network context profile from a VDC group.
func (c *Client) DeleteNetworkContextProfile(ctx context.Context, params types.ParamsDeleteNetworkContextProfile) error {
	vdcGroupID := params.VdcGroupID
	if params.ID == "" && vdcGroupID == "" && params.VdcGroupName != "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return fmt.Errorf("%s: resolve vdc group: %w", opDeleteNetworkContextProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	current, err := findNetworkContextProfile(ctx, c, params.ID, params.Name, vdcGroupID)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteNetworkContextProfile, err)
	}

	ep := endpoints.DeleteNetworkContextProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteNetworkContextProfile, err)
	}

	return nil
}

// GetNetworkContextProfileAttributes returns live APP_ID and DOMAIN_NAME values for a VDC group.
func (c *Client) GetNetworkContextProfileAttributes(ctx context.Context, params types.ParamsGetNetworkContextProfileAttributes) (*types.ModelGetNetworkContextProfileAttributesCatalog, error) {
	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opGetNetworkContextProfileAttributes, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	ep := endpoints.GetNetworkContextProfileAttributes()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("vdcGroupId==%s", vdcGroupID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: get attributes: %w", opGetNetworkContextProfileAttributes, err)
	}

	catalog := &types.ModelGetNetworkContextProfileAttributesCatalog{}
	for _, attr := range resp.Result().(*itypes.ApiNetworkContextProfileAttributesResponse).Attributes {
		switch attr.Type {
		case types.NetworkContextProfileAttributeTypeAppID:
			catalog.AppIDValues = append(catalog.AppIDValues, attr.Values...)
		case types.NetworkContextProfileAttributeTypeDomainName:
			catalog.DomainNameValues = append(catalog.DomainNameValues, attr.Values...)
		}
	}

	return catalog, nil
}
