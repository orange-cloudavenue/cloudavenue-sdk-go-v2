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
	"slices"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func resolveVDCGroupParticipatingVdcs(ctx context.Context, c *Client, vdcs []types.ParamsCreateVDCGroupVDC, siteID string) ([]itypes.APIResponseVDCGroupParticipatingVDC, error) {
	resolved := make([]itypes.APIResponseVDCGroupParticipatingVDC, 0, len(vdcs))

	needsLookup := slices.ContainsFunc(vdcs, func(vdc types.ParamsCreateVDCGroupVDC) bool {
		return vdc.ID == ""
	})

	vdcIDsByName := make(map[string]string)
	if needsLookup {
		respListVDC, err := c.c.Do(ctx, endpoints.ListVDC())
		if err != nil {
			return nil, err
		}

		for _, vdc := range respListVDC.Result().(*itypes.APIResponseListVDC).Records {
			vdcIDsByName[vdc.Name] = vdc.ID
		}
	}

	for _, vdc := range vdcs {
		resolvedID := vdc.ID
		if resolvedID == "" {
			resolvedID = vdcIDsByName[vdc.Name]
			if resolvedID == "" {
				return nil, fmt.Errorf("vdc %q not found", vdc.Name)
			}
		}

		resolved = append(resolved, itypes.APIResponseVDCGroupParticipatingVDC{
			VDC:                  itypes.APIResponseVDCGroupParticipatingVDCRef{ID: resolvedID, Name: vdc.Name},
			Site:                 itypes.APIResponseVDCGroupParticipatingSiteRef{ID: siteID},
			FaultDomainTag:       "AZ01",
			NetworkProviderScope: "AZ01",
		})
	}

	return resolved, nil
}

func mergeVDCGroupMembers(current []types.ModelGetVDCGroupVDC, additions []types.ParamsCreateVDCGroupVDC) []types.ParamsCreateVDCGroupVDC {
	vdcs := make([]types.ParamsCreateVDCGroupVDC, 0, len(additions)+len(current))
	vdcs = append(vdcs, additions...)

	for _, vdc := range current {
		vdcs = append(vdcs, types.ParamsCreateVDCGroupVDC(vdc))
	}

	return vdcs
}

func remainingVDCGroupMembers(current []types.ModelGetVDCGroupVDC, removals []types.ParamsCreateVDCGroupVDC) []types.ParamsCreateVDCGroupVDC {
	filtered := append([]types.ModelGetVDCGroupVDC(nil), current...)
	for _, vdc := range removals {
		filtered = slices.DeleteFunc(filtered, func(item types.ModelGetVDCGroupVDC) bool {
			return (vdc.ID != "" && item.ID == vdc.ID) || (vdc.Name != "" && item.Name == vdc.Name)
		})
	}

	vdcs := make([]types.ParamsCreateVDCGroupVDC, 0, len(filtered))
	for _, vdc := range filtered {
		vdcs = append(vdcs, types.ParamsCreateVDCGroupVDC{ID: vdc.ID, Name: vdc.Name})
	}

	return vdcs
}

// AddVDCToVDCGroup adds VDCs to an existing VDC group.
func (c *Client) AddVDCToVDCGroup(ctx context.Context, params types.ParamsAddVDCToVDCGroup) error {
	vdcGroup, err := c.GetVDCGroup(ctx, types.ParamsGetVDCGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opAddVDCToVDCGroup, err)
	}

	var description *string
	if vdcGroup.Description != "" {
		description = &vdcGroup.Description
	}

	_, err = c.UpdateVDCGroup(ctx, types.ParamsUpdateVDCGroup{
		ID:          vdcGroup.ID,
		Name:        vdcGroup.Name,
		Description: description,
		Vdcs:        mergeVDCGroupMembers(vdcGroup.Vdcs, params.Vdcs),
	})
	if err != nil {
		return fmt.Errorf("%s: update: %w", opAddVDCToVDCGroup, err)
	}

	return nil
}

// RemoveVDCFromVDCGroup removes VDCs from an existing VDC group.
func (c *Client) RemoveVDCFromVDCGroup(ctx context.Context, params types.ParamsRemoveVDCFromVDCGroup) error {
	vdcGroup, err := c.GetVDCGroup(ctx, types.ParamsGetVDCGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opRemoveVDCFromVDCGroup, err)
	}

	var description *string
	if vdcGroup.Description != "" {
		description = &vdcGroup.Description
	}

	_, err = c.UpdateVDCGroup(ctx, types.ParamsUpdateVDCGroup{
		ID:          vdcGroup.ID,
		Name:        vdcGroup.Name,
		Description: description,
		Vdcs:        remainingVDCGroupMembers(vdcGroup.Vdcs, params.Vdcs),
	})
	if err != nil {
		return fmt.Errorf("%s: update: %w", opRemoveVDCFromVDCGroup, err)
	}

	return nil
}
