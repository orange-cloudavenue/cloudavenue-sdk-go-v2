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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/inetworkobjects"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListAppPortProfile   = "AppPortProfile.List"
	opGetAppPortProfile    = "AppPortProfile.Get"
	opCreateAppPortProfile = "AppPortProfile.Create"
	opUpdateAppPortProfile = "AppPortProfile.Update"
	opDeleteAppPortProfile = "AppPortProfile.Delete"
)

// findAppPortProfile resolves an application port profile by ID or by name.
// Name lookups search tenant, provider, then system scopes.
func findAppPortProfile(ctx context.Context, cc *Client, idOrName, vdcGroupID string) (*itypes.ApiResponseAppPortProfile, error) {
	return inetworkobjects.FindAppPortProfile(ctx, cc.c, idOrName, vdcGroupID)
}

// ListAppPortProfile lists application port profiles visible from a VDC group.
func (c *Client) ListAppPortProfile(ctx context.Context, params types.ParamsListAppPortProfile) (*types.ModelListAppPortProfile, error) {
	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opListAppPortProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	ep := endpoints.ListAppPortProfile()
	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("_context==%s", vdcGroupID)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListAppPortProfile, err)
	}

	return resp.Result().(*itypes.ApiResponseListAppPortProfile).ToModel(), nil
}

// GetAppPortProfile returns an application port profile by ID or name for a VDC group.
func (c *Client) GetAppPortProfile(ctx context.Context, params types.ParamsGetAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" && params.VdcGroupName != "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdc group: %w", opGetAppPortProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	profile, err := findAppPortProfile(ctx, c, idOrName, vdcGroupID)
	if err != nil {
		return nil, fmt.Errorf("%s: get: %w", opGetAppPortProfile, err)
	}

	model := profile.ToModel()
	return &model, nil
}

// UpdateAppPortProfile replaces application port profile fields for a VDC group.
func (c *Client) UpdateAppPortProfile(ctx context.Context, params types.ParamsUpdateAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	if len(params.ApplicationPorts) != 0 {
		if err := inetworkobjects.ValidateAppPortProfileApplicationPorts(params.ApplicationPorts); err != nil {
			return nil, fmt.Errorf("%s: validate: %w", opUpdateAppPortProfile, err)
		}
	}

	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	current, err := findAppPortProfile(ctx, c, idOrName, "")
	if err != nil {
		return nil, fmt.Errorf("%s: resolve: %w", opUpdateAppPortProfile, err)
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
		return nil, fmt.Errorf("%s: update: %w", opUpdateAppPortProfile, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteAppPortProfile deletes an application port profile from a VDC group.
func (c *Client) DeleteAppPortProfile(ctx context.Context, params types.ParamsDeleteAppPortProfile) error {
	idOrName := params.ID
	if idOrName == "" {
		idOrName = params.Name
	}

	vdcGroupID := params.VdcGroupID
	if vdcGroupID == "" && params.VdcGroupName != "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.VdcGroupName})
		if err != nil {
			return fmt.Errorf("%s: resolve vdc group: %w", opDeleteAppPortProfile, err)
		}
		vdcGroupID = vdcGroup.ID
	}

	current, err := findAppPortProfile(ctx, c, idOrName, vdcGroupID)
	if err != nil {
		return fmt.Errorf("%s: resolve: %w", opDeleteAppPortProfile, err)
	}

	ep := endpoints.DeleteAppPortProfile()
	if _, err = c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], current.ID),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteAppPortProfile, err)
	}

	return nil
}

// CreateAppPortProfile creates a tenant-scoped application port profile for a VDC group.
func (c *Client) CreateAppPortProfile(ctx context.Context, params types.ParamsCreateAppPortProfile) (*types.ModelGetAppPortProfile, error) {
	vdcGroupID := params.VdcGroupID
	epList := endpoints.ListVdcGroup()
	filter := "id==" + vdcGroupID
	if vdcGroupID == "" {
		filter = "name==" + params.VdcGroupName
	}

	respList, err := c.c.Do(ctx, epList, cav.WithQueryParam(epList.QueryParams[0], filter))
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdc group: %w", opCreateAppPortProfile, err)
	}

	rL := respList.Result().(*itypes.ApiResponseListVdcGroup)
	if len(rL.Values) == 0 {
		return nil, errors.Newf("vdc group not found")
	}
	if vdcGroupID == "" {
		vdcGroupID = rL.Values[0].ID
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())

	body := itypes.ApiRequestAppPortProfile{
		Name:             params.Name,
		Description:      params.Description,
		ApplicationPorts: inetworkobjects.ToApiAppPortProfilePorts(params.ApplicationPorts),
		OrgRef:           &itypes.ApiObjectReference{ID: cd.OrganizationID},
		ContextEntityId:  vdcGroupID,
		Scope:            types.AppPortProfileScopeTenant,
	}

	ep := endpoints.CreateAppPortProfile()
	resp, err := c.c.Do(ctx, ep, cav.SetBody(body))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opCreateAppPortProfile, err)
	}

	profile, ok := resp.Result().(*itypes.ApiResponseAppPortProfile)
	if !ok || profile == nil {
		return nil, fmt.Errorf("%s: unexpected create response type %T", opCreateAppPortProfile, resp.Result())
	}

	model := profile.ToModel()
	return &model, nil
}
