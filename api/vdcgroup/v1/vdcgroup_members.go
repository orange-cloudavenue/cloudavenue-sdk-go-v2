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

func resolveVdcGroupParticipatingVdcs(ctx context.Context, c *Client, vdcs []types.ParamsCreateVdcGroupVdc, siteID string) ([]itypes.ApiResponseVdcGroupParticipatingVdc, error) {
	resolved := make([]itypes.ApiResponseVdcGroupParticipatingVdc, 0, len(vdcs))

	needsLookup := slices.ContainsFunc(vdcs, func(vdc types.ParamsCreateVdcGroupVdc) bool {
		return vdc.ID == ""
	})

	vdcIDsByName := make(map[string]string)
	if needsLookup {
		respListVdc, err := c.c.Do(ctx, endpoints.ListVdc())
		if err != nil {
			return nil, err
		}

		for _, vdc := range respListVdc.Result().(*itypes.ApiResponseListVDC).Records {
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

		resolved = append(resolved, itypes.ApiResponseVdcGroupParticipatingVdc{
			Vdc:                  itypes.ApiResponseVdcGroupParticipatingVdcRef{ID: resolvedID, Name: vdc.Name},
			Site:                 itypes.ApiResponseVdcGroupParticipatingSiteRef{ID: siteID},
			FaultDomainTag:       "AZ01",
			NetworkProviderScope: "AZ01",
		})
	}

	return resolved, nil
}

func mergeVdcGroupMembers(current []types.ModelGetVdcGroupVdc, additions []types.ParamsCreateVdcGroupVdc) []types.ParamsCreateVdcGroupVdc {
	vdcs := make([]types.ParamsCreateVdcGroupVdc, 0, len(additions)+len(current))
	vdcs = append(vdcs, additions...)

	for _, vdc := range current {
		vdcs = append(vdcs, types.ParamsCreateVdcGroupVdc{ID: vdc.ID, Name: vdc.Name})
	}

	return vdcs
}

func remainingVdcGroupMembers(current []types.ModelGetVdcGroupVdc, removals []types.ParamsCreateVdcGroupVdc) []types.ParamsCreateVdcGroupVdc {
	filtered := append([]types.ModelGetVdcGroupVdc(nil), current...)
	for _, vdc := range removals {
		filtered = slices.DeleteFunc(filtered, func(item types.ModelGetVdcGroupVdc) bool {
			return (vdc.ID != "" && item.ID == vdc.ID) || (vdc.Name != "" && item.Name == vdc.Name)
		})
	}

	vdcs := make([]types.ParamsCreateVdcGroupVdc, 0, len(filtered))
	for _, vdc := range filtered {
		vdcs = append(vdcs, types.ParamsCreateVdcGroupVdc{ID: vdc.ID, Name: vdc.Name})
	}

	return vdcs
}

// AddVdcToVdcGroup adds VDCs to an existing VDC group.
func (c *Client) AddVdcToVdcGroup(ctx context.Context, params types.ParamsAddVdcToVdcGroup) error {
	vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opAddVdcToVdcGroup, err)
	}

	var description *string
	if vdcGroup.Description != "" {
		description = &vdcGroup.Description
	}

	_, err = c.UpdateVdcGroup(ctx, types.ParamsUpdateVdcGroup{
		ID:          vdcGroup.ID,
		Name:        vdcGroup.Name,
		Description: description,
		Vdcs:        mergeVdcGroupMembers(vdcGroup.Vdcs, params.Vdcs),
	})
	if err != nil {
		return fmt.Errorf("%s: update: %w", opAddVdcToVdcGroup, err)
	}

	return nil
}

// RemoveVdcFromVdcGroup removes VDCs from an existing VDC group.
func (c *Client) RemoveVdcFromVdcGroup(ctx context.Context, params types.ParamsRemoveVdcFromVdcGroup) error {
	vdcGroup, err := c.GetVdcGroup(ctx, types.ParamsGetVdcGroup{ID: params.ID, Name: params.Name})
	if err != nil {
		return fmt.Errorf("%s: resolve target: %w", opRemoveVdcFromVdcGroup, err)
	}

	var description *string
	if vdcGroup.Description != "" {
		description = &vdcGroup.Description
	}

	_, err = c.UpdateVdcGroup(ctx, types.ParamsUpdateVdcGroup{
		ID:          vdcGroup.ID,
		Name:        vdcGroup.Name,
		Description: description,
		Vdcs:        remainingVdcGroupMembers(vdcGroup.Vdcs, params.Vdcs),
	})
	if err != nil {
		return fmt.Errorf("%s: update: %w", opRemoveVdcFromVdcGroup, err)
	}

	return nil
}
