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
	opListVDCGroup          = "VDCGroup.List"
	opGetVDCGroup           = "VDCGroup.Get"
	opCreateVDCGroup        = "VDCGroup.Create"
	opUpdateVDCGroup        = "VDCGroup.Update"
	opDeleteVDCGroup        = "VDCGroup.Delete"
	opAddVDCToVDCGroup      = "VDCGroup.Vdc.Add"
	opRemoveVDCFromVDCGroup = "VDCGroup.Vdc.Remove"
)

// ListVDCGroup lists VDC groups visible to current organization.
func (c *Client) ListVDCGroup(ctx context.Context, params types.ParamsListVDCGroup) (*types.ModelListVDCGroup, error) {
	ep := endpoints.ListVDCGroup()

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
		return nil, fmt.Errorf("%s: list: %w", opListVDCGroup, err)
	}

	return resp.Result().(*itypes.APIResponseListVDCGroup).ToModel(), nil
}

// GetVDCGroup returns a VDC group by ID or name.
func (c *Client) GetVDCGroup(ctx context.Context, params types.ParamsGetVDCGroup) (*types.ModelGetVDCGroup, error) {
	vdcgroups, err := c.ListVDCGroup(ctx, types.ParamsListVDCGroup(params))
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opGetVDCGroup, err)
	}

	if vdcgroups == nil || len(vdcgroups.VDCGroups) == 0 {
		return nil, fmt.Errorf("%s: vdc group not found", opGetVDCGroup)
	}

	matches := vdcgroups.VDCGroups
	if params.ID != "" || params.Name != "" {
		matches = make([]types.ModelGetVDCGroup, 0, len(vdcgroups.VDCGroups))
		for _, vdcGroup := range vdcgroups.VDCGroups {
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
		return nil, fmt.Errorf("%s: vdc group not found", opGetVDCGroup)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("%s: multiple vdc groups found", opGetVDCGroup)
	}

	return &matches[0], nil
}

// CreateVDCGroup creates a VDC group.
func (c *Client) CreateVDCGroup(ctx context.Context, params types.ParamsCreateVDCGroup) (*types.ModelGetVDCGroup, error) {
	epList := endpoints.ListVDCGroup()
	respList, err := c.c.Do(
		ctx,
		epList,
		cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("name==%s", params.Name)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: list existing: %w", opCreateVDCGroup, err)
	}

	listed := respList.Result().(*itypes.APIResponseListVDCGroup).ToModel()
	for _, existing := range listed.VDCGroups {
		if existing.Name == params.Name {
			return nil, fmt.Errorf("%s: vdc group already exists", opCreateVDCGroup)
		}
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())
	vdcs, err := resolveVDCGroupParticipatingVdcs(ctx, c, params.VDCs, cd.SiteID)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdcs: %w", opCreateVDCGroup, err)
	}

	body := itypes.APIRequestCreateVDCGroup{
		OrgID:               cd.OrganizationID,
		Name:                params.Name,
		Description:         params.Description,
		Vdcs:                vdcs,
		NetworkProviderType: "NSX_T",
		Type:                "LOCAL",
	}

	if _, err := c.c.Do(ctx, endpoints.CreateVDCGroup(), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: create: %w", opCreateVDCGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// UpdateVDCGroup updates VDC group metadata or membership.
//
// When params.Vdcs is set, it becomes full desired membership.
// VDCs omitted from params.Vdcs are removed from group.
func (c *Client) UpdateVDCGroup(ctx context.Context, params types.ParamsUpdateVDCGroup) (*types.ModelGetVDCGroup, error) {
	current, err := c.GetVDCGroup(ctx, types.ParamsGetVDCGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: resolve target: %w", opUpdateVDCGroup, err)
	}

	epList := endpoints.ListVDCGroup()
	param := cav.WithQueryParam(epList.QueryParams[0], fmt.Sprintf("id==%s", current.ID))
	respList, err := c.c.Do(ctx, epList, param)
	if err != nil {
		return nil, fmt.Errorf("%s: list target: %w", opUpdateVDCGroup, err)
	}

	cd := cav.GetExtraDataFromContext(respList.Request.Context())
	rL := respList.Result().(*itypes.APIResponseListVDCGroup)
	if len(rL.Values) == 0 {
		return nil, fmt.Errorf("%s: target not found", opUpdateVDCGroup)
	}

	selected := rL.Values[0]
	for _, candidate := range rL.Values {
		if candidate.ID == current.ID {
			selected = candidate
			break
		}
	}

	body := itypes.APIRequestUpdateVDCGroup{
		ID:    current.ID,
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
		body.Vdcs, err = resolveVDCGroupParticipatingVdcs(ctx, c, params.Vdcs, cd.SiteID)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve vdcs: %w", opUpdateVDCGroup, err)
		}
	}

	ep := endpoints.UpdateVDCGroup()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], current.ID), cav.SetBody(body)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateVDCGroup, err)
	}

	model := body.ToModel()
	return &model, nil
}

// DeleteVDCGroup deletes a VDC group.
func (c *Client) DeleteVDCGroup(ctx context.Context, params types.ParamsDeleteVDCGroup) error {
	id := params.ID
	if id == "" {
		vdcGroup, err := c.GetVDCGroup(ctx, types.ParamsGetVDCGroup{Name: params.Name})
		if err != nil {
			return fmt.Errorf("%s: resolve target: %w", opDeleteVDCGroup, err)
		}
		id = vdcGroup.ID
	}

	ep := endpoints.DeleteVDCGroup()
	if _, err := c.c.Do(
		ctx,
		ep,
		cav.WithPathParam(ep.PathParams[0], id),
		cav.WithQueryParam(ep.QueryParams[0], fmt.Sprintf("%t", params.Force)),
	); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteVDCGroup, err)
	}

	return nil
}
