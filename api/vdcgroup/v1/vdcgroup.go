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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListVdcGroup          = "VdcGroup.List"
	opGetVdcGroup           = "VdcGroup.Get"
	opCreateVdcGroup        = "VdcGroup.Create"
	opUpdateVdcGroup        = "VdcGroup.Update"
	opDeleteVdcGroup        = "VdcGroup.Delete"
	opAddVdcToVdcGroup      = "VdcGroup.Vdc.Add"
	opRemoveVdcFromVdcGroup = "VdcGroup.Vdc.Remove"
)

// ListVdcGroup lists VDC groups visible to current organization.
func (c *Client) ListVdcGroup(ctx context.Context, params types.ParamsListVdcGroup) (*types.ModelListVdcGroup, error) {
	ep := endpoints.ListVdcGroup()

	query := ""
	if params.Name != "" {
		query = fmt.Sprintf("name==%s", params.Name)
	}
	if params.ID != "" {
		query = fmt.Sprintf("id==%s", params.ID)
	}

	resp, err := c.c.Do(
		ctx,
		ep,
		cav.WithQueryParam(ep.QueryParams[0], query),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListVdcGroup, err)
	}

	return resp.Result().(*itypes.ApiResponseListVdcGroup).ToModel(), nil
}

// GetVdcGroup returns a VDC group by ID or name.
func (c *Client) GetVdcGroup(ctx context.Context, params types.ParamsGetVdcGroup) (*types.ModelGetVdcGroup, error) {
	vdcgroups, err := c.ListVdcGroup(ctx, types.ParamsListVdcGroup{
		ID:   params.ID,
		Name: params.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opGetVdcGroup, err)
	}

	if vdcgroups == nil || len(vdcgroups.VdcGroups) == 0 {
		return nil, fmt.Errorf("%s: vdc group not found", opGetVdcGroup)
	}

	matches := vdcgroups.VdcGroups
	if params.ID != "" || params.Name != "" {
		matches = make([]types.ModelGetVdcGroup, 0, len(vdcgroups.VdcGroups))
		for _, vdcGroup := range vdcgroups.VdcGroups {
			if params.ID != "" && vdcGroup.ID != params.ID {
				continue
			}
			if params.Name != "" && vdcGroup.Name != params.Name {
				continue
			}
			matches = append(matches, vdcGroup)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("%s: vdc group not found", opGetVdcGroup)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("%s: multiple vdc groups found", opGetVdcGroup)
	}

	return &matches[0], nil
}

// CreateVdcGroup creates a VDC group.
func (c *Client) CreateVdcGroup(ctx context.Context, params types.ParamsCreateVdcGroup) (*types.ModelGetVdcGroup, error) {
	epList := endpoints.ListVdcGroup()
	respList, err := c.c.Do(
		ctx,
		epList,
		cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", params.Name)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list existing: %w", opCreateVdcGroup, err)
	}

	listed := respList.Result().(*itypes.ApiResponseListVdcGroup).ToModel()
	for _, existing := range listed.VdcGroups {
		if existing.Name == params.Name {
			return nil, fmt.Errorf("%s: vdc group already exists", opCreateVdcGroup)
		}
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())
	vdcs, err := resolveVdcGroupParticipatingVdcs(ctx, c, params.Vdcs, cd.SiteID)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdcs: %w", opCreateVdcGroup, err)
	}

	body := itypes.ApiRequestCreateVdcGroup{
		OrgID:               cd.OrganizationID,
		Name:                params.Name,
		Description:         params.Description,
		Vdcs:                vdcs,
		NetworkProviderType: "NSX_T",
		Type:                "LOCAL",
	}

	if _, err := c.c.Do(ctx, endpoints.CreateVdcGroup(), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: create: %w", opCreateVdcGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// UpdateVdcGroup updates VDC group metadata or membership.
//
// When params.Vdcs is set, it becomes full desired membership.
// VDCs omitted from params.Vdcs are removed from group.
func (c *Client) UpdateVdcGroup(ctx context.Context, params types.ParamsUpdateVdcGroup) (*types.ModelGetVdcGroup, error) {
	current, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVdcGroup, err)
	}

	epList := endpoints.ListVdcGroup()
	param := cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("id==%s", current.ID))
	respList, err := c.c.Do(ctx, epList, param)
	if err != nil {
		return nil, fmt.Errorf("%s: list target: %w", opUpdateVdcGroup, err)
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())
	rL := respList.Result().(*itypes.ApiResponseListVdcGroup)
	if len(rL.Values) == 0 {
		return nil, fmt.Errorf("%s: target not found", opUpdateVdcGroup)
	}

	selected := rL.Values[0]
	for _, candidate := range rL.Values {
		if candidate.ID == current.ID {
			selected = candidate
			break
		}
	}

	body := itypes.ApiRequestUpdateVdcGroup{
		Id:    current.ID,
		OrgID: selected.OrgID,
		Name: func() string {
			if params.Name != "" {
				return params.Name
			}
			return current.Name
		}(),
		Description: func() string {
			if params.Description != nil {
				return *params.Description
			}
			return current.Description
		}(),
		NetworkProviderType: "NSX_T",
		Type:                "LOCAL",
	}

	if len(params.Vdcs) == 0 {
		body.Vdcs = selected.Vdcs
	} else {
		body.Vdcs, err = resolveVdcGroupParticipatingVdcs(ctx, c, params.Vdcs, cd.SiteID)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdcs: %w", opUpdateVdcGroup, err)
		}
	}

	ep := endpoints.UpdateVdcGroup()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateVdcGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVdcGroup deletes a VDC group.
func (c *Client) DeleteVdcGroup(ctx context.Context, params types.ParamsDeleteVdcGroup) error {
	id := params.ID
	if id == "" {
		vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{Name: params.Name})
		if err != nil {
			return fmt.Errorf("%s: resolve target: %w", opDeleteVdcGroup, err)
		}
		id = vdcGroup.ID
	}

	ep := endpoints.DeleteVdcGroup()
	if _, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], id),
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%t", params.Force)),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVdcGroup, err)
	}

	return nil
}
