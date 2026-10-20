/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/itypes"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	opListStorageProfile   = "VDC.StorageProfile.List"
	opAddStorageProfile    = "VDC.StorageProfile.Add"
	opDeleteStorageProfile = "VDC.StorageProfile.Delete"
	opUpdateStorageProfile = "VDC.StorageProfile.Update"
)

// ListStorageProfile lists storage profiles for one or more VDCs.
func (c *Client) ListStorageProfile(ctx context.Context, params types.ParamsListStorageProfile) (*types.ModelListStorageProfiles, error) {
	ep := endpoints.ListStorageProfile()

	filters := make([]string, 0, 4)
	if params.ID != "" {
		filters = append(filters, "id=="+params.ID)
	}
	if params.Class != "" {
		filters = append(filters, "name=="+params.Class)
	}
	if params.VDCID != "" {
		filters = append(filters, "vdc=="+params.VDCID)
	}
	if params.VDCName != "" {
		filters = append(filters, "vdcName=="+params.VDCName)
	}

	var value strings.Builder
	for i, f := range filters {
		if i > 0 {
			value.WriteString(";")
		}
		value.WriteString(f)
	}

	resp, err := c.c.Do(ctx, ep, cav.WithQueryParam(ep.QueryParams[0], value.String()))
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opListStorageProfile, err)
	}

	model := resp.Result().(*itypes.APIResponseListStorageProfiles).ToModel()
	for _, vdc := range model.VDCS {
		if vdc.ID == "" {
			return nil, fmt.Errorf("%s: vdc ID is empty", opListStorageProfile)
		}
		for _, sp := range vdc.StorageProfiles {
			if sp.ID == "" {
				return nil, fmt.Errorf("%s: storage profile ID is empty", opListStorageProfile)
			}
		}
	}
	return model, nil
}

func uniqueVDCFromStorageProfiles(list *types.ModelListStorageProfiles) (*types.ModelListStorageProfilesVDC, error) {
	if len(list.VDCS) == 0 {
		return nil, errors.New("no VDC found with the provided ID or Name")
	}
	if len(list.VDCS) > 1 {
		return nil, errors.New("multiple VDCs found with the provided ID or Name, please specify a unique VDC")
	}
	return &list.VDCS[0], nil
}

// AddStorageProfile adds storage profiles to a VDC.
func (c *Client) AddStorageProfile(ctx context.Context, params types.ParamsAddStorageProfile) error {
	vdc, err := c.GetVDC(ctx, types.ParamsGetVDC{ID: params.VDCID, Name: params.VDCName})
	if err != nil {
		return fmt.Errorf("%s: get vdc: %w", opAddStorageProfile, err)
	}

	apiR := itypes.APIRequestUpdateVDC{
		VDC: itypes.APIRequestUpdateVDCVDC{Name: vdc.Name},
	}
	for _, sp := range params.StorageProfiles {
		apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.APIRequestVDCStorageProfile{
			Class:   sp.Class,
			Limit:   sp.Limit,
			Default: sp.Default,
		})
	}

	ep := endpoints.UpdateVDC()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], vdc.Name), cav.SetBody(apiR)); err != nil {
		return fmt.Errorf("%s: add: %w", opAddStorageProfile, err)
	}

	return nil
}

// DeleteStorageProfile removes storage profiles from a VDC.
func (c *Client) DeleteStorageProfile(ctx context.Context, params types.ParamsDeleteStorageProfile) error {
	listSP, err := c.ListStorageProfile(ctx, types.ParamsListStorageProfile{VDCName: params.VDCName, VDCID: params.VDCID})
	if err != nil {
		return fmt.Errorf("%s: list: %w", opDeleteStorageProfile, err)
	}
	vdc, err := uniqueVDCFromStorageProfiles(listSP)
	if err != nil {
		return fmt.Errorf("%s: resolve vdc: %w", opDeleteStorageProfile, err)
	}

	if len(vdc.StorageProfiles) == 1 {
		return fmt.Errorf("%s: %w", opDeleteStorageProfile, errors.New("cannot delete storage profile, at least one storage profile must exist for a VDC"))
	}

	apiR := itypes.APIRequestUpdateVDC{VDC: itypes.APIRequestUpdateVDCVDC{Name: vdc.Name}}
	apiR.VDC.StorageProfiles = make([]itypes.APIRequestVDCStorageProfile, 0, len(params.StorageProfiles))
	for _, pSP := range params.StorageProfiles {
		found := false
		for _, sp := range vdc.StorageProfiles {
			if sp.Class != pSP.Class {
				continue
			}
			found = true
			if sp.Default {
				return fmt.Errorf("%s: %w", opDeleteStorageProfile, errors.Newf("cannot delete the default storage profile %s from VDC %s", sp.Class, vdc.Name))
			}
			if sp.Used > 0 {
				return fmt.Errorf("%s: %w", opDeleteStorageProfile, errors.Newf("cannot delete a non-empty storage profile %s from VDC %s", sp.Class, vdc.Name))
			}
			apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.APIRequestVDCStorageProfile{Class: pSP.Class, Limit: 0, Default: false})
		}
		if !found {
			return fmt.Errorf("%s: %w", opDeleteStorageProfile, errors.Newf("storage profile class %s not found in VDC %s", pSP.Class, vdc.Name))
		}
	}

	ep := endpoints.UpdateVDC()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], vdc.Name), cav.SetBody(apiR)); err != nil {
		return fmt.Errorf("%s: delete: %w", opDeleteStorageProfile, err)
	}

	return nil
}

// UpdateStorageProfile updates storage profile limits and default flags for a VDC.
func (c *Client) UpdateStorageProfile(ctx context.Context, params types.ParamsUpdateStorageProfile) (*types.ModelListStorageProfilesVDC, error) {
	listSP, err := c.ListStorageProfile(ctx, types.ParamsListStorageProfile{VDCName: params.VDCName, VDCID: params.VDCID})
	if err != nil {
		return nil, fmt.Errorf("%s: list: %w", opUpdateStorageProfile, err)
	}
	vdc, err := uniqueVDCFromStorageProfiles(listSP)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve vdc: %w", opUpdateStorageProfile, err)
	}

	apiR := itypes.APIRequestUpdateVDC{VDC: itypes.APIRequestUpdateVDCVDC{Name: vdc.Name}}
	apiR.VDC.StorageProfiles = make([]itypes.APIRequestVDCStorageProfile, 0, len(params.StorageProfiles))

	currentStorageProfiles := make(map[string]types.ModelListStorageProfile, len(vdc.StorageProfiles))
	for _, sp := range vdc.StorageProfiles {
		currentStorageProfiles[sp.Class] = sp
	}

	defaultCount := 0
	for _, sp := range params.StorageProfiles {
		current, ok := currentStorageProfiles[sp.Class]
		if !ok {
			return nil, fmt.Errorf("%s: storage profile class %s not found in VDC %s", opUpdateStorageProfile, sp.Class, vdc.Name)
		}
		if sp.Limit < current.Used && sp.Limit > 0 {
			return nil, fmt.Errorf("%s: new limit for storage profile %s cannot be less than the current used (%d GiB)", opUpdateStorageProfile, sp.Class, current.Used)
		}

		defaultValue := current.Default
		if sp.Default != nil {
			defaultValue = *sp.Default
			if *sp.Default {
				defaultCount++
			}
		}

		apiR.VDC.StorageProfiles = append(apiR.VDC.StorageProfiles, itypes.APIRequestVDCStorageProfile{
			Class: sp.Class,
			Limit: func() int {
				if sp.Limit > 0 {
					return sp.Limit
				}
				return current.Limit
			}(),
			Default: defaultValue,
		})
	}

	if defaultCount > 1 {
		return nil, fmt.Errorf("%s: multiple storage profiles have default=true, only one is allowed", opUpdateStorageProfile)
	}

	ep := endpoints.UpdateVDC()
	if _, err := c.c.Do(ctx, ep, cav.WithPathParam(ep.PathParams[0], vdc.Name), cav.SetBody(apiR)); err != nil {
		return nil, fmt.Errorf("%s: update: %w", opUpdateStorageProfile, err)
	}

	updated, err := c.ListStorageProfile(ctx, types.ParamsListStorageProfile{VDCID: vdc.ID, VDCName: vdc.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: list updated: %w", opUpdateStorageProfile, err)
	}
	resolved, err := uniqueVDCFromStorageProfiles(updated)
	if err != nil {
		return nil, fmt.Errorf("%s: resolve updated: %w", opUpdateStorageProfile, err)
	}

	return resolved, nil
}
